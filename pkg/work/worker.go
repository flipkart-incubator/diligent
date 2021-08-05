package work

import (
	"context"
	log "github.com/sirupsen/logrus"
	"sync"
)

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

func (w *Worker) Run(ctx context.Context, waitGroup *sync.WaitGroup) {
	log.Infof("Worker %d: Started\n", w.id)
	defer log.Infof("Worker %d: Stopped\n", w.id)
	defer waitGroup.Done()

	w.runParams.Metrics.IncConcurrencyForWorkload()
	defer w.runParams.Metrics.DecConcurrencyForWorkload()

	err := ConnCheck(w.runParams.DB)
	if err != nil {
		log.Fatal(err)
	}

	more := true
	for more {
		select {
		case <-ctx.Done():
			log.Infof("Worker %d: Context marked done (%s)\n", w.id, ctx.Err())
			return
		default:
		}
		w.runParams.Metrics.ObserveDbConn(w.runParams.DB)
		more = w.work.DoNext()
	}
}
