package main

import (
	"context"
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/buildinfo"
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

// BossServer represents a diligent boss gRPC server and associated state
// It is thread safe
type BossServer struct {
	proto.UnimplementedBossServer

	listenAddr string
	startTime  time.Time

	mut         *sync.Mutex
	minionPools map[string]*grpcpool.Pool
}

func NewBossServer(listenAddr string) *BossServer {
	return &BossServer{
		listenAddr:  listenAddr,
		startTime:   time.Now(),
		mut:         &sync.Mutex{},
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
		UptimeInfo: &proto.UpTimeInfo{
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

func (s *BossServer) RunWorkload(ctx context.Context, in *proto.BossRunWorkloadRequest) (*proto.BossRunWorkloadResponse, error) {
	log.Infof("GRPC: RunWorkload()")

	// Partition the data among the number of minions
	dataSpec := proto.DataSpecFromProto(in.DataSpec)
	numRecs := dataSpec.KeyGenSpec.NumKeys()
	fullRange := intgen.NewRange(0, numRecs)
	numMinions := len(s.minionPools)
	var assignedRanges []*intgen.Range

	workloadName := in.GetWlSpec().GetWorkloadName()
	switch workloadName {
	case "insert", "insert-txn", "delete", "delete-txn":
		assignedRanges = fullRange.Partition(numMinions)
	case "select-pk", "select-pk-txn", "select-uk", "select-uk-txn", "update", "update-txn":
		assignedRanges = fullRange.Duplicate(numMinions)
	default:
		reason := fmt.Sprintf("invalid workload '%s'", workloadName)
		log.Infof("RunWorkload(): %s", reason)
		return &proto.BossRunWorkloadResponse{
			OverallStatus: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}, fmt.Errorf(reason)
	}
	log.Infof("RunWorkload(): Data partitioned among minions: %v", assignedRanges)

	// Run the workload
	log.Infof("RunWorkload(): Now starting to run workload")
	minionStatuses := make([]*proto.MinionStatus, len(s.minionPools))
	ch := make(chan *proto.MinionStatus)

	i := 0
	for addr := range s.minionPools {
		log.Infof("RunWorkload(): Triggering on Minion %s", addr)
		wlSpec := &proto.WorkloadSpec{
			WorkloadName:  in.GetWlSpec().GetWorkloadName(),
			AssignedRange: proto.RangeToProto(assignedRanges[i]),
			TableName:     in.GetWlSpec().GetTableName(),
			DurationSec:   in.GetWlSpec().GetDurationSec(),
			Concurrency:   in.GetWlSpec().GetConcurrency(),
			BatchSize:     in.GetWlSpec().GetBatchSize(),
		}

		go s.runWorkloadOnMinion(ctx, addr, in.GetDataSpec(), in.DbSpec, wlSpec, ch)
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

	log.Infof("RunWorkload(): completed")
	return &proto.BossRunWorkloadResponse{
		OverallStatus:  &overallStatus,
		MinionStatuses: minionStatuses,
	}, nil
}

func (s *BossServer) StopWorkload(ctx context.Context, _ *proto.BossStopWorkloadRequest) (*proto.BossStopWorkloadResponse, error) {
	log.Infof("GRPC: StopWorkload()")

	minionStatuses := make([]*proto.MinionStatus, len(s.minionPools))
	ch := make(chan *proto.MinionStatus)

	for addr := range s.minionPools {
		log.Infof("StopWorkload(): Stopping on Minion %s", addr)

		go s.stopWorkloadOnMinion(ctx, addr, ch)
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

	log.Infof("StopWorkload(): completed")
	return &proto.BossStopWorkloadResponse{
		OverallStatus:  &overallStatus,
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
		BuildInfo:  res.GetBuildInfo(),
		UptimeInfo: res.GetUptimeInfo(),
	}
}

func (s *BossServer) loadDataSpecOnMinion(ctx context.Context, dataSpec *proto.DataSpec, minionAddr string, minionClient proto.MinionClient) error {
	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	loadDataSpecResponse, err := minionClient.LoadDataSpec(grpcCtx, &proto.MinionLoadDataSpecRequest{
		DataSpec: dataSpec,
	})
	grpcCtxCancel()

	if err != nil {
		return fmt.Errorf("LoadDataSpec request to %s failed (%v)", minionAddr, err)
	}

	if loadDataSpecResponse.GetStatus().GetIsOk() != true {
		return fmt.Errorf("LoadDataSpec request to %s failed (%s)", minionAddr, loadDataSpecResponse.GetStatus().GetFailureReason())
	}

	log.Infof("RunWorkload(): LoadDataSpec request on %s successful", minionAddr)
	return nil
}

func (s *BossServer) openDbConnectionOnMinion(ctx context.Context, dbSpec *proto.DBSpec, minionAddr string, minionClient proto.MinionClient) error {
	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	openDbConnResponse, err := minionClient.OpenDBConnection(grpcCtx, &proto.MinionOpenDBConnectionRequest{
		DbSpec: dbSpec,
	})
	grpcCtxCancel()

	if err != nil {
		return fmt.Errorf("OpenDBConnection request to %s failed (%v)", minionAddr, err)
	}
	if openDbConnResponse.GetStatus().GetIsOk() != true {
		return fmt.Errorf("OpenDBConnection request to %s failed (%s)", minionAddr, openDbConnResponse.GetStatus().GetFailureReason())
	}
	log.Infof("RunWorkload(): OpenDBConnection on %s successful", minionAddr)
	return nil
}

func (s *BossServer) startWorkloadOnMinion(ctx context.Context, wlSpec *proto.WorkloadSpec, minionAddr string, minionClient proto.MinionClient) error {
	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	openDbConnResponse, err := minionClient.RunWorkload(grpcCtx, &proto.MinionRunWorkloadRequest{
		WorkloadSpec: wlSpec,
	})
	grpcCtxCancel()

	if err != nil {
		return fmt.Errorf("RunWorkload request to %s failed (%v)", minionAddr, err)
	}
	if openDbConnResponse.GetStatus().GetIsOk() != true {
		return fmt.Errorf("RunWorkload request to %s failed (%s)", minionAddr, openDbConnResponse.GetStatus().GetFailureReason())
	}
	log.Infof("RunWorkload(): RunWorkload on %s successful", minionAddr)
	return nil
}

func (s *BossServer) runWorkloadOnMinion(ctx context.Context, addr string,
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

	err = s.loadDataSpecOnMinion(ctx, dataSpec, addr, minionClient)
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

	err = s.openDbConnectionOnMinion(ctx, dbSpec, addr, minionClient)
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

	err = s.startWorkloadOnMinion(ctx, wlSpec, addr, minionClient)
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

	ch <- &proto.MinionStatus{
		Addr: addr,
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
	}
	return
}

func (s *BossServer) stopWorkloadOnMinion(ctx context.Context, addr string, ch chan *proto.MinionStatus) {
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
	stopWorkloadResponse, err := minionClient.StopWorkload(grpcCtx, &proto.MinionStopWorkloadRequest{})
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

	ch <- &proto.MinionStatus{
		Addr: addr,
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
	}
	return
}
