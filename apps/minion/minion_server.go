package main

import (
	"context"
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/buildinfo"
	"github.com/flipkart-incubator/diligent/pkg/idgen"
	"github.com/flipkart-incubator/diligent/pkg/metrics"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

const (
	bossConnectionTimoutSecs = 5
	bossRequestTimeoutSecs   = 5
)

// MinionServer represents a diligent minion gRPC server and associated state
// It is thread safe
type MinionServer struct {
	proto.UnimplementedMinionServer
	mut sync.Mutex // Must be taken for all operations on the MinionServer

	grpcAddr      string
	metricsAddr   string
	advertiseAddr string
	bossAddr      string

	pid       string    //pid is a string that represents a unique run of the minion server
	startTime time.Time //startTime is when this server started executing

	job     *Job
	metrics *metrics.DiligentMetrics
}

func NewMinionServer(grpcAddr, metricsAddr, advertiseAddr, bossAddr string) *MinionServer {
	return &MinionServer{
		grpcAddr:      grpcAddr,
		metricsAddr:   metricsAddr,
		advertiseAddr: advertiseAddr,
		bossAddr:      bossAddr,

		pid:       idgen.GenerateId16(),
		startTime: time.Now(),

		metrics: nil,
	}
}

func (s *MinionServer) RegisterWithBoss() error {
	log.Infof("Registering with boss")
	var conn *grpc.ClientConn
	var err error
	for {
		dialCtx, dialCancel := context.WithTimeout(context.Background(), bossConnectionTimoutSecs*time.Second)
		conn, err = grpc.DialContext(dialCtx, s.bossAddr, grpc.WithInsecure(), grpc.WithBlock())
		dialCancel()
		if err != nil {
			log.Errorf("failed to connect to boss %s (%v)", s.bossAddr, err)
			continue
		}

		bossClient := proto.NewBossClient(conn)

		grcpCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
		_, err = bossClient.RegisterMinion(grcpCtx, &proto.BossRegisterMinionRequest{Addr: s.advertiseAddr})
		grpcCancel()
		if err != nil {
			log.Errorf("Registration failed with error: (%v)\n", err)
			continue
		} else {
			break
		}
	}

	log.Infof("Successfully registered with boss")
	return nil
}

func (s *MinionServer) Serve() error {
	log.Infof("Starting to serve...")

	// Start metrics serving
	s.metrics = metrics.NewDiligentMetrics(s.metricsAddr)
	s.metrics.Register()

	// Create listening port
	listener, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		return err
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register our this server object with the gRPC server
	proto.RegisterMinionServer(grpcServer, s)

	// Start serving
	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func (s *MinionServer) Ping(_ context.Context, in *proto.MinionPingRequest) (*proto.MinionPingResponse, error) {
	log.Infof("GRPC: Ping")
	s.mut.Lock()
	defer s.mut.Unlock()

	log.Infof("GRPC: Ping completed successfully")
	var protoJobInfo *proto.JobInfo = nil
	currentJob := s.job
	if currentJob != nil {
		protoJobInfo = currentJob.Info().ToProto()
	}
	return &proto.MinionPingResponse{
		BuildInfo: &proto.BuildInfo{
			AppName:    buildinfo.AppName,
			AppVersion: buildinfo.AppVersion,
			CommitHash: buildinfo.CommitHash,
			GoVersion:  buildinfo.GoVersion,
			BuildTime:  buildinfo.BuildTime,
		},
		ProcessInfo: &proto.ProcessInfo{
			Pid:       s.pid,
			StartTime: s.startTime.Format(time.UnixDate),
			Uptime:    time.Since(s.startTime).String(),
		},
		JobInfo: protoJobInfo,
	}, nil
}

func (s *MinionServer) PrepareJob(ctx context.Context, in *proto.MinionPrepareJobRequest) (*proto.MinionPrepareJobResponse, error) {
	log.Infof("GRPC: PrepareJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	// Cannot run prepare if we have a job in running state. Either no job, or job in an ended state are acceptable
	if s.job != nil && s.job.State() == Running {
		return &proto.MinionPrepareJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: fmt.Sprintf("job with id=%s is currently running", s.job.Id()),
			},
			Pid: s.pid,
		}, nil
	}

	job, err := PrepareJob(ctx, in.GetJobId(), in.GetJobSpec(), s.metrics)
	if err != nil {
		return &proto.MinionPrepareJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
			Pid: s.pid,
		}, nil
	}
	s.job = job

	// Preparation was successful
	log.Infof("GRPC: PrepareJob(jobId=%s) completed successfully", in.GetJobId())
	return &proto.MinionPrepareJobResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
		Pid: s.pid,
	}, nil
}

func (s *MinionServer) RunJob(ctx context.Context, in *proto.MinionRunJobRequest) (*proto.MinionRunJobResponse, error) {
	log.Infof("GRPC: RunJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	// Cannot run if we have no current job
	if s.job == nil {
		return &proto.MinionRunJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "no current job",
			},
			Pid: s.pid,
		}, nil
	}

	_, err := s.job.Run(ctx)
	if err != nil {
		return &proto.MinionRunJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
			Pid: s.pid,
		}, nil
	}

	// Run was successful
	log.Infof("GRPC: RunJob() completed successfully")
	return &proto.MinionRunJobResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
		Pid: s.pid,
	}, nil
}

func (s *MinionServer) AbortJob(ctx context.Context, in *proto.MinionAbortJobRequest) (*proto.MinionAbortJobResponse, error) {
	log.Infof("GRPC: AbortJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	// Cannot abort if we have no current job
	if s.job == nil {
		return &proto.MinionAbortJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "no current job",
			},
			Pid: s.pid,
		}, nil
	}

	err := s.job.Abort(ctx)
	if err != nil {
		return &proto.MinionAbortJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
			Pid: s.pid,
		}, nil
	}

	// Abort was successful
	log.Infof("GRPC: AbortJob() completed successfully")
	return &proto.MinionAbortJobResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
		Pid: s.pid,
	}, nil
}

func (s *MinionServer) QueryJob(ctx context.Context, in *proto.MinionQueryJobRequest) (*proto.MinionQueryJobResponse, error) {
	log.Infof("GRPC: QueryJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	log.Infof("GRPC: QueryJob() completed successfully")
	var protoJobInfo *proto.JobInfo = nil
	currentJob := s.job
	if currentJob != nil {
		protoJobInfo = currentJob.Info().ToProto()
	}
	return &proto.MinionQueryJobResponse{
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
		Pid:     s.pid,
		JobInfo: protoJobInfo,
	}, nil
}
