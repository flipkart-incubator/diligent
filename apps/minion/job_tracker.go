package main

import (
	"context"
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/metrics"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	"sync"
)

// JobTracker manages the lifecycle of jobs run by the minion
// It is thread safe
type JobTracker struct {
	mut        sync.Mutex
	currentJob *Job
}

func NewJobTracker() *JobTracker {
	return &JobTracker{
		currentJob: nil,
	}
}

func (t *JobTracker) CurrentJob() *Job {
	t.mut.Lock()
	defer t.mut.Unlock()
	return t.currentJob
}

func (t *JobTracker) CurrentJobInfo() *JobInfo {
	t.mut.Lock()
	defer t.mut.Unlock()
	if t.currentJob == nil {
		return nil
	}
	return t.currentJob.Info()
}

func (t *JobTracker) Prepare(ctx context.Context, id string, spec *proto.JobSpec, metrics *metrics.DiligentMetrics) error {
	t.mut.Lock()
	defer t.mut.Unlock()

	if t.currentJob != nil {
		err := t.currentJob.Abort(ctx)
		if err != nil {
			panic(fmt.Errorf("failed to abort current job %s", err.Error()))
		}
	}

	j, err := PrepareJob(ctx, id, spec, metrics)
	if err != nil {
		return err
	}
	t.currentJob = j
	return nil
}

func (t *JobTracker) Run(ctx context.Context) error {
	t.mut.Lock()
	defer t.mut.Unlock()

	if t.currentJob == nil {
		return fmt.Errorf("no current job")
	}
	err := t.currentJob.Run(ctx)
	return err
}

func (t *JobTracker) Abort(ctx context.Context) error {
	t.mut.Lock()
	defer t.mut.Unlock()

	if t.currentJob == nil {
		return fmt.Errorf("no current job")
	}
	err := t.currentJob.Abort(ctx)
	return err
}
