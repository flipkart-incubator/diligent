package main

import (
	"context"
	"fmt"
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

	mut         *sync.Mutex
	minionPools map[string]*grpcpool.Pool
}

func NewBossServer(listenAddr string) *BossServer {
	return &BossServer{
		listenAddr:  listenAddr,
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

	return &proto.BossPingResponse{}, nil
}

func (s *BossServer) RegisterMinion(_ context.Context, in *proto.BossRegisterMinionRequest) (*proto.BossRegisterMinionResponse, error) {
	url := in.GetUrl()
	log.Infof("GRPC: RegisterMinion(%s)", url)
	s.mut.Lock()
	defer s.mut.Unlock()

	// Acknowledge the request by making an entry for the minion, but first close any existing pool
	pool, ok := s.minionPools[url]
	if ok && pool != nil {
		log.Infof("GRPC: RegisterMinion(%s): Existing entry found. Closing existing pool", url)
		pool.Close()
	}
	s.minionPools[url] = nil

	// Factory method for pool
	var factory grpcpool.Factory = func() (*grpc.ClientConn, error) {
		log.Infof("grpcpool.Factory(): Trying to connect to minion %s", url)
		dialCtx, dialCancel := context.WithTimeout(context.Background(), minionDialTimeoutSecs*time.Second)
		defer dialCancel()
		conn, err := grpc.DialContext(dialCtx, url, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Errorf("grpcpool.Factory(): Failed to connect to minion %s (%s)", url, err.Error())
			return nil, err
		}
		log.Infof("grpcpool.Factory(): Successfully connected to %s", url)
		return conn, nil
	}

	// Create an empty connection pool (don't try to establish connection right now as minion may not be ready)
	pool, err := grpcpool.New(factory, 0, 1, minionConnIdleTimeoutSecs*time.Second, minionConnMaxLifetimeSecs*time.Second)
	if err != nil {
		log.Errorf("GRPC: RegisterMinion(%s): Failed to create pool (%s)", url, err.Error())
		return &proto.BossRegisterMinionResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}, err
	}

	log.Infof("GRPC: RegisterMinion(%s): Successfully registered", url)
	s.minionPools[url] = pool
	return &proto.BossRegisterMinionResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
	}, nil
}

func (s *BossServer) UnregisterMinion(_ context.Context, in *proto.BossUnregisterMinonRequest) (*proto.BossUnregisterMinionResponse, error) {
	url := in.GetUrl()
	log.Infof("GRPC: UnregisterMinion(%s)", url)
	s.mut.Lock()
	defer s.mut.Unlock()

	pool, ok := s.minionPools[url]
	if ok {
		log.Infof("GRPC: UnregisterMinion(%s): Found registered minion - removing", url)
		delete(s.minionPools, url)
		if pool != nil {
			log.Infof("GRPC: UnregisterMinion(%s): Found connection pool - closing", url)
			pool.Close()
		} else {
			log.Infof("GRPC: UnregisterMinion(%s): Pool was nil, ignoring close", url)
		}
	} else {
		log.Infof("GRPC: UnregisterMinion(%s): No such minion. Ignoring request", url)
	}

	log.Infof("GRPC: Minion unregistered: %s", url)
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

	minionStatuses := make([]*proto.MinionStatus, len(s.minionPools))
	ch := make(chan *proto.MinionStatus)

	for url := range s.minionPools {
		go s.getMinionStatus(ctx, url, ch)
	}

	for i, _ := range minionStatuses {
		minionStatuses[i] = <-ch
	}

	return &proto.BossShowMinionResponse{
		Minions: minionStatuses,
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
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}, fmt.Errorf(reason)
	}
	log.Infof("RunWorkload(): Data partitioned among minions: %v", assignedRanges)

	// Run the workload
	log.Infof("RunWorkload(): Now starting to run workload")
	i := 0
	for url := range s.minionPools {
		log.Infof("RunWorkload(): Triggering on Minion %s", url)
		wlSpec := &proto.WorkloadSpec{
			WorkloadName:  in.GetWlSpec().GetWorkloadName(),
			AssignedRange: proto.RangeToProto(assignedRanges[i]),
			TableName:     in.GetWlSpec().GetTableName(),
			DurationSec:   in.GetWlSpec().GetDurationSec(),
			Concurrency:   in.GetWlSpec().GetConcurrency(),
			BatchSize:     in.GetWlSpec().GetBatchSize(),
		}

		err := s.runWorkloadOnMinion(ctx, url, in.GetDataSpec(), in.DbSpec, wlSpec)
		i++

		if err != nil {
			reason := fmt.Sprintf("Failed to run workload on %s (%s)", url, err.Error())
			log.Errorf("RunWorkload(): %s", reason)
			return &proto.BossRunWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}

		log.Infof("RunWorkload(): Run workload request to %s successful", url)
	}

	log.Infof("RunWorkload(): Run started successfully")
	return &proto.BossRunWorkloadResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
	}, nil
}

func (s *BossServer) StopWorkload(ctx context.Context, _ *proto.BossStopWorkloadRequest) (*proto.BossStopWorkloadResponse, error) {
	log.Infof("GRPC: StopWorkload()")

	for url := range s.minionPools {
		err := s.stopWorkloadOnMinion(ctx, url)
		if err != nil {
			return &proto.BossStopWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: err.Error(),
				},
			}, err
		}
	}

	log.Infof("StopWorkload(): Stopped successfully")
	return &proto.BossStopWorkloadResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
	}, nil
}

func (s *BossServer) getConnectionForMinion(ctx context.Context, minionUrl string) (*grpcpool.ClientConn, error) {
	pool, ok := s.minionPools[minionUrl]

	if !ok {
		return nil, fmt.Errorf("no such minion: %s", minionUrl)
	}

	if pool == nil {
		return nil, fmt.Errorf("no pool created for minion: %s", minionUrl)
	}

	conn, err := pool.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection for minion: %s (%s)", minionUrl, err.Error())
	}

	return conn, nil
}

func (s *BossServer) getMinionStatus(ctx context.Context, minionUrl string, ch chan *proto.MinionStatus) {
	minionConnection, err := s.getConnectionForMinion(ctx, minionUrl)
	if err != nil {
		reason := fmt.Sprintf("failed to get client for minion %s (%s)", minionUrl, err.Error())
		log.Errorf("getMinionStatus(): %s", reason)
		ch <- &proto.MinionStatus{
			Url: minionUrl,
			Status: &proto.GeneralStatus{
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
	_, err = minionClient.Ping(pingCtx, &proto.MinionPingRequest{})
	if err != nil {
		reason := fmt.Sprintf("gRPC request to %s failed (%v)", minionUrl, err)
		log.Errorf("getMinionStatus(): %s", reason)
		ch <- &proto.MinionStatus{
			Url: minionUrl,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}
		return
	}

	log.Infof("getMinionStatus(): Ping request to %s successful", minionUrl)
	ch <- &proto.MinionStatus{
		Url: minionUrl,
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
	}
}

func (s *BossServer) loadDataSpecOnMinion(ctx context.Context, dataSpec *proto.DataSpec, minionUrl string, minionClient proto.MinionClient) error {
	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	loadDataSpecResponse, err := minionClient.LoadDataSpec(grpcCtx, &proto.MinionLoadDataSpecRequest{
		DataSpec: dataSpec,
	})
	grpcCtxCancel()

	if err != nil {
		return fmt.Errorf("LoadDataSpec request to %s failed (%v)", minionUrl, err)
	}

	if loadDataSpecResponse.GetStatus().GetIsOk() != true {
		return fmt.Errorf("LoadDataSpec request to %s failed (%s)", minionUrl, loadDataSpecResponse.GetStatus().GetFailureReason())
	}

	log.Infof("RunWorkload(): LoadDataSpec request on %s successful", minionUrl)
	return nil
}

func (s *BossServer) openDbConnectionOnMinion(ctx context.Context, dbSpec *proto.DBSpec, minionUrl string, minionClient proto.MinionClient) error {
	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	openDbConnResponse, err := minionClient.OpenDBConnection(grpcCtx, &proto.MinionOpenDBConnectionRequest{
		DbSpec: dbSpec,
	})
	grpcCtxCancel()

	if err != nil {
		return fmt.Errorf("OpenDBConnection request to %s failed (%v)", minionUrl, err)
	}
	if openDbConnResponse.GetStatus().GetIsOk() != true {
		return fmt.Errorf("OpenDBConnection request to %s failed (%s)", minionUrl, openDbConnResponse.GetStatus().GetFailureReason())
	}
	log.Infof("RunWorkload(): OpenDBConnection on %s successful", minionUrl)
	return nil
}

func (s *BossServer) startWorkloadOnMinion(ctx context.Context, wlSpec *proto.WorkloadSpec, minionUrl string, minionClient proto.MinionClient) error {
	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	openDbConnResponse, err := minionClient.RunWorkload(grpcCtx, &proto.MinionRunWorkloadRequest{
		WorkloadSpec: wlSpec,
	})
	grpcCtxCancel()

	if err != nil {
		return fmt.Errorf("RunWorkload request to %s failed (%v)", minionUrl, err)
	}
	if openDbConnResponse.GetStatus().GetIsOk() != true {
		return fmt.Errorf("RunWorkload request to %s failed (%s)", minionUrl, openDbConnResponse.GetStatus().GetFailureReason())
	}
	log.Infof("RunWorkload(): RunWorkload on %s successful", minionUrl)
	return nil
}

func (s *BossServer) runWorkloadOnMinion(ctx context.Context, url string,
	dataSpec *proto.DataSpec, dbSpec *proto.DBSpec, wlSpec *proto.WorkloadSpec) error {

	minionConnection, err := s.getConnectionForMinion(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to get client for minion %s (%s)", url, err.Error())
	}
	defer minionConnection.Close()

	minionClient := proto.NewMinionClient(minionConnection)

	err = s.loadDataSpecOnMinion(ctx, dataSpec, url, minionClient)
	if err != nil {
		return err
	}

	err = s.openDbConnectionOnMinion(ctx, dbSpec, url, minionClient)
	if err != nil {
		return err
	}

	err = s.startWorkloadOnMinion(ctx, wlSpec, url, minionClient)
	if err != nil {
		return err
	}

	return nil
}

func (s *BossServer) stopWorkloadOnMinion(ctx context.Context, url string) error {
	minionConnection, err := s.getConnectionForMinion(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to get client for minion %s (%s)", url, err.Error())
	}
	defer minionConnection.Close()

	minionClient := proto.NewMinionClient(minionConnection)

	// Stop workload
	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	stopWorkloadResponse, err := minionClient.StopWorkload(grpcCtx, &proto.MinionStopWorkloadRequest{})
	grpcCtxCancel()

	if err != nil {
		return fmt.Errorf("request to %s failed (%v)", url, err)
	}

	if stopWorkloadResponse.GetStatus().GetIsOk() != true {
		return fmt.Errorf("request to %s failed (%s)", url, stopWorkloadResponse.GetStatus().GetFailureReason())
	}

	return nil
}
