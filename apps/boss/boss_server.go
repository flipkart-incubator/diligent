package main

import (
	"context"
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/intgen"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

type BossServer struct {
	proto.UnimplementedBossServer
	host     string
	grpcPort string

	mut        *sync.Mutex
	minionUrls map[string]bool
}

func NewBossServer(host, grpcPort string) *BossServer {
	return &BossServer{
		host:       host,
		grpcPort:   grpcPort,
		mut:        &sync.Mutex{},
		minionUrls: make(map[string]bool),
	}
}

func (s *BossServer) Serve() error {
	// Create listening port
	address := s.host + s.grpcPort
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register our this server object with the gRPC server
	proto.RegisterBossServer(grpcServer, s)

	// Start serving
	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

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
	log.Infof("GRPC: MinionRegister(%s)", url)
	s.mut.Lock()
	defer s.mut.Unlock()
	s.minionUrls[url] = true
	return &proto.BossRegisterMinionResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
	}, nil
}

func (s *BossServer) UnregisterMinion(_ context.Context, in *proto.BossUnregisterMinonRequest) (*proto.BossUnregisterMinionResponse, error) {
	url := in.GetUrl()
	log.Infof("GRPC: MinionUnregister(%s)", url)
	s.mut.Lock()
	defer s.mut.Unlock()
	delete(s.minionUrls, url)
	return &proto.BossUnregisterMinionResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
	}, nil
}

func (s *BossServer) ShowMinions(context.Context, *proto.BossShowMinionRequest) (*proto.BossShowMinionResponse, error) {
	log.Infof("GRPC: MinionShow()")
	s.mut.Lock()
	defer s.mut.Unlock()

	minionStatuses := make([]*proto.MinionStatus, 0)

	for url := range s.minionUrls {
		log.Infof("RunWorkload(): Trying ping for minion %s", url)

		// Connect to Minion
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		conn, err := grpc.DialContext(ctx, url, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			reason := fmt.Sprintf("connection to %s failed (%v)", url, err)
			log.Infof("RunWorkload(): TCP connection to %s failed (%s)", url, reason)
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

		log.Infof("RunWorkload(): TCP connection to %s successful", url)
		minionClient := proto.NewMinionClient(conn)

		// Ping Minion
		ctx, _ = context.WithTimeout(context.Background(), 1*time.Second)
		_, err = minionClient.Ping(ctx, &proto.MinionPingRequest{})
		if err != nil {
			reason := fmt.Sprintf("gRPC request to %s failed (%v)", url, err)
			log.Infof("RunWorkload(): gRPC request to %s failed (%s)", url, reason)
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

		log.Infof("RunWorkload(): gRPC request to %s successful", url)

		minionStatuses = append(minionStatuses, &proto.MinionStatus{
			Url:         url,
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

	for url := range s.minionUrls {
		log.Infof("RunWorkload(): Minion %s", url)

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		conn, err := grpc.DialContext(ctx, url, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			// TODO: Simple impl for now
			reason := fmt.Sprintf("failed to connect to %s (%v)", url, err)
			return &proto.BossRunWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}
		log.Infof("RunWorkload(): TCP connection to %s successful", url)
		minionClient := proto.NewMinionClient(conn)

		// Load Data Spec
		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
		loadDataSpecResponse, err := minionClient.LoadDataSpec(ctx, &proto.MinionLoadDataSpecRequest{
			DataSpec: in.GetDataSpec(),
		})
		if err != nil {
			reason := fmt.Sprintf("request to %s failed (%v)", url, err)
			return &proto.BossRunWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}
		if loadDataSpecResponse.GetStatus().GetIsOk() != true {
			reason := fmt.Sprintf("request to %s failed (%s)", url, loadDataSpecResponse.GetStatus().GetFailureReason())
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
			return &proto.BossRunWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}
		if openDbResponse.GetStatus().GetIsOk() != true {
			reason := fmt.Sprintf("request to %s failed (%s)", url, openDbResponse.GetStatus().GetFailureReason())
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
	numMinions := len(s.minionUrls)
	var assignedRanges []*intgen.Range
	workloadName := in.GetWlSpec().GetWorkloadName()
	switch workloadName {
	case "insert", "insert-txn", "delete", "delete-txn":
		assignedRanges = fullRange.Partition(numMinions)
	case "select-pk", "select-pk-txn", "select-uk", "select-uk-txn", "update", "update-txn":
		assignedRanges = fullRange.Duplicate(numMinions)
	default:
		reason := fmt.Sprintf("invalid workload '%s'", workloadName)
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
	for url := range s.minionUrls {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		conn, err := grpc.DialContext(ctx, url, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			// TODO: Simple impl for now
			reason := fmt.Sprintf("failed to connect to %s (%v)", url, err)
			return &proto.BossRunWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}
		log.Infof("RunWorkload(): TCP connection to %s successful", url)

		minionClient := proto.NewMinionClient(conn)

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
			return &proto.BossRunWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}
		if runWorkloadResponse.GetStatus().GetIsOk() != true {
			reason := fmt.Sprintf("request to %s failed (%s)", url, runWorkloadResponse.GetStatus().GetFailureReason())
			return &proto.BossRunWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}
		log.Infof("RunWorkload(): Run workload request to %s successful", url)
	}

	return &proto.BossRunWorkloadResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
	}, nil
}

func (s *BossServer) StopWorkload(context.Context, *proto.BossStopWorkloadRequest) (*proto.BossStopWorkloadResponse, error) {
	for url := range s.minionUrls {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		conn, err := grpc.DialContext(ctx, url, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			// TODO: Simple impl for now
			reason := fmt.Sprintf("failed to connect to %s (%v)", url, err)
			return &proto.BossStopWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}
		minionClient := proto.NewMinionClient(conn)

		// Stop workload
		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
		stopWorkloadResponse, err := minionClient.StopWorkload(ctx, &proto.MinionStopWorkloadRequest{})
		if err != nil {
			reason := fmt.Sprintf("request to %s failed (%v)", url, err)
			return &proto.BossStopWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}
		if stopWorkloadResponse.GetStatus().GetIsOk() != true {
			reason := fmt.Sprintf("request to %s failed (%s)", url, stopWorkloadResponse.GetStatus().GetFailureReason())
			return &proto.BossStopWorkloadResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: reason,
				},
			}, fmt.Errorf(reason)
		}
	}

	return &proto.BossStopWorkloadResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
	}, nil
}
