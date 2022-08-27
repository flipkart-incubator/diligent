package main

import (
	"context"
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	log "github.com/sirupsen/logrus"
	"sync"
)

// MinionManager manages a single minion
// It uses a MinionClient to execute gRPC operations on the minion
// It uses a MinionWatcher to watch for the health and completion state of a minion once the minion is prepared
// It is thread safe
type MinionManager struct {
	mut     sync.Mutex // Must be taken for all operations on the MinionManager
	addr    string
	client  *MinionClient
	watcher *MinionWatcher
}

func NewMinionManager(addr string) (*MinionManager, error) {
	client, err := NewMinionClient(addr)
	if err != nil {
		log.Errorf("Failed to create client for minion %s (%s)", addr, err.Error())
		return nil, err
	}
	return &MinionManager{
		addr:   addr,
		client: client,
	}, nil
}

func (m *MinionManager) Close() {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.client.Close()
}

func (m *MinionManager) GetAddr() string {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.addr
}

func (m *MinionManager) GetMinionStatus(ctx context.Context, ch chan *proto.MinionInfo) {
	res, err := m.client.Ping(ctx)
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

	res, err := m.client.PrepareJob(ctx, jobId, jobDesc, dataSpec, dbSpec, wlSpec)
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

	m.startWatchingMinion(res.GetPid(), jobId)
	log.Infof("Successfully prepared minion %s", m.addr)
	ch <- &proto.MinionStatus{
		Addr: m.addr,
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
	}
	return
}

func (m *MinionManager) startWatchingMinion(pid, jobId string) {
	if m.watcher != nil {
		m.watcher.StopWatching()
	}
	m.watcher = NewMinionWatcher(m.addr, m.client, pid, jobId)
	m.watcher.StartWatching()
}

func (m *MinionManager) RunJobOnMinion(ctx context.Context, ch chan *proto.MinionStatus) {
	if m.watcher.HasFailed() {
		ch <- &proto.MinionStatus{
			Addr: m.addr,
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: fmt.Sprintf("minion %s failed since it was prepared", m.addr),
			},
		}
	}

	res, err := m.client.RunJob(ctx)
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

func (m *MinionManager) AbortJobOnMinion(ctx context.Context, ch chan *proto.MinionStatus) {
	res, err := m.client.AbortJob(ctx)

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

	log.Infof("Successfully triggered abort on minion %s", m.addr)
	ch <- &proto.MinionStatus{
		Addr: m.addr,
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
	}
	return
}
