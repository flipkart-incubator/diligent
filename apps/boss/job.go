package main

import (
	"context"
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/intgen"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type JobState int

const (
	_ JobState = iota
	JobNew
	JobPrepared
	JobRunning
	JobEndedSuccess
	JobEndedFailure
	JobEndedAborted
)

func (j JobState) String() string {
	switch j {
	case JobNew:
		return "New"
	case JobPrepared:
		return "Prepared"
	case JobRunning:
		return "Running"
	case JobEndedSuccess:
		return "EndedSuccess"
	case JobEndedFailure:
		return "EndedFailure"
	case JobEndedAborted:
		return "EndedAborted"
	default:
		panic(fmt.Sprintf("unknown job state %d", j))
	}
}

type Job struct {
	mut sync.Mutex

	name            string
	minions         map[string]*MinionProxy
	jobState        JobState
	minionEndStates []MinionWatchStatus

	createTime  time.Time
	prepareTime time.Time
	runTime     time.Time
	endTime     time.Time
}

func NewJob(name string, minions map[string]*MinionProxy) *Job {
	return &Job{
		name:            name,
		minions:         minions,
		jobState:        JobNew,
		minionEndStates: make([]MinionWatchStatus, len(minions)),
		createTime:      time.Now(),
	}
}

func (j *Job) HasEnded() bool {
	j.mut.Lock()
	defer j.mut.Unlock()
	log.Infof("job.HasEnded(): name=%s, current state = %s", j.name, j.jobState.String())
	switch j.jobState {
	case JobEndedSuccess, JobEndedFailure, JobEndedAborted:
		log.Infof("job.HasEnded(): true")
		return true
	default:
		log.Infof("job.HasEnded(): false")
		return false
	}
}

func (j *Job) Name() string {
	j.mut.Lock()
	defer j.mut.Unlock()
	return j.name
}

func (j *Job) State() JobState {
	j.mut.Lock()
	defer j.mut.Unlock()
	return j.jobState
}

func (j *Job) Prepare(ctx context.Context, in *proto.BossPrepareJobRequest) (*proto.BossPrepareJobResponse, error) {
	j.mut.Lock()
	defer j.mut.Unlock()

	// Partition the data among the number of minions
	dataSpec := proto.DataSpecFromProto(in.GetJobSpec().GetDataSpec())
	numRecs := dataSpec.KeyGenSpec.NumKeys()
	fullRange := intgen.NewRange(0, numRecs)
	numMinions := len(j.minions)
	var assignedRanges []*intgen.Range

	// Get the job level workload spec, and do homework to determine per-minion workload spec
	workloadName := in.GetJobSpec().GetWorkloadSpec().GetWorkloadName()
	switch workloadName {
	case "insert", "insert-txn", "delete", "delete-txn":
		assignedRanges = fullRange.Partition(numMinions)
	case "select-pk", "select-pk-txn", "select-uk", "select-uk-txn", "update", "update-txn":
		assignedRanges = fullRange.Duplicate(numMinions)
	default:
		reason := fmt.Sprintf("invalid workload '%s'", workloadName)
		log.Infof("RunWorkload(): %s", reason)
		return &proto.BossPrepareJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}, fmt.Errorf(reason)
	}
	log.Infof("PrepareJob(): Data partitioned among minions: %v", assignedRanges)

	// Slices to gather data on individual minions
	addrs := make([]string, numMinions)
	rchs := make([]chan *proto.MinionPrepareJobResponse, numMinions)
	echs := make([]chan error, numMinions)

	// Invoke on individual minions
	i := 0
	for addr, proxy := range j.minions {
		log.Infof("PrepareJob(): Preparing Minion %s", addr)
		wlSpec := &proto.WorkloadSpec{
			WorkloadName:  in.GetJobSpec().GetWorkloadSpec().GetWorkloadName(),
			AssignedRange: proto.RangeToProto(assignedRanges[i]),
			TableName:     in.GetJobSpec().GetWorkloadSpec().GetTableName(),
			DurationSec:   in.GetJobSpec().GetWorkloadSpec().GetDurationSec(),
			Concurrency:   in.GetJobSpec().GetWorkloadSpec().GetConcurrency(),
			BatchSize:     in.GetJobSpec().GetWorkloadSpec().GetBatchSize(),
		}

		addrs[i] = addr
		rchs[i], echs[i] = proxy.PrepareJobAsync(ctx, in.GetJobSpec().GetJobName(), in.GetJobSpec().GetDataSpec(), in.GetJobSpec().GetDbSpec(), wlSpec)
		i++
	}

	// For holding overall result of this call and status of individual minions
	overallStatus := proto.GeneralStatus{
		IsOk:          true,
		FailureReason: "",
	}
	minionStatuses := make([]*proto.MinionStatus, numMinions)

	// Collect execution results
	for i, _ = range minionStatuses {
		select {
		case res := <-rchs[i]:
			minionStatuses[i] = &proto.MinionStatus{
				Addr: addrs[i],
				Status: &proto.GeneralStatus{
					IsOk:          res.GetStatus().GetIsOk(),
					FailureReason: res.GetStatus().GetFailureReason(),
				},
			}
		case err := <-echs[i]:
			overallStatus = proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "errors encountered on one or more minions",
			}
			minionStatuses[i] = &proto.MinionStatus{
				Addr: addrs[i],
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: err.Error(),
				},
			}
		}
	}

	j.jobState = JobPrepared
	j.prepareTime = time.Now()
	go j.WatchForEnd()
	log.Infof("PrepareJob(): completed")
	return &proto.BossPrepareJobResponse{
		Status:         &overallStatus,
		MinionStatuses: minionStatuses,
	}, nil
}

func (j *Job) Run(ctx context.Context) (*proto.BossRunJobResponse, error) {
	log.Infof("GRPC: RunJob()")
	j.mut.Lock()
	defer j.mut.Unlock()

	// Slices to gather data on individual minions
	numMinions := len(j.minions)
	addrs := make([]string, numMinions)
	rchs := make([]chan *proto.MinionRunJobResponse, numMinions)
	echs := make([]chan error, numMinions)

	// Invoke on individual minions
	i := 0
	for addr, proxy := range j.minions {
		log.Infof("RunJob(): Triggering run on Minion %s", addr)
		addrs[i] = addr
		rchs[i], echs[i] = proxy.RunJobAsync(ctx)
		i++
	}

	// For holding overall result of this call and status of individual minions
	overallStatus := proto.GeneralStatus{
		IsOk:          true,
		FailureReason: "",
	}
	minionStatuses := make([]*proto.MinionStatus, numMinions)

	// Collect execution results
	for i, _ = range minionStatuses {
		select {
		case res := <-rchs[i]:
			minionStatuses[i] = &proto.MinionStatus{
				Addr: addrs[i],
				Status: &proto.GeneralStatus{
					IsOk:          res.GetStatus().GetIsOk(),
					FailureReason: res.GetStatus().GetFailureReason(),
				},
			}
		case err := <-echs[i]:
			overallStatus = proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "errors encountered on one or more minions",
			}
			minionStatuses[i] = &proto.MinionStatus{
				Addr: addrs[i],
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: err.Error(),
				},
			}
		}
	}

	j.jobState = JobRunning
	j.runTime = time.Now()
	log.Infof("RunJob(): completed")
	return &proto.BossRunJobResponse{
		Status:         &overallStatus,
		MinionStatuses: minionStatuses,
	}, nil
}

func (j *Job) Abort(ctx context.Context) (*proto.BossAbortJobResponse, error) {
	log.Infof("GRPC: AbortJob()")
	j.mut.Lock()
	defer j.mut.Unlock()

	// Slices to gather data on individual minions
	numMinions := len(j.minions)
	addrs := make([]string, numMinions)
	rchs := make([]chan *proto.MinionAbortJobResponse, numMinions)
	echs := make([]chan error, numMinions)

	// Invoke on individual minions
	i := 0
	for addr, proxy := range j.minions {
		log.Infof("AbortJob(): Triggering run on Minion %s", addr)
		addrs[i] = addr
		rchs[i], echs[i] = proxy.AbortJobAsync(ctx)
		i++
	}

	// For holding overall result of this call and status of individual minions
	overallStatus := proto.GeneralStatus{
		IsOk:          true,
		FailureReason: "",
	}
	minionStatuses := make([]*proto.MinionStatus, numMinions)

	// Collect execution results
	for i, _ = range minionStatuses {
		select {
		case res := <-rchs[i]:
			minionStatuses[i] = &proto.MinionStatus{
				Addr: addrs[i],
				Status: &proto.GeneralStatus{
					IsOk:          res.GetStatus().GetIsOk(),
					FailureReason: res.GetStatus().GetFailureReason(),
				},
			}
		case err := <-echs[i]:
			overallStatus = proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "errors encountered on one or more minions",
			}
			minionStatuses[i] = &proto.MinionStatus{
				Addr: addrs[i],
				Status: &proto.GeneralStatus{
					IsOk:          false,
					FailureReason: err.Error(),
				},
			}
		}
	}

	log.Infof("AbortJob(): completed")
	return &proto.BossAbortJobResponse{
		Status:         &overallStatus,
		MinionStatuses: minionStatuses,
	}, nil
}

func (j *Job) WatchForEnd() {
	i := 0
	jobEndStatus := JobEndedSuccess
	for addr, proxy := range j.minions {
		minionWatchStatus := <-proxy.Watch()
		log.Warnf("Minion %s: got watch status: %s", addr, minionWatchStatus.String())

		j.minionEndStates[i] = minionWatchStatus

		// Job end state is set based the first non-success status if any
		if jobEndStatus == JobEndedSuccess {
			switch minionWatchStatus {
			case EndedSuccess:
				// Nothing to do
			case EndedFailure:
				jobEndStatus = JobEndedFailure
			case EndedAborted:
				jobEndStatus = JobEndedAborted
			case Unreachable:
				jobEndStatus = JobEndedFailure
			case Restarted:
				jobEndStatus = JobEndedFailure
			case Errored:
				jobEndStatus = JobEndedFailure
			}
		}
		i++
	}

	j.mut.Lock()
	defer j.mut.Unlock()
	j.jobState = jobEndStatus
	j.endTime = time.Now()
	log.Warnf("Job %s ended: status=%s, time=%s", j.name, j.jobState.String(), j.endTime.Format(time.UnixDate))
}
