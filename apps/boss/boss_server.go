package main

import (
	"context"
	"github.com/flipkart-incubator/diligent/pkg/buildinfo"
	"github.com/flipkart-incubator/diligent/pkg/idgen"
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

	registry   *MinionRegistry
	job        *Job
	experiment *Experiment
}

func NewBossServer(listenAddr string) *BossServer {
	return &BossServer{
		listenAddr: listenAddr,
		pid:        idgen.GenerateId16(),
		startTime:  time.Now(),
		registry:   NewMinionRegistry(),
		job:        nil,
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
			StartTime: s.startTime.UnixMilli(),
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

func (s *BossServer) GetMinionInfo(ctx context.Context, _ *proto.BossGetMinionInfoRequest) (*proto.BossGetMinionInfoResponse, error) {
	log.Infof("GRPC: GetMinionInfo()")
	s.mut.Lock()
	defer s.mut.Unlock()

	// Slices to gather data on individual minions
	numMinions := s.registry.GetNumMinions()
	addrs := make([]string, numMinions)
	rchs := make([]chan *proto.MinionPingResponse, numMinions)
	echs := make([]chan error, numMinions)

	// Invoke on individual minions
	i := 0
	for addr, proxy := range s.registry.Minions() {
		log.Infof("GetMinionInfo(): Getting info from Minion %s", addr)
		addrs[i] = addr
		rchs[i], echs[i] = proxy.PingAsync(ctx)
		i++
	}

	// For holding info of individual minions
	minionInfos := make([]*proto.MinionInfo, numMinions)

	// Collect execution results
	for i, _ = range minionInfos {
		select {
		case res := <-rchs[i]:
			minionInfos[i] = &proto.MinionInfo{
				Addr: addrs[i],
				Reachability: &proto.GeneralStatus{
					IsOk:          true,
					FailureReason: "",
				},
				BuildInfo:   res.GetBuildInfo(),
				ProcessInfo: res.GetProcessInfo(),
				JobInfo:     res.GetJobInfo(),
			}
		case err := <-echs[i]:
			minionInfos[i] = &proto.MinionInfo{
				Addr: addrs[i],
				Reachability: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: err.Error(),
				},
			}
		}
	}

	log.Infof("GetMinionInfo(): completed")
	return &proto.BossGetMinionInfoResponse{
		MinionInfos: minionInfos,
	}, nil
}

func (s *BossServer) PrepareJob(ctx context.Context, in *proto.BossPrepareJobRequest) (*proto.BossPrepareJobResponse, error) {
	log.Infof("GRPC: PrepareJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.job != nil {
		log.Infof("PrepareJob(): current job=%s, state=%s", s.job.Name(), s.job.State())
		if !s.job.HasEnded() {
			return &proto.BossPrepareJobResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: "current job has not ended",
				},
			}, nil
		}
	}

	s.job = NewJob(in.JobSpec, s.registry.Minions())
	return s.job.Prepare(ctx)
}

func (s *BossServer) RunJob(ctx context.Context, in *proto.BossRunJobRequest) (*proto.BossRunJobResponse, error) {
	log.Infof("GRPC: RunJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.job == nil {
		return &proto.BossRunJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "no current job",
			},
		}, nil
	}

	if s.job.HasEnded() {
		return &proto.BossRunJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "current job has already ended",
			},
		}, nil
	}

	return s.job.Run(ctx)
}

func (s *BossServer) AbortJob(ctx context.Context, in *proto.BossAbortJobRequest) (*proto.BossAbortJobResponse, error) {
	log.Infof("GRPC: AbortJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.job == nil {
		return &proto.BossAbortJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "no current job",
			},
		}, nil
	}

	if s.job.HasEnded() {
		return &proto.BossAbortJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "current job has already ended",
			},
		}, nil
	}

	return s.job.Abort(ctx)
}

func (s *BossServer) GetJobInfo(ctx context.Context, in *proto.BossGetJobInfoRequest) (*proto.BossGetJobInfoResponse, error) {
	log.Infof("GRPC: GetJobInfo()")
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.job == nil {
		return &proto.BossGetJobInfoResponse{
			Status: &proto.GeneralStatus{
				IsOk: true,
			},
			JobInfo: nil,
		}, nil
	}

	log.Infof("GetJobInfo(): completed")
	return &proto.BossGetJobInfoResponse{
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
		JobInfo: s.job.Info().ToProto(),
	}, nil
}

func (s *BossServer) StartExperiment(ctx context.Context, in *proto.BossStartExperimentRequest) (*proto.BossStartExperimentResponse, error) {
	log.Infof("GRPC: StartExperiment()")
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.experiment != nil {
		log.Infof("StartExperiment(): current experiment=%s, state=%s", s.experiment.Name(), s.experiment.State())
		if !s.experiment.HasEnded() {
			return &proto.BossStartExperimentResponse{
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: "current experiment has not ended",
				},
			}, nil
		}
	}

	// TODO: Lock registry
	exp := NewExperiment(in.GetExperimentName())
	err := exp.Start()

	if err != nil {
		return &proto.BossStartExperimentResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}, nil
	}

	s.experiment = exp
	return &proto.BossStartExperimentResponse{
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
	}, nil
}

func (s *BossServer) StopExperiment(ctx context.Context, in *proto.BossStopExperimentRequest) (*proto.BossStopExperimentResponse, error) {
	log.Infof("GRPC: StopExperiment()")
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.experiment == nil {
		return &proto.BossStopExperimentResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "no current experiment",
			},
		}, nil
	}

	err := s.experiment.Stop()
	if err != nil {
		return &proto.BossStopExperimentResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}, nil
	}

	// TODO: Unlock registry
	return &proto.BossStopExperimentResponse{
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
	}, nil
}

func (s *BossServer) GetExperimentInfo(ctx context.Context, in *proto.BossGetExperimentInfoRequest) (*proto.BossGetExperimentInfoResponse, error) {
	log.Infof("GRPC: GetExperimentInfo()")
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.experiment == nil {
		return &proto.BossGetExperimentInfoResponse{
			Status: &proto.GeneralStatus{
				IsOk: true,
			},
			ExperimentInfo: nil,
		}, nil
	}

	log.Infof("GetExperimentInfo(): completed")
	return &proto.BossGetExperimentInfoResponse{
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
		ExperimentInfo: s.experiment.ToProto(),
	}, nil
}
