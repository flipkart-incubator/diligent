package work

import (
	"context"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type WorkloadResult struct {
	IsAborted      bool
	FatalErrors    int
	NonFatalErrors int
}

type Workload struct {
	name         string
	runParams    *RunParams
	workers      []*Worker
	mutex        *sync.Mutex
	ctx          context.Context
	cancelFunc   context.CancelFunc
	isRunStarted bool
	isRunEnded   bool
	isAborted    bool
}

func NewWorkload(name string, rp *RunParams) *Workload {
	return &Workload{
		name:      name,
		runParams: rp,
		workers:   make([]*Worker, 0),
		mutex:     &sync.Mutex{},
	}
}

func (wl *Workload) CreateWorkerWithWork(work CompositeWork) {
	wl.mutex.Lock()
	defer wl.mutex.Unlock()
	workerId := len(wl.workers)
	worker := NewWorker(workerId, wl.runParams, work)
	wl.workers = append(wl.workers, worker)
}

func (wl *Workload) Start(duration time.Duration, ch chan *WorkloadResult) {
	wl.mutex.Lock()
	defer wl.mutex.Unlock()

	// Return from here if already running
	if wl.isRunStarted {
		log.Fatal("attempt to run a workload again")
	}

	// Note that run was started
	wl.isRunStarted = true

	// Update metrics that reflect configured values
	wl.runParams.Metrics.SetConfigMetricsForWorkload(wl.runParams.EnableTxn, wl.runParams.BatchSize, wl.runParams.Concurrency)
	defer wl.runParams.Metrics.UnsetConfigMetricsForWorkload()

	// Create a new context
	if duration == 0 {
		wl.ctx, wl.cancelFunc = context.WithCancel(context.Background())
	} else {
		wl.ctx, wl.cancelFunc = context.WithTimeout(context.Background(), duration)
	}

	// Launch the workers
	workerResultCh := make(chan *WorkerResult)
	for _, w := range wl.workers {
		go w.Run(wl.ctx, workerResultCh)
	}

	// Launch goroutine that will collect result of individual workers and send result of the workload to caller
	go wl.handleCompletion(workerResultCh, ch)
}

func (wl *Workload) Abort() {
	wl.mutex.Lock()
	defer wl.mutex.Unlock()

	if wl.isRunEnded {
		return
	}

	if wl.cancelFunc != nil {
		wl.isAborted = true
		wl.cancelFunc()
	}
}

func (wl *Workload) handleCompletion(inCh chan *WorkerResult, outCh chan *WorkloadResult) {
	// Collect one result from each worker
	N := len(wl.workers)
	fatalErrors := 0
	nonFatalErrors := 0
	for i := 0; i < N; i++ {
		wr := <-inCh

		if wr.FatalError {
			fatalErrors++
		}
		nonFatalErrors += wr.NonFatalErrors
	}

	// All workers are done. Lock the workload and update its state
	wl.mutex.Lock()
	defer wl.mutex.Unlock()

	// Note that run has finished
	wl.isRunEnded = true

	outCh <- &WorkloadResult{
		IsAborted:      wl.isAborted,
		FatalErrors:    fatalErrors,
		NonFatalErrors: nonFatalErrors,
	}
}
