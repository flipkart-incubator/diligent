package main

import (
	"context"
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/metrics"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	"sync"
)

type JobTrackerStateName int

const (
	_ JobTrackerStateName = iota
	JTIdle
	JTPrepared
	JTRunning
)

type JobTrackerState interface {
	State() JobTrackerStateName
	CurrentJob() *Job
	Prepare(ctx context.Context, jt *JobTracker, id string, spec *proto.JobSpec, metrics *metrics.DiligentMetrics) (JobTrackerState, error)
	Run(ctx context.Context, jt *JobTracker) (JobTrackerState, error)
	Abort(ctx context.Context, jt *JobTracker) (JobTrackerState, error)
}

type idleState struct {
}

func NewIdleState() *idleState {
	return &idleState{}
}

func (s *idleState) State() JobTrackerStateName {
	return JTIdle
}

func (s *idleState) CurrentJob() *Job {
	return nil
}

func (s *idleState) Prepare(ctx context.Context, jt *JobTracker, id string, spec *proto.JobSpec, metrics *metrics.DiligentMetrics) (JobTrackerState, error) {
	newJob, err := PrepareJob(ctx, id, spec, metrics)
	if err != nil {
		return s, err
	}
	jt.addJobEntry(newJob)
	return NewPreparedState(newJob), nil
}

func (s *idleState) Run(ctx context.Context, jt *JobTracker) (JobTrackerState, error) {
	return s, fmt.Errorf("no current job")
}

func (s *idleState) Abort(ctx context.Context, jt *JobTracker) (JobTrackerState, error) {
	return s, fmt.Errorf("no current job")
}

type preparedState struct {
	job *Job
}

func NewPreparedState(job *Job) *preparedState {
	return &preparedState{job: job}
}

func (s *preparedState) State() JobTrackerStateName {
	return JTPrepared
}

func (s *preparedState) CurrentJob() *Job {
	return s.job
}

func (s *preparedState) Prepare(ctx context.Context, jt *JobTracker, id string, spec *proto.JobSpec, metrics *metrics.DiligentMetrics) (JobTrackerState, error) {
	_ = s.job.Abort(ctx)
	newJob, err := PrepareJob(ctx, id, spec, metrics)
	if err != nil {
		return s, err
	}
	return NewPreparedState(newJob), nil
}

func (s *preparedState) Run(ctx context.Context, jt *JobTracker) (JobTrackerState, error) {
	done, err := s.job.Run(ctx)
	if err != nil {
		return s, err
	}
	return NewRunningState(s.job, jt, done), nil
}

func (s *preparedState) Abort(ctx context.Context, jt *JobTracker) (JobTrackerState, error) {
	_ = s.job.Abort(ctx)
	return NewIdleState(), nil
}

type runningState struct {
	job        *Job
	jobDone    chan int
	cancelWait context.CancelFunc
}

func NewRunningState(job *Job, jt *JobTracker, jobDone chan int) *runningState {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-jobDone:
				jt.MakeIdle()
			case <-ctx.Done():
				break
			}
		}
	}()
	return &runningState{
		job:        job,
		jobDone:    jobDone,
		cancelWait: cancel,
	}
}

func (s *runningState) State() JobTrackerStateName {
	return JTRunning
}

func (s *runningState) CurrentJob() *Job {
	return s.job
}

func (s *runningState) Prepare(ctx context.Context, jt *JobTracker, id string, spec *proto.JobSpec, metrics *metrics.DiligentMetrics) (JobTrackerState, error) {
	return s, fmt.Errorf("job %s is running, cannot prepare another job", s.job.Id())
}

func (s *runningState) Run(ctx context.Context, jt *JobTracker) (JobTrackerState, error) {
	return s, fmt.Errorf("job %s is already running", s.job.Id())
}

func (s *runningState) Abort(ctx context.Context, jt *JobTracker) (JobTrackerState, error) {
	// First stop waiting for job to finish. This is to avoid a race condition where the goroutine that
	// is waiting for the job to finish will also attempt to move the JobTracker to idle state
	s.cancelWait()
	_ = s.job.Abort(ctx)
	return NewIdleState(), nil
}

// JobTracker manages the lifecycle of jobs run by the minion
// It is thread safe
type JobTracker struct {
	mut        sync.Mutex
	state      JobTrackerState
	jobHistory map[string]*Job
}

func NewJobTracker() *JobTracker {
	return &JobTracker{
		state:      NewIdleState(),
		jobHistory: make(map[string]*Job),
	}
}

// CurrentJob Returns the current job. Can be nil
func (t *JobTracker) CurrentJob() *Job {
	t.mut.Lock()
	defer t.mut.Unlock()
	return t.state.CurrentJob()
}

func (t *JobTracker) Prepare(ctx context.Context, id string, spec *proto.JobSpec, metrics *metrics.DiligentMetrics) error {
	t.mut.Lock()
	defer t.mut.Unlock()
	nextState, err := t.state.Prepare(ctx, t, id, spec, metrics)
	t.state = nextState
	return err
}

func (t *JobTracker) Run(ctx context.Context) error {
	t.mut.Lock()
	defer t.mut.Unlock()
	nextState, err := t.state.Run(ctx, t)
	t.state = nextState
	return err
}

func (t *JobTracker) Abort(ctx context.Context) error {
	t.mut.Lock()
	defer t.mut.Unlock()
	nextState, err := t.state.Abort(ctx, t)
	t.state = nextState
	return err
}

func (t *JobTracker) GetJobInfo(id string) *JobInfo {
	t.mut.Lock()
	defer t.mut.Unlock()
	job := t.jobHistory[id]
	if job == nil {
		return nil
	}
	return job.Info()
}

func (t *JobTracker) MakeIdle() {
	t.mut.Lock()
	defer t.mut.Unlock()
	t.state = NewIdleState()
}

// Assumed to be called internally with mutex already held
func (t *JobTracker) addJobEntry(job *Job) {
	t.jobHistory[job.Id()] = job
}
