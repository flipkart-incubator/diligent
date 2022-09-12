package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/datagen"
	"github.com/flipkart-incubator/diligent/pkg/intgen"
	"github.com/flipkart-incubator/diligent/pkg/metrics"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	"github.com/flipkart-incubator/diligent/pkg/work"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type JobState int

const (
	_ JobState = iota
	Prepared
	Running
	EndedSuccess
	EndedFailure
	EndedAborted
	EndedNeverRan
)

func (j JobState) String() string {
	switch j {
	case Prepared:
		return "Prepared"
	case Running:
		return "Running"
	case EndedSuccess:
		return "EndedSuccess"
	case EndedFailure:
		return "EndedFailure"
	case EndedAborted:
		return "EndedAborted"
	case EndedNeverRan:
		return "EndedNeverRan"
	}
	panic(fmt.Errorf("unknown job state %d", j))
}

func (j JobState) ToProto() proto.JobState {
	switch j {
	case Prepared:
		return proto.JobState_PREPARED
	case Running:
		return proto.JobState_RUNNING
	case EndedSuccess:
		return proto.JobState_ENDED_SUCCESS
	case EndedFailure:
		return proto.JobState_ENDED_FAILURE
	case EndedAborted:
		return proto.JobState_ENDED_ABORTED
	case EndedNeverRan:
		return proto.JobState_ENDED_NEVER_RAN
	}
	panic(fmt.Errorf("unknown job state %d", j))
}

type JobInfo struct {
	name           string
	state          JobState
	spec           *proto.JobSpec
	prepareTime    time.Time
	runTime        time.Time
	endTime        time.Time
	fatalErrors    int
	nonFatalErrors int
}

func (j *JobInfo) ToProto() *proto.JobInfo {
	if j == nil {
		return nil
	}
	return &proto.JobInfo{
		JobName:        j.name,
		JobSpec:        j.spec,
		JobState:       j.state.ToProto(),
		PrepareTime:    j.prepareTime.UnixMilli(),
		RunTime:        j.runTime.UnixMilli(),
		EndTime:        j.endTime.UnixMilli(),
		FatalErrors:    int32(j.fatalErrors),
		NonFatalErrors: int32(j.nonFatalErrors),
	}
}

type DataContext struct {
	dataSpec *datagen.Spec
	dataGen  *datagen.DataGen
}

type DBContext struct {
	driver string
	url    string
	db     *sql.DB
}

type WorkloadContext struct {
	workloadName  string
	assignedRange *intgen.Range
	tableName     string
	durationSec   int
	concurrency   int
	batchSize     int
	workload      *work.Workload
}

// Job represents a particular run of a workload
// It is thread safe
type Job struct {
	mut sync.Mutex

	name string

	spec     *proto.JobSpec
	data     *DataContext
	db       *DBContext
	workload *WorkloadContext
	metrics  *metrics.DiligentMetrics

	state          JobState
	prepareTime    time.Time
	runTime        time.Time
	endTime        time.Time
	fatalErrors    int
	nonFatalErrors int
}

func PrepareJob(ctx context.Context, spec *proto.JobSpec, metrics *metrics.DiligentMetrics) (*Job, error) {
	log.Infof("PrepareJob(%s)", spec.GetJobName())

	job := &Job{
		name:        spec.GetJobName(),
		state:       Prepared,
		spec:        spec,
		data:        nil,
		db:          nil,
		workload:    nil,
		metrics:     metrics,
		prepareTime: time.Now(),
	}

	// Load the dataspec
	err := job.loadDataSpec(ctx, spec.GetDataSpec())
	if err != nil {
		return nil, err
	}

	// Open connection with DB
	err = job.openDBConnection(ctx, spec.GetDbSpec())
	if err != nil {
		return nil, err
	}

	// Prepare workload
	err = job.prepareWorkload(ctx, spec.GetWorkloadSpec())
	if err != nil {
		return nil, err
	}

	log.Infof("PrepareJob(%s) completed successfully", spec.GetJobName())
	return job, nil
}

func (j *Job) Name() string {
	j.mut.Lock()
	defer j.mut.Unlock()
	return j.name
}

func (j *Job) State() JobState {
	j.mut.Lock()
	defer j.mut.Unlock()
	return j.state
}

func (j *Job) Info() *JobInfo {
	j.mut.Lock()
	defer j.mut.Unlock()
	return &JobInfo{
		name:           j.name,
		state:          j.state,
		spec:           j.spec,
		prepareTime:    j.prepareTime,
		runTime:        j.runTime,
		endTime:        j.endTime,
		fatalErrors:    j.fatalErrors,
		nonFatalErrors: j.nonFatalErrors,
	}
}

func (j *Job) Run(ctx context.Context) (chan int, error) {
	log.Infof("Run(%s)", j.name)
	j.mut.Lock()
	defer j.mut.Unlock()

	// Process only job is in prepared state
	if j.state != Prepared {
		return nil, fmt.Errorf("job is in %s state. cannot run", j.state.String())
	}

	notifyCh := make(chan int, 1) // Buffered so that we can notify and move on
	j.state = Running
	j.runTime = time.Now()

	go j.runWorkload(notifyCh)

	log.Infof("Run(%s) initiated successfully", j.name)
	return notifyCh, nil
}

func (j *Job) runWorkload(ch chan int) {
	j.mut.Lock()
	j.mut.Unlock()
	log.Infof("Starting workload for job %s...", j.name)

	resultCh := make(chan *work.WorkloadResult)
	j.workload.workload.Start(time.Duration(j.workload.durationSec)*time.Second, resultCh)

	log.Infof("Waiting for workload to complete for job %s...", j.name)
	result := <-resultCh
	log.Infof("Workload completed for job %s. Result=%v", j.name, result)

	j.mut.Lock()
	defer j.mut.Unlock()

	j.endTime = time.Now()
	j.fatalErrors = result.FatalErrors
	j.nonFatalErrors = result.NonFatalErrors

	switch {
	case result.IsAborted:
		j.state = EndedAborted
	case result.FatalErrors > 0 || result.NonFatalErrors > 0:
		j.state = EndedFailure
	default:
		j.state = EndedSuccess
	}

	log.Infof("Cleaning runtime state of job %s...", j.name)
	j.data = nil
	if j.db != nil {
		j.db.db.Close()
		j.db = nil
	}
	j.workload = nil

	// Notify waiter
	ch <- 0
}

func (j *Job) Abort(ctx context.Context) error {
	log.Infof("Abort(%s)", j.name)
	j.mut.Lock()
	defer j.mut.Unlock()

	if j.state == Prepared {
		j.state = EndedNeverRan
		j.endTime = time.Now()
	}

	if j.state == Running {
		j.workload.workload.Abort()
	}

	log.Infof("Abort(%s) completed successfully", j.name)
	return nil
}

func (j *Job) loadDataSpec(ctx context.Context, protoDs *proto.DataSpec) error {
	log.Infof("Loading dataspec...")
	dataSpec := proto.DataSpecFromProto(protoDs)
	j.data = &DataContext{
		dataSpec: dataSpec,
		dataGen:  datagen.NewDataGen(dataSpec),
	}
	log.Infof("Data spec loaded successfully")
	return nil
}

func (j *Job) openDBConnection(_ context.Context, protoDBSpec *proto.DBSpec) error {
	log.Infof("Opening DB connection...")
	// Validate driver
	driver := protoDBSpec.GetDriver()
	switch driver {
	case "mysql", "pgx":
	default:
		return fmt.Errorf("invalid driver: '%s'. Allowed values are 'mysql', 'pgx'", driver)
	}

	// Validate URL
	url := protoDBSpec.GetUrl()
	if url == "" {
		return fmt.Errorf("please specify the connection url")
	}

	// Close any existing connection
	if j.db != nil {
		j.db.db.Close()
		j.db = nil
	}

	// Open new connection
	db, err := sql.Open(driver, url)
	if err != nil {
		return err
	}

	err = work.ConnCheck(db)
	if err != nil {
		return err
	}

	j.db = &DBContext{
		driver: driver,
		url:    url,
		db:     db,
	}
	log.Infof("DB Connection successful (driver=%s, url=%s)", driver, url)
	return nil
}

func (j *Job) prepareWorkload(_ context.Context, protoWl *proto.WorkloadSpec) error {
	log.Infof("Preparing workload...")

	durationSec := int(protoWl.GetDurationSec())
	if durationSec < 0 {
		return fmt.Errorf("invalid duration %d. Must be >= 0", durationSec)
	}

	batchSize := int(protoWl.GetBatchSize())
	if batchSize < 1 {
		return fmt.Errorf("invalid batch size %d. Must be >= 1", batchSize)
	}

	concurrency := int(protoWl.GetConcurrency())
	if concurrency < 1 {
		return fmt.Errorf("invalid concurrency %d. Must be >= 1", concurrency)
	}

	table := protoWl.GetTableName()
	if table == "" {
		return fmt.Errorf("table name is missing")
	}

	inputRange := protoWl.GetAssignedRange()
	if inputRange == nil {
		return fmt.Errorf("assigned range is missing")
	}

	totalRows := j.data.dataSpec.KeyGenSpec.NumKeys()
	rangeStart := int(inputRange.GetStart())
	rangeLimit := int(inputRange.GetLimit())
	if rangeStart < 0 || rangeLimit < rangeStart {
		return fmt.Errorf("invalid range [%d, %d)", rangeStart, rangeLimit)
	}
	if rangeStart >= totalRows {
		return fmt.Errorf("invalid range [%d, %d)", rangeStart, rangeLimit)
	}
	if rangeLimit > totalRows {
		return fmt.Errorf("invalid range [%d, %d)", rangeStart, rangeLimit)
	}
	assignedRange := proto.RangeFromProto(inputRange)

	// Configure the db connections to maintain the specified concurrency
	j.db.db.SetMaxOpenConns(concurrency)
	j.db.db.SetMaxIdleConns(concurrency)

	rp := &work.RunParams{
		DB:          j.db.db,
		DataGen:     j.data.dataGen,
		Metrics:     j.metrics,
		Table:       table,
		Concurrency: concurrency,
		BatchSize:   batchSize,
		DurationSec: durationSec,
	}

	workloadName := protoWl.GetWorkloadName()
	var workload *work.Workload
	switch workloadName {
	case "insert":
		workload = work.NewInsertRowWorkload(assignedRange, rp)
	case "insert-txn":
		rp.EnableTxn = true
		workload = work.NewInsertTxnWorkload(assignedRange, rp)
	case "select-pk":
		workload = work.NewSelectByPkRowWorkload(assignedRange, rp)
	case "select-pk-txn":
		rp.EnableTxn = true
		workload = work.NewSelectByPkTxnWorkload(assignedRange, rp)
	case "select-uk":
		workload = work.NewSelectByUkRowWorkload(assignedRange, rp)
	case "select-uk-txn":
		rp.EnableTxn = true
		workload = work.NewSelectByUkTxnWorkload(assignedRange, rp)
	case "update":
		workload = work.NewUpdateRowWorkload(assignedRange, rp)
	case "update-txn":
		rp.EnableTxn = true
		workload = work.NewUpdateTxnWorkload(assignedRange, rp)
	case "delete":
		workload = work.NewDeleteRowWorkload(assignedRange, rp)
	case "delete-txn":
		rp.EnableTxn = true
		workload = work.NewDeleteTxnWorkload(assignedRange, rp)
	default:
		return fmt.Errorf("invalid workload '%s'", workloadName)
	}

	j.workload = &WorkloadContext{
		workloadName:  workloadName,
		assignedRange: assignedRange,
		tableName:     table,
		durationSec:   durationSec,
		concurrency:   concurrency,
		batchSize:     batchSize,
		workload:      workload,
	}

	log.Infof("Workload prepared successfully")
	return nil
}
