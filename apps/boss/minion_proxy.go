package main

import (
	"context"
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

// MinionProxy helps us run gRPC operations on a minion
// It internally maintains the connection pool for the minion
// It is thread safe
type MinionProxy struct {
	mut  sync.Mutex // Must be taken for all operations on the MinionProxy
	addr string
	pool *grpcpool.Pool
}

func NewMinionProxy(addr string) (*MinionProxy, error) {
	// Factory method for pool
	var factory grpcpool.Factory = func() (*grpc.ClientConn, error) {
		log.Infof("grpcpool.Factory(): Trying to connect to minion %s", addr)
		dialCtx, dialCancel := context.WithTimeout(context.Background(), minionDialTimeoutSecs*time.Second)
		defer dialCancel()
		conn, err := grpc.DialContext(dialCtx, addr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Errorf("Failed to connect to minion %s (%s)", addr, err.Error())
			return nil, err
		}
		log.Infof("Successfully connected to %s", addr)
		return conn, nil
	}

	// Create an empty connection pool (don't try to establish connection right now as minion may not be ready)
	pool, err := grpcpool.New(factory, 0, 1, minionConnIdleTimeoutSecs*time.Second, minionConnMaxLifetimeSecs*time.Second)
	if err != nil {
		log.Errorf("Failed to create pool for minion %s (%s)", addr, err.Error())
		return nil, err
	}

	return &MinionProxy{
		addr: addr,
		pool: pool,
	}, nil
}

func (p *MinionProxy) Close() {
	p.mut.Lock()
	defer p.mut.Unlock()

	if p.pool != nil {
		p.pool.Close()
		p.pool = nil
	}
}

func (p *MinionProxy) GetAddr() string {
	p.mut.Lock()
	defer p.mut.Unlock()
	return p.addr
}

// PingAsync invokes Ping on a minion asynchronously
func (p *MinionProxy) PingAsync(ctx context.Context) (rch chan *proto.MinionPingResponse, ech chan error) {
	ech = make(chan error, 1)
	rch = make(chan *proto.MinionPingResponse, 1)
	go func() {
		res, err := p.PingSync(ctx)
		if err != nil {
			ech <- err
		} else {
			rch <- res
		}
	}()
	return rch, ech
}

// PingSync invokes Ping on a minion synchronously
func (p *MinionProxy) PingSync(ctx context.Context) (*proto.MinionPingResponse, error) {
	p.mut.Lock()
	defer p.mut.Unlock()

	conn, err := p.pool.Get(ctx)
	if err != nil {
		log.Errorf("PrepareJob(): %s", err.Error())
		return nil, err
	}
	defer conn.Close()
	grpcClient := proto.NewMinionClient(conn)

	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	res, err := grpcClient.Ping(grpcCtx, &proto.MinionPingRequest{})
	grpcCtxCancel()
	if err != nil {
		log.Errorf("Ping(): %s", err.Error())
		return nil, err
	}
	return res, nil
}

// PrepareJobAsync invokes prepare-job on a minion asynchronously
func (p *MinionProxy) PrepareJobAsync(ctx context.Context, jobName string,
	dataSpec *proto.DataSpec, dbSpec *proto.DBSpec, wlSpec *proto.WorkloadSpec) (rch chan *proto.MinionPrepareJobResponse, ech chan error) {
	ech = make(chan error, 1)
	rch = make(chan *proto.MinionPrepareJobResponse, 1)
	go func() {
		res, err := p.PrepareJobSync(ctx, jobName, dataSpec, dbSpec, wlSpec)
		if err != nil {
			ech <- err
		} else {
			rch <- res
		}
	}()
	return rch, ech
}

// PrepareJobSync invokes prepare-job on a minion synchronously
func (p *MinionProxy) PrepareJobSync(ctx context.Context, jobName string,
	dataSpec *proto.DataSpec, dbSpec *proto.DBSpec, wlSpec *proto.WorkloadSpec) (*proto.MinionPrepareJobResponse, error) {
	p.mut.Lock()
	defer p.mut.Unlock()

	conn, err := p.pool.Get(ctx)
	if err != nil {
		log.Errorf("PrepareJob(): %s", err.Error())
		return nil, err
	}
	defer conn.Close()
	grpcClient := proto.NewMinionClient(conn)

	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	res, err := grpcClient.PrepareJob(grpcCtx, &proto.MinionPrepareJobRequest{
		JobSpec: &proto.JobSpec{
			JobName:      jobName,
			DataSpec:     dataSpec,
			DbSpec:       dbSpec,
			WorkloadSpec: wlSpec,
		},
	})
	grpcCtxCancel()
	if err != nil {
		log.Errorf("PrepareJob(): %s", err.Error())
		return nil, err
	}
	return res, nil
}

// RunJobAsync invokes run-job on a minion asynchronously
func (p *MinionProxy) RunJobAsync(ctx context.Context) (rch chan *proto.MinionRunJobResponse, ech chan error) {
	ech = make(chan error, 1)
	rch = make(chan *proto.MinionRunJobResponse, 1)
	go func() {
		res, err := p.RunJobSync(ctx)
		if err != nil {
			ech <- err
		} else {
			rch <- res
		}
	}()
	return rch, ech
}

// RunJobSync invokes run-job on a minion synchronously
func (p *MinionProxy) RunJobSync(ctx context.Context) (*proto.MinionRunJobResponse, error) {
	p.mut.Lock()
	defer p.mut.Unlock()

	conn, err := p.pool.Get(ctx)
	if err != nil {
		log.Errorf("PrepareJob(): %s", err.Error())
		return nil, err
	}
	defer conn.Close()
	grpcClient := proto.NewMinionClient(conn)

	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	res, err := grpcClient.RunJob(grpcCtx, &proto.MinionRunJobRequest{})
	grpcCtxCancel()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// AbortJobAsync invokes abort-job on a minion asynchronously
func (p *MinionProxy) AbortJobAsync(ctx context.Context) (rch chan *proto.MinionAbortJobResponse, ech chan error) {
	ech = make(chan error, 1)
	rch = make(chan *proto.MinionAbortJobResponse, 1)
	go func() {
		res, err := p.AbortJobSync(ctx)
		if err != nil {
			ech <- err
		} else {
			rch <- res
		}
	}()
	return rch, ech
}

// AbortJobSync invokes abort-job on a minion synchronously
func (p *MinionProxy) AbortJobSync(ctx context.Context) (*proto.MinionAbortJobResponse, error) {
	p.mut.Lock()
	defer p.mut.Unlock()

	conn, err := p.pool.Get(ctx)
	if err != nil {
		log.Errorf("PrepareJob(): %s", err.Error())
		return nil, err
	}
	defer conn.Close()
	grpcClient := proto.NewMinionClient(conn)

	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	res, err := grpcClient.AbortJob(grpcCtx, &proto.MinionAbortJobRequest{})
	grpcCtxCancel()
	if err != nil {
		log.Errorf("AbortJob(): %s", err.Error())
		return nil, err
	}
	return res, nil
}

// QueryJobAsync invokes query-job on a minion asynchronously
func (p *MinionProxy) QueryJobAsync(ctx context.Context) (rch chan *proto.MinionQueryJobResponse, ech chan error) {
	ech = make(chan error, 1)
	rch = make(chan *proto.MinionQueryJobResponse, 1)
	go func() {
		res, err := p.QueryJobSync(ctx)
		if err != nil {
			ech <- err
		} else {
			rch <- res
		}
	}()
	return rch, ech
}

// QueryJobSync invokes query-job on a minion synchronously
func (p *MinionProxy) QueryJobSync(ctx context.Context) (*proto.MinionQueryJobResponse, error) {
	p.mut.Lock()
	defer p.mut.Unlock()

	conn, err := p.pool.Get(ctx)
	if err != nil {
		log.Errorf("PrepareJob(): %s", err.Error())
		return nil, err
	}
	defer conn.Close()
	grpcClient := proto.NewMinionClient(conn)

	grpcCtx, grpcCtxCancel := context.WithTimeout(ctx, minionRequestTimeout*time.Second)
	res, err := grpcClient.QueryJob(grpcCtx, &proto.MinionQueryJobRequest{})
	grpcCtxCancel()
	if err != nil {
		log.Errorf("QueryJob(): %s", err.Error())
		return nil, err
	}
	return res, nil
}
