package main

import (
	"context"
	"database/sql"
	"github.com/flipkart-incubator/diligent/pkg/buildinfo"
	"github.com/flipkart-incubator/diligent/pkg/datagen"
	"github.com/flipkart-incubator/diligent/pkg/idgen"
	"github.com/flipkart-incubator/diligent/pkg/intgen"
	"github.com/flipkart-incubator/diligent/pkg/metrics"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	"github.com/flipkart-incubator/diligent/pkg/work"
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

type DataContext struct {
	dataSpec *datagen.Spec
	dataGen  *datagen.DataGen
}

type DBContext struct {
	driver string
	url    string
	db     *sql.DB
}

type WorkloadContext struct {
	workloadName  string
	assignedRange *intgen.Range
	tableName     string
	durationSec   int
	concurrency   int
	batchSize     int
	workload      *work.Workload
}

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

	jobTracker *JobTracker
	data       *DataContext
	db         *DBContext
	workload   *WorkloadContext

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

		jobTracker: NewJobTracker(),

		data:     nil,
		db:       nil,
		workload: nil,
		metrics:  nil,
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
	currentJob := s.jobTracker.CurrentJob()
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

	err := s.jobTracker.Prepare(ctx, in.GetJobId(), in.GetJobSpec(), s.metrics)
	if err != nil {
		return &proto.MinionPrepareJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
			Pid: s.pid,
		}, nil
	}

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

	err := s.jobTracker.Run(ctx)
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

	err := s.jobTracker.Abort(ctx)
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

	jobInfo := s.jobTracker.GetJobInfo(in.GetJobId())
	if jobInfo == nil {
		return &proto.MinionQueryJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "no information found for job: " + in.GetJobId(),
			},
			Pid: s.pid,
		}, nil
	}

	// Query was successful
	log.Infof("GRPC: QueryJob() completed successfully")
	return &proto.MinionQueryJobResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
		Pid:     s.pid,
		JobInfo: jobInfo.ToProto(),
	}, nil
}
