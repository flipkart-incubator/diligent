package main

import (
	"context"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	MinionPollPeriodSecs = 5
)

// MinionWatcher polls a minion to detect any failures, or detect job completion
// It uses a MinionClient to execute the gRPC operations on the minion
// It is thread safe
type MinionWatcher struct {
	mut          sync.Mutex // Must be taken for all operations on the MinionWatcher
	addr         string
	client       *MinionClient
	pid          string
	jobId        string
	pollCount    int
	watchStarted bool
	hasFinished  bool
	hasFailed    bool
	cancel       context.CancelFunc
}

func NewMinionWatcher(addr string, client *MinionClient, pid string, jobId string) *MinionWatcher {
	return &MinionWatcher{
		addr:         addr,
		client:       client,
		pid:          pid,
		jobId:        jobId,
		pollCount:    0,
		watchStarted: false,
		hasFailed:    false,
		hasFinished:  false,
		cancel:       nil,
	}
}

func (f *MinionWatcher) StartWatching() {
	log.Infof("StartWatching() on minion %s", f.addr)
	f.mut.Lock()
	defer f.mut.Unlock()
	if f.watchStarted == true {
		panic("attempt to restart a MinionWatcher")
	}
	ctx, cancel := context.WithCancel(context.Background())
	f.cancel = cancel
	f.watchStarted = true
	go f.pollContinuously(ctx)
}

func (f *MinionWatcher) StopWatching() {
	log.Infof("StopWatching() on minion %s", f.addr)
	f.mut.Lock()
	defer f.mut.Unlock()
	if f.watchStarted == false {
		panic("attempt to stop a MinionWatcher that was never started")
	}
	if f.cancel != nil {
		f.cancel()
	}
}

func (f *MinionWatcher) HasFailed() bool {
	f.mut.Lock()
	defer f.mut.Unlock()
	if f.watchStarted == false {
		panic("attempt to query a MinionWatcher that was never started")
	}
	return f.hasFailed
}

func (f *MinionWatcher) HasFinished() bool {
	f.mut.Lock()
	defer f.mut.Unlock()
	if f.watchStarted == false {
		panic("attempt to query a MinionWatcher that was never started")
	}
	return f.hasFinished
}

func (f *MinionWatcher) pollContinuously(ctx context.Context) {
	for {
		if ctx.Err() != nil || f.HasFailed() || f.HasFinished() {
			break
		}
		time.Sleep(MinionPollPeriodSecs * time.Second)
		f.pollOnce()
	}
	f.cancel() // To ensure cleanup of context
}

func (f *MinionWatcher) pollOnce() {
	log.Infof("Querying job status for minion %s. Count=%d", f.addr, f.pollCount)
	res, err := f.client.QueryJob(context.TODO(), f.jobId)

	f.mut.Lock()
	defer f.mut.Unlock()

	f.pollCount++

	if err != nil {
		// The query failed
		log.Infof("Job status query failed for minion %s, marking failed (reason=%s)", f.addr, err.Error())
		f.hasFailed = true
	} else {
		// We got a response
		if res.GetStatus().GetIsOk() == false {
			// The response is negative
			log.Infof("Job status query failed for minion %s, marking failed (reason=%s)", f.addr, res.GetStatus().GetFailureReason())
			f.hasFailed = true
		} else {
			// The response is positive
			if res.GetPid() != f.pid {
				// Though we got a response, it looks like the minion has restarted
				log.Infof("Pid mismatch on minion %s (expected=%s, actual=%s), marking failed", f.addr, f.pid, res.GetPid())
				f.hasFailed = true
			} else {
				// The response is from the correct minion
				switch res.GetJobInfo().GetJobState() {
				case proto.JobState_ENDED_SUCCESS:
					log.Infof("Got ENDED_SUCCESS for minion %s, marking finished", f.addr)
					f.hasFinished = true
				case proto.JobState_ENDED_FAILURE:
					log.Infof("Got ENDED_FAILURE for minion %s, marking finished", f.addr)
					f.hasFinished = true
				case proto.JobState_ENDED_ABORTED:
					log.Infof("Got ENDED_ABORTED for minion %s, marking finished", f.addr)
					f.hasFinished = true
				case proto.JobState_ENDED_NEVER_RAN:
					log.Infof("Got ENDED_NEVER_RAN for minion %s, marking finished", f.addr)
					f.hasFinished = true
				}
			}
		}
	}
}
