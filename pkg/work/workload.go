package work

import (
	"context"
	"sync"
	"time"
)

type Workload struct {
	name       string
	runParams  *RunParams
	workers    []*Worker
	mutex      *sync.Mutex
	ctx        context.Context
	cancelFunc context.CancelFunc
	waitGroup  *sync.WaitGroup
	isRunning  bool
}

func NewWorkload(name string, rp *RunParams) *Workload {
	return &Workload{
		name:       name,
		runParams:  rp,
		workers:    make([]*Worker, 0),
		mutex:      &sync.Mutex{},
		ctx:        nil,
		cancelFunc: nil,
		waitGroup:  &sync.WaitGroup{},
		isRunning:  false,
	}
}

func (wl *Workload) CreateWorkerWithWork(work CompositeWork) {
	wl.mutex.Lock()
	defer wl.mutex.Unlock()
	workerId := len(wl.workers)
	worker := NewWorker(workerId, wl.runParams, work)
	wl.workers = append(wl.workers, worker)
}

func (wl *Workload) Start(duration time.Duration) {
	wl.mutex.Lock()
	defer wl.mutex.Unlock()

	// Return from here if already running
	if wl.isRunning {
		return
	}

	// Update metrics that reflect configured values
	wl.runParams.Metrics.SetConfigMetricsForWorkload(wl.runParams.EnableTxn, wl.runParams.BatchSize, wl.runParams.Concurrency)
	defer wl.runParams.Metrics.UnsetConfigMetricsForWorkload()

	// Create a new context
	if duration == 0 {
		wl.ctx, wl.cancelFunc = context.WithCancel(context.Background())
	} else {
		wl.ctx, wl.cancelFunc = context.WithTimeout(context.Background(), duration)
	}

	// Update waitGroup count
	wl.waitGroup.Add(len(wl.workers))

	// Set workload to running
	wl.isRunning = true

	// Launch the workers
	for _, w := range wl.workers {
		go w.Run(wl.ctx, wl.waitGroup)
	}

	// Create a goroutine which will wait in the background on the workload to finish
	// that will result in the isRunning state to get updated to false on completion
	go wl.Wait()
}

func (wl *Workload) Cancel() {
	wl.mutex.Lock()
	defer wl.mutex.Unlock()
	if wl.cancelFunc != nil {
		wl.cancelFunc()
	}
}

func (wl *Workload) Wait() {
	// Wait for workers to finish
	wl.waitGroup.Wait()

	// All workers are done. Lock the workload and update its state
	wl.mutex.Lock()
	defer wl.mutex.Unlock()
	wl.isRunning = false
}

func (wl *Workload) IsRunning() bool {
	wl.mutex.Lock()
	defer wl.mutex.Unlock()
	return wl.isRunning
}
