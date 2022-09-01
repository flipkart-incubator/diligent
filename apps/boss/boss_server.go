package main

import (
	"context"
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/buildinfo"
	"github.com/flipkart-incubator/diligent/pkg/idgen"
	"github.com/flipkart-incubator/diligent/pkg/intgen"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

// BossServer represents a diligent boss gRPC server and associated state
// It is thread safe
type BossServer struct {
	proto.UnimplementedBossServer
	mut sync.Mutex // Must be taken for all operations on the BossServer

	listenAddr string
	pid        string
	startTime  time.Time
	nextJobNum int

	registry *MinionRegistry
}

func NewBossServer(listenAddr string) *BossServer {
	return &BossServer{
		listenAddr: listenAddr,
		pid:        idgen.GenerateId16(),
		startTime:  time.Now(),
		nextJobNum: 1,
		registry:   NewMinionRegistry(),
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

	err := s.registry.RegisterMinion(addr)
	if err != nil {
		log.Errorf("GRPC: RegisterMinion(%s): Failed to register (%s)", addr, err.Error())
		return &proto.BossRegisterMinionResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}, err
	}

	log.Infof("GRPC: RegisterMinion(%s): Successfully registered", addr)
	return &proto.BossRegisterMinionResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
	}, nil
}

func (s *BossServer) UnregisterMinion(_ context.Context, in *proto.BossUnregisterMinionRequest) (*proto.BossUnregisterMinionResponse, error) {
	addr := in.GetAddr()
	log.Infof("GRPC: UnregisterMinion(%s)", addr)
	s.mut.Lock()
	defer s.mut.Unlock()

	err := s.registry.UnregisterMinion(addr)
	if err != nil {
		log.Errorf("GRPC: UnregisterMinion(%s): Failed to unregister (%s)", addr, err.Error())
		return &proto.BossUnregisterMinionResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}, err
	}

	log.Infof("GRPC: UnregisterMinion(%s): Successfully unregistered", addr)
	return &proto.BossUnregisterMinionResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
	}, nil
}

func (s *BossServer) GetMinions(ctx context.Context, _ *proto.BossGetMinionsRequest) (*proto.BossGetMinionsResponse, error) {
	log.Infof("GRPC: ShowMinions()")
	s.mut.Lock()
	defer s.mut.Unlock()

	minionInfos := make([]*proto.MinionInfo, s.registry.GetNumMinions())
	ch := make(chan *proto.MinionInfo)

	for _, mm := range s.registry.GetMinionManagers() {
		go mm.GetMinionStatus(ctx, ch)
	}

	for i, _ := range minionInfos {
		minionInfos[i] = <-ch
	}

	return &proto.BossGetMinionsResponse{
		MinionInfos: minionInfos,
	}, nil
}

func (s *BossServer) PrepareJob(ctx context.Context, in *proto.BossPrepareJobRequest) (*proto.BossPrepareJobResponse, error) {
	log.Infof("GRPC: PrepareJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	// Generate the ID for this job
	jobId := s.getNextJobId()
	log.Infof("PrepareJob(): Assigning JobId=%s", jobId)

	// Partition the data among the number of minions
	dataSpec := proto.DataSpecFromProto(in.GetJobSpec().GetDataSpec())
	numRecs := dataSpec.KeyGenSpec.NumKeys()
	fullRange := intgen.NewRange(0, numRecs)
	numMinions := s.registry.GetNumMinions()
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
	minionStatuses := make([]*proto.MinionStatus, numMinions)
	ch := make(chan *proto.MinionStatus)

	i := 0
	for _, mm := range s.registry.GetMinionManagers() {
		log.Infof("PrepareJob(): Preparing Minion %s", mm.GetAddr())
		wlSpec := &proto.WorkloadSpec{
			WorkloadName:  in.GetJobSpec().GetWorkloadSpec().GetWorkloadName(),
			AssignedRange: proto.RangeToProto(assignedRanges[i]),
			TableName:     in.GetJobSpec().GetWorkloadSpec().GetTableName(),
			DurationSec:   in.GetJobSpec().GetWorkloadSpec().GetDurationSec(),
			Concurrency:   in.GetJobSpec().GetWorkloadSpec().GetConcurrency(),
			BatchSize:     in.GetJobSpec().GetWorkloadSpec().GetBatchSize(),
		}

		go mm.PrepareJobOnMinion(ctx, jobId, in.GetJobDesc(), in.GetJobSpec().GetDataSpec(), in.GetJobSpec().GetDbSpec(), wlSpec, ch)
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

	// Run on individual minions
	minionStatuses := make([]*proto.MinionStatus, s.registry.GetNumMinions())
	ch := make(chan *proto.MinionStatus)

	for _, mm := range s.registry.GetMinionManagers() {
		log.Infof("RunJob(): Triggering run on Minion %s", mm.GetAddr())
		go mm.RunJobOnMinion(ctx, ch)
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

	log.Infof("RunJob(): completed")
	return &proto.BossRunJobResponse{
		Status:         &overallStatus,
		MinionStatuses: minionStatuses,
	}, nil
}

func (s *BossServer) AbortJob(ctx context.Context, in *proto.BossAbortJobRequest) (*proto.BossAbortJobResponse, error) {
	log.Infof("GRPC: AbortJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	minionStatuses := make([]*proto.MinionStatus, s.registry.GetNumMinions())
	ch := make(chan *proto.MinionStatus)

	for _, mm := range s.registry.GetMinionManagers() {
		log.Infof("AbortJob(): Aborting on Minion %s", mm.GetAddr())
		go mm.AbortJobOnMinion(ctx, ch)
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

	log.Infof("AbortJob(): completed")
	return &proto.BossAbortJobResponse{
		Status:         &overallStatus,
		MinionStatuses: minionStatuses,
	}, nil
}

func (s *BossServer) QueryJob(ctx context.Context, in *proto.BossQueryJobRequest) (*proto.BossQueryJobResponse, error) {
	log.Infof("GRPC: QueryJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	minionJobInfos := make([]*proto.MinionJobInfo, s.registry.GetNumMinions())
	ch := make(chan *proto.MinionJobInfo)

	for _, mm := range s.registry.GetMinionManagers() {
		log.Infof("QueryJob(): Querying on Minion %s", mm.GetAddr())
		go mm.QueryJobOnMinion(ctx, in.GetJobId(), ch)
	}

	// Collect execution results
	for i, _ := range minionJobInfos {
		minionJobInfos[i] = <-ch
	}

	// Build overall status
	overallStatus := proto.GeneralStatus{
		IsOk:          true,
		FailureReason: "",
	}
	for _, ms := range minionJobInfos {
		if !ms.Status.IsOk {
			overallStatus.IsOk = false
			overallStatus.FailureReason = "errors encountered on one or more minions"
		}
	}

	log.Infof("QueryJob(): completed")
	return &proto.BossQueryJobResponse{
		Status:         &overallStatus,
		MinionJobInfos: minionJobInfos,
	}, nil
}

func (s *BossServer) getNextJobId() string {
	id := fmt.Sprintf("%d", s.nextJobNum)
	s.nextJobNum++
	return id
}
