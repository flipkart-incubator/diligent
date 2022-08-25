package main

import (
	"context"
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	grpcpool "github.com/processout/grpc-go-pool"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sync"
	"time"
)

const (
	minionConnIdleTimeoutSecs = 600
	minionConnMaxLifetimeSecs = 600
	minionDialTimeoutSecs     = 3
	minionRequestTimeout      = 5
)

// MinionManager manages a single minion
// It takes care of the connection pool, and executing the gRPC operations on the minion it manages
// It is thread safe
type MinionManager struct {
	mut  sync.Mutex // Must be taken for all operations on the MinionManager
	addr string
	pool *grpcpool.Pool
}

func NewMinionManager(addr string) (*MinionManager, error) {
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
		log.Errorf("Failed to create pool for minion %s (%s)", addr, err.Error())
		return nil, err
	}

	return &MinionManager{
		addr: addr,
		pool: pool,
	}, nil
}

func (m *MinionManager) Close() {
	m.mut.Lock()
	defer m.mut.Unlock()

	if m.pool != nil {
		m.pool.Close()
		m.pool = nil
	}
}

func (m *MinionManager) GetAddr() string {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.addr
}

func (m *MinionManager) GetMinionStatus(ctx context.Context, ch chan *proto.MinionInfo) {
	minionConnection, err := m.getConnectionForMinion(ctx)
	if err != nil {
		reason := fmt.Sprintf("failed to get client for minion %s (%s)", m.addr, err.Error())
		log.Errorf("getMinionStatus(): %s", reason)
		ch <- &proto.MinionInfo{
			Addr: m.addr,
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
		reason := fmt.Sprintf("gRPC request to %s failed (%v)", m.addr, err)
		log.Errorf("getMinionStatus(): %s", reason)
		ch <- &proto.MinionInfo{
			Addr: m.addr,
			Reachability: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}
		return
	}

	log.Infof("getMinionStatus(): Ping request to %s successful", m.addr)
	ch <- &proto.MinionInfo{
		Addr: m.addr,
		Reachability: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
		BuildInfo:   res.GetBuildInfo(),
		ProcessInfo: res.GetProcessInfo(),
		JobInfo:     res.GetJobInfo(),
	}
}

func (m *MinionManager) PrepareJobOnMinion(ctx context.Context, jobId, jobDesc string,
	dataSpec *proto.DataSpec, dbSpec *proto.DBSpec, wlSpec *proto.WorkloadSpec, ch chan *proto.MinionStatus) {

	minionConnection, err := m.getConnectionForMinion(ctx)
	if err != nil {
		ch <- &proto.MinionStatus{
			Addr: m.addr,
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
			Addr: m.addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}
		return
	}

	if res.GetStatus().GetIsOk() != true {
		ch <- &proto.MinionStatus{
			Addr: m.addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: res.GetStatus().GetFailureReason(),
			},
		}
		return
	}

	log.Infof("Successfully prepared minion %s", m.addr)
	ch <- &proto.MinionStatus{
		Addr: m.addr,
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
	}
	return
}

func (m *MinionManager) RunJobOnMinion(ctx context.Context, ch chan *proto.MinionStatus) {
	minionConnection, err := m.getConnectionForMinion(ctx)
	if err != nil {
		ch <- &proto.MinionStatus{
			Addr: m.addr,
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
			Addr: m.addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}
		return
	}

	if res.GetStatus().GetIsOk() != true {
		ch <- &proto.MinionStatus{
			Addr: m.addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: res.GetStatus().GetFailureReason(),
			},
		}
		return
	}

	log.Infof("Successfully triggered run on minion %s", m.addr)
	ch <- &proto.MinionStatus{
		Addr: m.addr,
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
	}
	return
}

func (m *MinionManager) StopJobOnMinion(ctx context.Context, ch chan *proto.MinionStatus) {
	minionConnection, err := m.getConnectionForMinion(ctx)
	if err != nil {
		ch <- &proto.MinionStatus{
			Addr: m.addr,
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
			Addr: m.addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}
		return
	}

	if stopWorkloadResponse.GetStatus().GetIsOk() != true {
		ch <- &proto.MinionStatus{
			Addr: m.addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: stopWorkloadResponse.GetStatus().GetFailureReason(),
			},
		}
		return
	}

	log.Infof("Successfully triggered stop on minion %s", m.addr)
	ch <- &proto.MinionStatus{
		Addr: m.addr,
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
	}
	return
}

func (m *MinionManager) getConnectionForMinion(ctx context.Context) (*grpcpool.ClientConn, error) {
	if m.pool == nil {
		return nil, fmt.Errorf("no pool for minion: %s", m.addr)
	}

	conn, err := m.pool.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection for minion: %s (%s)", m.addr, err.Error())
	}

	return conn, nil
}
