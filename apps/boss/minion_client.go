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

// MinionClient helps us run gRPC operations on a minion
// It internally maintains the connection pool for the minion
// It is thread safe
type MinionClient struct {
	mut  sync.Mutex // Must be taken for all operations on the MinionClient
	addr string
	pool *grpcpool.Pool
}

func NewMinionClient(addr string) (*MinionClient, error) {
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

	return &MinionClient{
		addr: addr,
		pool: pool,
	}, nil
}

func (m *MinionClient) Close() {
	m.mut.Lock()
	defer m.mut.Unlock()

	if m.pool != nil {
		m.pool.Close()
		m.pool = nil
	}
}

func (m *MinionClient) GetAddr() string {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.addr
}

func (m *MinionClient) Ping(ctx context.Context) (*proto.MinionPingResponse, error) {
	m.mut.Lock()
	defer m.mut.Unlock()

	minionConnection, err := m.getConnectionForMinion(ctx)
	if err != nil {
		log.Errorf("Ping(): %s", err.Error())
		return nil, err
	}
	defer minionConnection.Close()
	grpcClient := proto.NewMinionClient(minionConnection)

	// Ping Minion
	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	defer grpcCtxCancel()
	res, err := grpcClient.Ping(grpcCtx, &proto.MinionPingRequest{})
	if err != nil {
		e := fmt.Errorf("gRPC request to %s failed (%v)", m.addr, err)
		log.Errorf("getMinionStatus(): %s", e.Error())
		return nil, e
	}

	log.Infof("getMinionStatus(): Ping request to %s successful", m.addr)
	return res, nil
}

func (m *MinionClient) PrepareJob(ctx context.Context, jobId, jobDesc string,
	dataSpec *proto.DataSpec, dbSpec *proto.DBSpec, wlSpec *proto.WorkloadSpec) (*proto.MinionPrepareJobResponse, error) {
	m.mut.Lock()
	defer m.mut.Unlock()

	minionConnection, err := m.getConnectionForMinion(ctx)
	if err != nil {
		log.Errorf("PrepareJob(): %s", err.Error())
		return nil, err
	}
	defer minionConnection.Close()
	grpcClient := proto.NewMinionClient(minionConnection)

	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	res, err := grpcClient.PrepareJob(grpcCtx, &proto.MinionPrepareJobRequest{
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
		return nil, err
	}
	return res, nil
}

func (m *MinionClient) RunJob(ctx context.Context) (*proto.MinionRunJobResponse, error) {
	m.mut.Lock()
	defer m.mut.Unlock()

	minionConnection, err := m.getConnectionForMinion(ctx)
	if err != nil {
		log.Errorf("RunJob(): %s", err.Error())
		return nil, err
	}
	defer minionConnection.Close()
	grpcClient := proto.NewMinionClient(minionConnection)

	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	res, err := grpcClient.RunJob(grpcCtx, &proto.MinionRunJobRequest{})
	grpcCtxCancel()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *MinionClient) AbortJob(ctx context.Context) (*proto.MinionAbortJobResponse, error) {
	m.mut.Lock()
	defer m.mut.Unlock()

	minionConnection, err := m.getConnectionForMinion(ctx)
	if err != nil {
		log.Errorf("AbortJob(): %s", err.Error())
		return nil, err
	}
	defer minionConnection.Close()
	grpcClient := proto.NewMinionClient(minionConnection)

	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	res, err := grpcClient.AbortJob(grpcCtx, &proto.MinionAbortJobRequest{})
	grpcCtxCancel()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *MinionClient) getConnectionForMinion(ctx context.Context) (*grpcpool.ClientConn, error) {
	if m.pool == nil {
		return nil, fmt.Errorf("no pool for minion: %s", m.addr)
	}

	conn, err := m.pool.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection for minion: %s (%s)", m.addr, err.Error())
	}

	return conn, nil
}
