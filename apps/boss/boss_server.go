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
	minionMaxLifetimeSecs = 600
)

type BossServer struct {
	proto.UnimplementedBossServer
	host     string
	grpcPort string

	mut         *sync.Mutex
	minionPools map[string]*grpcpool.Pool
}

func NewBossServer(host, grpcPort string) *BossServer {
	return &BossServer{
		host:        host,
		grpcPort:    grpcPort,
		mut:         &sync.Mutex{},
		minionPools: make(map[string]*grpcpool.Pool),
	}
}

func (s *BossServer) Serve() error {
	// Create listening port
	address := s.host + s.grpcPort
	log.Infof("Creating listening port: %s", address)
	listener, err := net.Listen("tcp", address)
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

func (s *BossServer) Ping(context.Context, *proto.BossPingRequest) (*proto.BossPingResponse, error) {
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
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		conn, err := grpc.DialContext(ctx, url, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Errorf("grpcpool.Factory(): Failed to connect to minion %s (%s)", url, err.Error())
			return nil, err
		}
		log.Infof("grpcpool.Factory(): Successfully connected to %s", url)
		return conn, nil
	}

	// Create an empty connection pool (don't try to establish connection right now as minion may not be ready)
	pool, err := grpcpool.New(factory, 0, 1, minionConnIdleTimeoutSecs * time.Second, minionMaxLifetimeSecs * time.Second)
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

func (s *BossServer) ShowMinions(context.Context, *proto.BossShowMinionRequest) (*proto.BossShowMinionResponse, error) {
	log.Infof("GRPC: ShowMinions()")
	s.mut.Lock()
	defer s.mut.Unlock()

	minionStatuses := make([]*proto.MinionStatus, 0)

	for url := range s.minionPools {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		minionClient, err := s.getClientForMinion(ctx, url)
		if err != nil {
			reason := fmt.Sprintf("failed to get client for minion %s (%s)", url, err.Error())
			log.Errorf("ShowMinions(): %s", reason)
			minionStatus := &proto.MinionStatus{
				Url: url,
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}
			minionStatuses = append(minionStatuses, minionStatus)
			continue
		}

		// Ping Minion
		ctx, _ = context.WithTimeout(context.Background(), 1*time.Second)
		_, err = minionClient.Ping(ctx, &proto.MinionPingRequest{})
		if err != nil {
			reason := fmt.Sprintf("gRPC request to %s failed (%v)", url, err)
			log.Errorf("ShowMinions(): %s", reason)
			minionStatus := &proto.MinionStatus{
				Url: url,
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}
			minionStatuses = append(minionStatuses, minionStatus)
			continue
		}

		log.Infof("ShowMinions(): Ping request to %s successful", url)
		minionStatuses = append(minionStatuses, &proto.MinionStatus{
			Url: url,
			Status: &proto.GeneralStatus{
				IsOk:          true,
				FailureReason: "",
			},
		})
	}

	return &proto.BossShowMinionResponse{
		Minions: minionStatuses,
	}, nil
}

func (s *BossServer) RunWorkload(_ context.Context, in *proto.BossRunWorkloadRequest) (*proto.BossRunWorkloadResponse, error) {
	log.Infof("GRPC: RunWorkload()")

	for url := range s.minionPools {
		log.Infof("RunWorkload(): Preparing Minion %s", url)

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		minionClient, err := s.getClientForMinion(ctx, url)
		if err != nil {
			reason := fmt.Sprintf("failed to get client for minion %s (%s)", url, err.Error())
			log.Errorf("RunWorkload(): %s", reason)
			return &proto.BossRunWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}

		// Load Data Spec
		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
		loadDataSpecResponse, err := minionClient.LoadDataSpec(ctx, &proto.MinionLoadDataSpecRequest{
			DataSpec: in.GetDataSpec(),
		})
		if err != nil {
			reason := fmt.Sprintf("request to %s failed (%v)", url, err)
			log.Errorf("RunWorkload(): %s", reason)
			return &proto.BossRunWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}
		if loadDataSpecResponse.GetStatus().GetIsOk() != true {
			reason := fmt.Sprintf("request to %s failed (%s)", url, loadDataSpecResponse.GetStatus().GetFailureReason())
			log.Errorf("RunWorkload(): %s", reason)
			return &proto.BossRunWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}
		log.Infof("RunWorkload(): Load dataspec on %s successful", url)

		// Connect to DB
		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
		openDbResponse, err := minionClient.OpenDBConnection(ctx, &proto.MinionOpenDBConnectionRequest{
			DbSpec: in.GetDbSpec(),
		})
		if err != nil {
			reason := fmt.Sprintf("request to %s failed (%v)", url, err)
			log.Errorf("RunWorkload(): %s", reason)
			return &proto.BossRunWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}
		if openDbResponse.GetStatus().GetIsOk() != true {
			reason := fmt.Sprintf("request to %s failed (%s)", url, openDbResponse.GetStatus().GetFailureReason())
			log.Errorf("RunWorkload(): %s", reason)
			return &proto.BossRunWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}
		log.Infof("RunWorkload(): Connect to DB on %s successful", url)
	}

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

	log.Infof("RunWorkload(): Now starting to run workload")
	// Run workload
	i := 0
	for url := range s.minionPools {
		log.Infof("RunWorkload(): Triggering on Minion %s", url)

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		minionClient, err := s.getClientForMinion(ctx, url)
		if err != nil {
			reason := fmt.Sprintf("failed to get client for minion %s (%s)", url, err.Error())
			log.Errorf("RunWorkload(): %s", reason)
			return &proto.BossRunWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}

		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
		runWorkloadResponse, err := minionClient.RunWorkload(ctx, &proto.MinionRunWorkloadRequest{
			WorkloadSpec: &proto.WorkloadSpec{
				WorkloadName:  in.GetWlSpec().GetWorkloadName(),
				AssignedRange: proto.RangeToProto(assignedRanges[i]),
				TableName:     in.GetWlSpec().GetTableName(),
				DurationSec:   in.GetWlSpec().GetDurationSec(),
				Concurrency:   in.GetWlSpec().GetConcurrency(),
				BatchSize:     in.GetWlSpec().GetBatchSize(),
			},
		})
		i++
		if err != nil {
			reason := fmt.Sprintf("request to %s failed (%v)", url, err)
			log.Errorf("RunWorkload(): %s", reason)
			return &proto.BossRunWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}
		if runWorkloadResponse.GetStatus().GetIsOk() != true {
			reason := fmt.Sprintf("request to %s failed (%s)", url, runWorkloadResponse.GetStatus().GetFailureReason())
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
		minionClient, err := s.getClientForMinion(ctx, url)
		if err != nil {
			reason := fmt.Sprintf("failed to get client for minion %s (%s)", url, err.Error())
			log.Errorf("StopWorkload(): %s", reason)
			return &proto.BossStopWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}

		// Stop workload
		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
		stopWorkloadResponse, err := minionClient.StopWorkload(ctx, &proto.MinionStopWorkloadRequest{})
		if err != nil {
			reason := fmt.Sprintf("request to %s failed (%v)", url, err)
			log.Errorf("StopWorkload(): %s", reason)
			return &proto.BossStopWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}
		if stopWorkloadResponse.GetStatus().GetIsOk() != true {
			reason := fmt.Sprintf("request to %s failed (%s)", url, stopWorkloadResponse.GetStatus().GetFailureReason())
			log.Errorf("StopWorkload(): %s", reason)
			return &proto.BossStopWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
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

func (s *BossServer) getClientForMinion(ctx context.Context, minionUrl string) (proto.MinionClient, error) {
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

	return proto.NewMinionClient(conn), nil
}
