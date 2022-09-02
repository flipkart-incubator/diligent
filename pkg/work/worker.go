package work

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type WorkerResult struct {
	FatalError     bool
	NonFatalErrors int
}

type Worker struct {
	id        int
	runParams *RunParams
	work      CompositeWork
}

func NewWorker(id int, rp *RunParams, work CompositeWork) *Worker {
	return &Worker{
		id:        id,
		runParams: rp,
		work:      work,
	}
}

func (w *Worker) Run(ctx context.Context, ch chan *WorkerResult) {
	log.Infof("Worker %d: Started\n", w.id)
	defer log.Infof("Worker %d: Stopped\n", w.id)

	w.runParams.Metrics.IncConcurrencyForWorkload()
	defer w.runParams.Metrics.DecConcurrencyForWorkload()

	err := ConnCheck(w.runParams.DB)
	if err != nil {
		log.Errorf("Encountered fatal error: %s", err.Error())
		ch <- &WorkerResult{FatalError: true}
		return
	}

	nonFatalErrors := 0
	more := true
	for more {
		if ctx.Err() != nil {
			log.Infof("Worker %d: Context marked done (%s)\n", w.id, ctx.Err())
			ch <- &WorkerResult{NonFatalErrors: nonFatalErrors}
			return
		}
		w.runParams.Metrics.ObserveDbConn(w.runParams.DB)
		more, err = w.work.DoNext()
		if err != nil {
			log.Errorf("Encountered non-fatal error: %s", err.Error())
			nonFatalErrors++
		}
	}
	ch <- &WorkerResult{NonFatalErrors: nonFatalErrors}
}
