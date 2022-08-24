package main

import (
	"context"
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/buildinfo"
	"github.com/flipkart-incubator/diligent/pkg/idgen"
	"github.com/flipkart-incubator/diligent/pkg/intgen"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	"github.com/processout/grpc-go-pool"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

const (
	minionConnIdleTimeoutSecs = 600
	minionConnMaxLifetimeSecs = 600
	minionDialTimeoutSecs     = 3
	minionRequestTimeout      = 5
)

//type JobState int
//const (
//	Prepared JobState = iota
//	Running
//	EndedSuccess
//	EndedFailure
//	EndedStopped
//	EndedNeverRan
//)
//
//type JobInfo struct {
//	jobId       string
//	jobSpec     *proto.JobSpec
//	jobState    JobState
//	prepareTime time.Time
//	runTime     time.Time
//	endTime     time.Time
//}

// BossServer represents a diligent boss gRPC server and associated state
// It is thread safe
type BossServer struct {
	proto.UnimplementedBossServer
	mut *sync.Mutex

	listenAddr string

	//currentJobInfo *JobInfo
	//pastJobInfos   []*JobInfo

	pid        string
	startTime  time.Time
	nextJobNum int

	minionPools map[string]*grpcpool.Pool
}

func NewBossServer(listenAddr string) *BossServer {
	return &BossServer{
		mut: &sync.Mutex{},

		listenAddr:  listenAddr,
		pid:         idgen.GenerateId16(),
		startTime:   time.Now(),
		nextJobNum:  1,
		minionPools: make(map[string]*grpcpool.Pool),
	}
}

func (s *BossServer) Serve() error {
	// Create listening port
	log.Infof("Creating listening port: %s", s.listenAddr)
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}

	// Create gRPC server
	log.Infof("Creating gRPC server")
	grpcServer := grpc.NewServer()

	// Register our this server object with the gRPC server
	log.Infof("Registering Boss Server")
	proto.RegisterBossServer(grpcServer, s)

	// Start serving
	log.Infof("Starting to serve")
	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	log.Infof("Diligent Boss server up and running")
	return nil
}

func (s *BossServer) Ping(_ context.Context, _ *proto.BossPingRequest) (*proto.BossPingResponse, error) {
	log.Infof("GRPC: Ping")
	s.mut.Lock()
	defer s.mut.Unlock()

	return &proto.BossPingResponse{
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
	}, nil
}

func (s *BossServer) RegisterMinion(_ context.Context, in *proto.BossRegisterMinionRequest) (*proto.BossRegisterMinionResponse, error) {
	addr := in.GetAddr()
	log.Infof("GRPC: RegisterMinion(%s)", addr)
	s.mut.Lock()
	defer s.mut.Unlock()

	// Acknowledge the request by making an entry for the minion, but first close any existing pool
	pool, ok := s.minionPools[addr]
	if ok && pool != nil {
		log.Infof("GRPC: RegisterMinion(%s): Existing entry found. Closing existing pool", addr)
		pool.Close()
	}
	s.minionPools[addr] = nil

	// Factory method for pool
	var factory grpcpool.Factory = func() (*grpc.ClientConn, error) {
		log.Infof("grpcpool.Factory(): Trying to connect to minion %s", addr)
		dialCtx, dialCancel := context.WithTimeout(context.Background(), minionDialTimeoutSecs*time.Second)
		defer dialCancel()
		conn, err := grpc.DialContext(dialCtx, addr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Errorf("grpcpool.Factory(): Failed to connect to minion %s (%s)", addr, err.Error())
			return nil, err
		}
		log.Infof("grpcpool.Factory(): Successfully connected to %s", addr)
		return conn, nil
	}

	// Create an empty connection pool (don't try to establish connection right now as minion may not be ready)
	pool, err := grpcpool.New(factory, 0, 1, minionConnIdleTimeoutSecs*time.Second, minionConnMaxLifetimeSecs*time.Second)
	if err != nil {
		log.Errorf("GRPC: RegisterMinion(%s): Failed to create pool (%s)", addr, err.Error())
		return &proto.BossRegisterMinionResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}, err
	}

	log.Infof("GRPC: RegisterMinion(%s): Successfully registered", addr)
	s.minionPools[addr] = pool
	return &proto.BossRegisterMinionResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
	}, nil
}

func (s *BossServer) UnregisterMinion(_ context.Context, in *proto.BossUnregisterMinonRequest) (*proto.BossUnregisterMinionResponse, error) {
	addr := in.GetAddr()
	log.Infof("GRPC: UnregisterMinion(%s)", addr)
	s.mut.Lock()
	defer s.mut.Unlock()

	pool, ok := s.minionPools[addr]
	if ok {
		log.Infof("GRPC: UnregisterMinion(%s): Found registered minion - removing", addr)
		delete(s.minionPools, addr)
		if pool != nil {
			log.Infof("GRPC: UnregisterMinion(%s): Found connection pool - closing", addr)
			pool.Close()
		} else {
			log.Infof("GRPC: UnregisterMinion(%s): Pool was nil, ignoring close", addr)
		}
	} else {
		log.Infof("GRPC: UnregisterMinion(%s): No such minion. Ignoring request", addr)
	}

	log.Infof("GRPC: Minion unregistered: %s", addr)
	return &proto.BossUnregisterMinionResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
	}, nil
}

func (s *BossServer) ShowMinions(ctx context.Context, _ *proto.BossShowMinionRequest) (*proto.BossShowMinionResponse, error) {
	log.Infof("GRPC: ShowMinions()")
	s.mut.Lock()
	defer s.mut.Unlock()

	minionInfos := make([]*proto.MinionInfo, len(s.minionPools))
	ch := make(chan *proto.MinionInfo)

	for addr := range s.minionPools {
		go s.getMinionStatus(ctx, addr, ch)
	}

	for i, _ := range minionInfos {
		minionInfos[i] = <-ch
	}

	return &proto.BossShowMinionResponse{
		MinionInfos: minionInfos,
	}, nil
}

func (s *BossServer) PrepareJob(ctx context.Context, in *proto.BossPrepareJobRequest) (*proto.BossPrepareJobResponse, error) {
	log.Infof("GRPC: PrepareJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	// Reject if we are running a job currently
	//if s.currentJobInfo != nil && s.currentJobInfo.jobState == Running {
	//	log.Infof("PrepareJob(): Cannot prepare. Job %s is running", s.currentJobInfo.jobId)
	//	return &proto.BossPrepareJobResponse{
	//		Status: &proto.GeneralStatus{
	//			IsOk:          false,
	//			FailureReason: "current job is running: " + s.currentJobInfo.jobId,
	//		},
	//	}, nil
	//}

	// Clean up earlier state if we were already in prepared state
	//if s.currentJobInfo != nil && s.currentJobInfo.jobState == Prepared {
	//	log.Infof("PrepareJob(): Boss was prepared for job %s. Cleaning up previous state", s.currentJobInfo.jobId)
	//	s.currentJobInfo.jobState = EndedNeverRan
	//	s.currentJobInfo.endTime = time.Now()
	//	s.pastJobInfos = append(s.pastJobInfos, s.currentJobInfo)
	//	s.currentJobInfo = nil
	//}

	// Generate the ID for this job
	jobId := s.getNextJobId()
	log.Infof("PrepareJob(): Assigning JobId=%s", jobId)

	// Partition the data among the number of minions
	dataSpec := proto.DataSpecFromProto(in.GetJobSpec().GetDataSpec())
	numRecs := dataSpec.KeyGenSpec.NumKeys()
	fullRange := intgen.NewRange(0, numRecs)
	numMinions := len(s.minionPools)
	var assignedRanges []*intgen.Range

	// Get the job level workload spec, and do homework to determine per-minion workload spec
	workloadName := in.GetJobSpec().GetWorkloadSpec().GetWorkloadName()
	switch workloadName {
	case "insert", "insert-txn", "delete", "delete-txn":
		assignedRanges = fullRange.Partition(numMinions)
	case "select-pk", "select-pk-txn", "select-uk", "select-uk-txn", "update", "update-txn":
		assignedRanges = fullRange.Duplicate(numMinions)
	default:
		reason := fmt.Sprintf("invalid workload '%s'", workloadName)
		log.Infof("RunWorkload(): %s", reason)
		return &proto.BossPrepareJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}, fmt.Errorf(reason)
	}
	log.Infof("PrepareJob(): Data partitioned among minions: %v", assignedRanges)

	// Prepare individual minions
	minionStatuses := make([]*proto.MinionStatus, len(s.minionPools))
	ch := make(chan *proto.MinionStatus)

	i := 0
	for addr := range s.minionPools {
		log.Infof("PrepareJob(): Preparing Minion %s", addr)
		wlSpec := &proto.WorkloadSpec{
			WorkloadName:  in.GetJobSpec().GetWorkloadSpec().GetWorkloadName(),
			AssignedRange: proto.RangeToProto(assignedRanges[i]),
			TableName:     in.GetJobSpec().GetWorkloadSpec().GetTableName(),
			DurationSec:   in.GetJobSpec().GetWorkloadSpec().GetDurationSec(),
			Concurrency:   in.GetJobSpec().GetWorkloadSpec().GetConcurrency(),
			BatchSize:     in.GetJobSpec().GetWorkloadSpec().GetBatchSize(),
		}

		go s.prepareJobOnMinion(ctx, addr, jobId, in.GetJobDesc(), in.GetJobSpec().GetDataSpec(), in.GetJobSpec().GetDbSpec(), wlSpec, ch)
		i++
	}

	// Collect execution results
	for i, _ = range minionStatuses {
		minionStatuses[i] = <-ch
	}

	// Build overall status
	overallStatus := proto.GeneralStatus{
		IsOk:          true,
		FailureReason: "",
	}
	for _, ms := range minionStatuses {
		if !ms.Status.IsOk {
			overallStatus.IsOk = false
			overallStatus.FailureReason = "errors encountered on one or more minions"
		}
	}

	// Preparation was successful
	//s.currentJobInfo = &JobInfo{
	//	jobId:       jobId,
	//	jobSpec:     in.GetJobSpec(),
	//	jobState:    Prepared,
	//	prepareTime: time.Now(),
	//}
	log.Infof("PrepareJob(): completed")
	return &proto.BossPrepareJobResponse{
		Status:         &overallStatus,
		JobId:          jobId,
		MinionStatuses: minionStatuses,
	}, nil
}

func (s *BossServer) RunJob(ctx context.Context, in *proto.BossRunJobRequest) (*proto.BossRunJobResponse, error) {
	log.Infof("GRPC: RunJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	// Process only if we are in prepared state
	//if s.currentJobInfo == nil || s.currentJobInfo.jobState != Prepared {
	//	return &proto.BossRunJobResponse{
	//		Status: &proto.GeneralStatus{
	//			IsOk:          false,
	//			FailureReason: "no current job in prepared state",
	//		},
	//	}, nil
	//}

	// Run on individual minions
	minionStatuses := make([]*proto.MinionStatus, len(s.minionPools))
	ch := make(chan *proto.MinionStatus)

	for addr := range s.minionPools {
		log.Infof("RunJob(): Triggering run on Minion %s", addr)
		go s.runJobOnMinion(ctx, addr, ch)
	}

	// Collect execution results
	for i, _ := range minionStatuses {
		minionStatuses[i] = <-ch
	}

	// Build overall status
	overallStatus := proto.GeneralStatus{
		IsOk:          true,
		FailureReason: "",
	}
	for _, ms := range minionStatuses {
		if !ms.Status.IsOk {
			overallStatus.IsOk = false
			overallStatus.FailureReason = "errors encountered on one or more minions"
		}
	}

	//s.currentJobInfo.jobState = Running
	//s.currentJobInfo.runTime = time.Now()
	log.Infof("RunJob(): completed")
	return &proto.BossRunJobResponse{
		Status:         &overallStatus,
		MinionStatuses: minionStatuses,
	}, nil
}

func (s *BossServer) StopJob(ctx context.Context, in *proto.BossStopJobRequest) (*proto.BossStopJobResponse, error) {
	log.Infof("GRPC: StopJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	// Process only if minion is running a job and id matches
	//if s.currentJobInfo == nil || s.currentJobInfo.jobState != Running {
	//	return &proto.BossStopJobResponse{
	//		Status: &proto.GeneralStatus{
	//			IsOk:          false,
	//			FailureReason: "no current job in running state",
	//		},
	//	}, nil
	//}

	minionStatuses := make([]*proto.MinionStatus, len(s.minionPools))
	ch := make(chan *proto.MinionStatus)

	for addr := range s.minionPools {
		log.Infof("StopJob(): Stopping on Minion %s", addr)
		go s.stopJobOnMinion(ctx, addr, ch)
	}

	// Collect execution results
	for i, _ := range minionStatuses {
		minionStatuses[i] = <-ch
	}

	// Build overall status
	overallStatus := proto.GeneralStatus{
		IsOk:          true,
		FailureReason: "",
	}
	for _, ms := range minionStatuses {
		if !ms.Status.IsOk {
			overallStatus.IsOk = false
			overallStatus.FailureReason = "errors encountered on one or more minions"
		}
	}

	//s.currentJobInfo.jobState = EndedStopped
	//s.currentJobInfo.endTime = time.Now()
	log.Infof("StopJob(): completed")
	return &proto.BossStopJobResponse{
		Status:         &overallStatus,
		MinionStatuses: minionStatuses,
	}, nil
}

func (s *BossServer) getConnectionForMinion(ctx context.Context, minionAddr string) (*grpcpool.ClientConn, error) {
	pool, ok := s.minionPools[minionAddr]

	if !ok {
		return nil, fmt.Errorf("no such minion: %s", minionAddr)
	}

	if pool == nil {
		return nil, fmt.Errorf("no pool created for minion: %s", minionAddr)
	}

	conn, err := pool.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection for minion: %s (%s)", minionAddr, err.Error())
	}

	return conn, nil
}

func (s *BossServer) getMinionStatus(ctx context.Context, minionAddr string, ch chan *proto.MinionInfo) {
	minionConnection, err := s.getConnectionForMinion(ctx, minionAddr)
	if err != nil {
		reason := fmt.Sprintf("failed to get client for minion %s (%s)", minionAddr, err.Error())
		log.Errorf("getMinionStatus(): %s", reason)
		ch <- &proto.MinionInfo{
			Addr: minionAddr,
			Reachability: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}
		return
	}
	defer minionConnection.Close()
	minionClient := proto.NewMinionClient(minionConnection)

	// Ping Minion
	pingCtx, pingCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	defer pingCtxCancel()
	res, err := minionClient.Ping(pingCtx, &proto.MinionPingRequest{})
	if err != nil {
		reason := fmt.Sprintf("gRPC request to %s failed (%v)", minionAddr, err)
		log.Errorf("getMinionStatus(): %s", reason)
		ch <- &proto.MinionInfo{
			Addr: minionAddr,
			Reachability: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}
		return
	}

	log.Infof("getMinionStatus(): Ping request to %s successful", minionAddr)
	ch <- &proto.MinionInfo{
		Addr: minionAddr,
		Reachability: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
		BuildInfo:   res.GetBuildInfo(),
		ProcessInfo: res.GetProcessInfo(),
		JobInfo:     res.GetJobInfo(),
	}
}

func (s *BossServer) prepareJobOnMinion(ctx context.Context, addr, jobId, jobDesc string,
	dataSpec *proto.DataSpec, dbSpec *proto.DBSpec, wlSpec *proto.WorkloadSpec, ch chan *proto.MinionStatus) {

	minionConnection, err := s.getConnectionForMinion(ctx, addr)
	if err != nil {
		ch <- &proto.MinionStatus{
			Addr: addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: fmt.Sprintf("failed to get client: %s", err.Error()),
			},
		}
		return
	}
	defer minionConnection.Close()

	minionClient := proto.NewMinionClient(minionConnection)
	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	res, err := minionClient.PrepareJob(grpcCtx, &proto.MinionPrepareJobRequest{
		JobId:   jobId,
		JobDesc: jobDesc,
		JobSpec: &proto.JobSpec{
			DataSpec:     dataSpec,
			DbSpec:       dbSpec,
			WorkloadSpec: wlSpec,
		},
	})
	grpcCtxCancel()

	if err != nil {
		ch <- &proto.MinionStatus{
			Addr: addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}
		return
	}

	if res.GetStatus().GetIsOk() != true {
		ch <- &proto.MinionStatus{
			Addr: addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: res.GetStatus().GetFailureReason(),
			},
		}
		return
	}

	log.Infof("Successfully prepared minion %s", addr)
	ch <- &proto.MinionStatus{
		Addr: addr,
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
	}
	return
}

func (s *BossServer) runJobOnMinion(ctx context.Context, addr string, ch chan *proto.MinionStatus) {
	minionConnection, err := s.getConnectionForMinion(ctx, addr)
	if err != nil {
		ch <- &proto.MinionStatus{
			Addr: addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: fmt.Sprintf("failed to get client: %s", err.Error()),
			},
		}
		return
	}
	defer minionConnection.Close()

	minionClient := proto.NewMinionClient(minionConnection)
	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	res, err := minionClient.RunJob(grpcCtx, &proto.MinionRunJobRequest{})
	grpcCtxCancel()

	if err != nil {
		ch <- &proto.MinionStatus{
			Addr: addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}
		return
	}

	if res.GetStatus().GetIsOk() != true {
		ch <- &proto.MinionStatus{
			Addr: addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: res.GetStatus().GetFailureReason(),
			},
		}
		return
	}

	log.Infof("Successfully triggered run on minion %s", addr)
	ch <- &proto.MinionStatus{
		Addr: addr,
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
	}
	return
}

func (s *BossServer) stopJobOnMinion(ctx context.Context, addr string, ch chan *proto.MinionStatus) {
	minionConnection, err := s.getConnectionForMinion(ctx, addr)
	if err != nil {
		ch <- &proto.MinionStatus{
			Addr: addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: fmt.Sprintf("failed to get client: %s", err.Error()),
			},
		}
		return
	}
	defer minionConnection.Close()

	minionClient := proto.NewMinionClient(minionConnection)

	// Stop workload
	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	stopWorkloadResponse, err := minionClient.StopJob(grpcCtx, &proto.MinionStopJobRequest{})
	grpcCtxCancel()

	if err != nil {
		ch <- &proto.MinionStatus{
			Addr: addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}
		return
	}

	if stopWorkloadResponse.GetStatus().GetIsOk() != true {
		ch <- &proto.MinionStatus{
			Addr: addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: stopWorkloadResponse.GetStatus().GetFailureReason(),
			},
		}
		return
	}

	log.Infof("Successfully triggered stop on minion %s", addr)
	ch <- &proto.MinionStatus{
		Addr: addr,
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
	}
	return
}

func (s *BossServer) getNextJobId() string {
	id := fmt.Sprintf("%d", s.nextJobNum)
	s.nextJobNum++
	return id
}
