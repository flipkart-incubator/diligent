package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/buildinfo"
	"github.com/flipkart-incubator/diligent/pkg/datagen"
	"github.com/flipkart-incubator/diligent/pkg/idgen"
	"github.com/flipkart-incubator/diligent/pkg/intgen"
	"github.com/flipkart-incubator/diligent/pkg/metrics"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	"github.com/flipkart-incubator/diligent/pkg/work"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

const (
	bossConnectionTimoutSecs = 5
	bossRequestTimeoutSecs   = 5
)

type MinionState int

const (
	Idle MinionState = iota
	Prepared
	Running
)

type DataContext struct {
	dataSpecName string
	dataSpec     *datagen.Spec
	dataGen      *datagen.DataGen
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

// MinionServer represents a diligent minion gRPC server and associated state
// It is thread safe
type MinionServer struct {
	proto.UnimplementedMinionServer
	mut *sync.Mutex //mut Is used to coordinate access to this struct

	grpcAddr      string
	metricsAddr   string
	advertiseAddr string
	bossAddr      string

	pid       string    //pid is a string that represents a unique run of the minion server
	startTime time.Time //startTime is when this server started executing

	state    MinionState
	data     *DataContext
	db       *DBContext
	workload *WorkloadContext

	metrics *metrics.DiligentMetrics
}

func NewMinionServer(grpcAddr, metricsAddr, advertiseAddr, bossAddr string) *MinionServer {
	return &MinionServer{
		mut: &sync.Mutex{},

		grpcAddr:      grpcAddr,
		metricsAddr:   metricsAddr,
		advertiseAddr: advertiseAddr,
		bossAddr:      bossAddr,

		pid:       idgen.GenerateId16(),
		startTime: time.Now(),

		state:    Idle,
		data:     nil,
		db:       nil,
		workload: nil,
		metrics:  nil,
	}
}

func (s *MinionServer) RegisterWithBoss() error {
	log.Infof("Registering with boss")
	var conn *grpc.ClientConn
	var err error
	for {
		dialCtx, dialCancel := context.WithTimeout(context.Background(), bossConnectionTimoutSecs*time.Second)
		conn, err = grpc.DialContext(dialCtx, s.bossAddr, grpc.WithInsecure(), grpc.WithBlock())
		dialCancel()
		if err != nil {
			log.Errorf("failed to connect to boss %s (%v)", s.bossAddr, err)
			continue
		}

		bossClient := proto.NewBossClient(conn)

		grcpCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
		_, err = bossClient.RegisterMinion(grcpCtx, &proto.BossRegisterMinionRequest{Addr: s.advertiseAddr})
		grpcCancel()
		if err != nil {
			log.Errorf("Registration failed with error: (%v)\n", err)
			continue
		} else {
			break
		}
	}

	log.Infof("Successfully registered with boss")
	return nil
}

func (s *MinionServer) Serve() error {
	log.Infof("Starting to serve...")

	// Start metrics serving
	s.metrics = metrics.NewDiligentMetrics(s.metricsAddr)
	s.metrics.Register()

	// Create listening port
	listener, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		return err
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register our this server object with the gRPC server
	proto.RegisterMinionServer(grpcServer, s)

	// Start serving
	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func (s *MinionServer) Ping(_ context.Context, in *proto.MinionPingRequest) (*proto.MinionPingResponse, error) {
	log.Infof("GRPC: Ping")
	s.mut.Lock()
	defer s.mut.Unlock()

	log.Infof("GRPC: Ping completed successfully")
	return &proto.MinionPingResponse{
		BuildInfo: &proto.BuildInfo{
			AppName:    buildinfo.AppName,
			AppVersion: buildinfo.AppVersion,
			CommitHash: buildinfo.CommitHash,
			GoVersion:  buildinfo.GoVersion,
			BuildTime:  buildinfo.BuildTime,
		},
		ProcessInfo: &proto.ProcessInfo{
			Pid:       s.pid,
			StartTime: s.startTime.Format(time.UnixDate),
			Uptime:    time.Since(s.startTime).String(),
		},
	}, nil
}

func (s *MinionServer) PrepareJob(ctx context.Context, in *proto.MinionPrepareJobRequest) (*proto.MinionPrepareJobResponse, error) {
	log.Infof("GRPC: PrepareJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	log.Infof("Minion state=%s", s.state.String())

	// Reject if minion is running a job currently
	if s.state == Running {
		log.Infof("PrepareJob(): Cannot prepare as minion is currently running a job")
		return &proto.MinionPrepareJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "minion is running a job",
			},
		}, nil
	}

	// Clean up earlier state if minion was already in prepared state
	if s.state == Prepared {
		log.Infof("PrepareJob(): Minion was alread prepared. Cleaning up previous state")
		s.data = nil
		if s.db != nil {
			s.db.db.Close()
			s.db = nil
		}
		s.workload = nil
	}

	// Load the dataspec
	err := s.loadDataSpec(ctx, "TODO", in.GetJobSpec().GetDataSpec())
	if err != nil {
		log.Errorf("GRPC: PrepareJob(jobId=%s) Failed: %s", in.GetJobId(), err.Error())
		return &proto.MinionPrepareJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}, nil
	}

	// Open connection with DB
	err = s.openDBConnection(ctx, in.GetJobSpec().GetDbSpec())
	if err != nil {
		log.Errorf("GRPC: PrepareJob(jobId=%s) Failed: %s", in.GetJobId(), err.Error())
		return &proto.MinionPrepareJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}, nil
	}

	// Prepare workload
	err = s.prepareWorkload(ctx, in.GetJobSpec().GetWorkloadSpec())
	if err != nil {
		log.Errorf("GRPC: PrepareJob(jobId=%s) Failed: %s", in.GetJobId(), err.Error())
		return &proto.MinionPrepareJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}, nil
	}

	// Preparation was successful
	s.state = Prepared
	log.Infof("GRPC: PrepareJob(jobId=%s) completed successfully", in.GetJobId())
	return &proto.MinionPrepareJobResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
	}, nil
}

func (s *MinionServer) RunJob(ctx context.Context, in *proto.MinionRunJobRequest) (*proto.MinionRunJobResponse, error) {
	log.Infof("GRPC: RunJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	// Process only if minion is prepared
	log.Infof("Minion state=%s", s.state.String())
	if s.state != Prepared {
		return &proto.MinionRunJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "minion is not in prepared state. current state: " + s.state.String(),
			},
		}, nil
	}

	go func() {
		s.mut.Lock()
		s.mut.Unlock()
		log.Infof("Starting workload...")
		s.workload.workload.Start(time.Duration(s.workload.durationSec) * time.Second)
		log.Infof("Waiting for workload to complete...")
		s.workload.workload.WaitForCompletion()
		log.Infof("Workload completed. Going into idle state...")
		s.mut.Lock()
		s.state = Idle
		s.mut.Unlock()
	}()

	s.state = Running
	log.Infof("GRPC: RunJob() completed successfully")
	return &proto.MinionRunJobResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
	}, nil
}

func (s *MinionServer) loadDataSpec(ctx context.Context, name string, protoDs *proto.DataSpec) error {
	log.Infof("Loading dataspec...")
	dataSpec := proto.DataSpecFromProto(protoDs)
	s.data = &DataContext{
		dataSpecName: name,
		dataSpec:     dataSpec,
		dataGen:      datagen.NewDataGen(dataSpec),
	}
	log.Infof("Data spec loaded successfully")
	return nil
}

func (s *MinionServer) openDBConnection(_ context.Context, protoDBSpec *proto.DBSpec) error {
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
	if s.db != nil {
		s.db.db.Close()
		s.db = nil
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

	s.db = &DBContext{
		driver: driver,
		url:    url,
		db:     db,
	}
	log.Infof("DB Connection successful (driver=%s, url=%s)", driver, url)
	return nil
}

func (s *MinionServer) prepareWorkload(_ context.Context, protoWl *proto.WorkloadSpec) error {
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

	totalRows := s.data.dataSpec.KeyGenSpec.NumKeys()
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
	s.db.db.SetMaxOpenConns(concurrency)
	s.db.db.SetMaxIdleConns(concurrency)

	rp := &work.RunParams{
		DB:          s.db.db,
		DataGen:     s.data.dataGen,
		Metrics:     s.metrics,
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
		workload = work.NewInsertTxnWorkload(assignedRange, rp)
	case "select-pk":
		workload = work.NewSelectByPkRowWorkload(assignedRange, rp)
	case "select-pk-txn":
		workload = work.NewSelectByPkTxnWorkload(assignedRange, rp)
	case "select-uk":
		workload = work.NewSelectByUkRowWorkload(assignedRange, rp)
	case "select-uk-txn":
		workload = work.NewSelectByUkTxnWorkload(assignedRange, rp)
	case "update":
		workload = work.NewUpdateRowWorkload(assignedRange, rp)
	case "update-txn":
		workload = work.NewUpdateTxnWorkload(assignedRange, rp)
	case "delete":
		workload = work.NewDeleteRowWorkload(assignedRange, rp)
	case "delete-txn":
		workload = work.NewDeleteTxnWorkload(assignedRange, rp)
	default:
		return fmt.Errorf("invalid workload '%s'", workloadName)
	}

	s.workload = &WorkloadContext{
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

func (s *MinionServer) StopJob(_ context.Context, in *proto.MinionStopJobRequest) (*proto.MinionStopJobResponse, error) {
	log.Infof("GRPC: StopJob()")
	s.mut.Lock()
	defer s.mut.Unlock()

	// Process only if minion is running a job and id matches
	if s.state != Running {
		return &proto.MinionStopJobResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "minion is not in running state. current state: " + s.state.String(),
			},
		}, nil
	}

	if s.workload.workload.IsRunning() {
		s.workload.workload.Cancel()
		s.workload.workload.WaitForCompletion()
	}

	log.Infof("Stopped job successfully")
	return &proto.MinionStopJobResponse{
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
	}, nil
}

func (s *MinionServer) GetDataSpecInfo(context.Context, *proto.MinionGetDataSpecInfoRequest) (*proto.MinionGetDataSpecInfoResponse, error) {
	log.Infof("GRPC: GetDataSpecInfo")
	s.mut.Lock()
	defer s.mut.Unlock()

	var info *proto.MinionGetDataSpecInfoResponse

	if s.data == nil {
		info = &proto.MinionGetDataSpecInfoResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "No data spec loaded",
			},
		}
	} else {
		info = &proto.MinionGetDataSpecInfoResponse{
			Status: &proto.GeneralStatus{
				IsOk:          true,
				FailureReason: "",
			},
		}
	}
	return info, nil
}

func (s *MinionServer) GetDBConnectionInfo(context.Context, *proto.MinionGetDBConnectionInfoRequest) (*proto.MinionGetDBConnectionInfoResponse, error) {
	log.Infof("GRPC: CheckDBConnection")
	s.mut.Lock()
	defer s.mut.Unlock()

	// Do we even have a connection?
	if s.db == nil {
		return &proto.MinionGetDBConnectionInfoResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "no connection info",
			},
			DbSpec: &proto.DBSpec{
				Driver: "",
				Url:    "",
			},
		}, nil
	}

	// Connection exists - check if it works
	err := work.ConnCheck(s.db.db)
	if err != nil {
		return &proto.MinionGetDBConnectionInfoResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "connection check failed",
			},
			DbSpec: &proto.DBSpec{
				Driver: s.db.driver,
				Url:    s.db.url,
			},
		}, nil
	}

	// Connection exists and check succeeded
	return &proto.MinionGetDBConnectionInfoResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
		DbSpec: &proto.DBSpec{
			Driver: s.db.driver,
			Url:    s.db.url,
		},
	}, nil
}

func (s *MinionServer) GetWorkloadInfo(_ context.Context, in *proto.MinionGetWorkloadInfoRequest) (*proto.MinionGetWorkloadInfoResponse, error) {
	log.Infof("GRPC: GetWorkloadInfo")
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.workload == nil {
		return &proto.MinionGetWorkloadInfoResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "no workload info",
			},
		}, nil
	} else {
		return &proto.MinionGetWorkloadInfoResponse{
			Status: &proto.GeneralStatus{
				IsOk:          true,
				FailureReason: "",
			},
			WorkloadSpec: &proto.WorkloadSpec{
				WorkloadName:  s.workload.workloadName,
				AssignedRange: proto.RangeToProto(s.workload.assignedRange),
				TableName:     s.workload.tableName,
				DurationSec:   int32(s.workload.durationSec),
				Concurrency:   int32(s.workload.concurrency),
				BatchSize:     int32(s.workload.batchSize),
			},
			IsRunning: s.workload.workload.IsRunning(),
		}, nil
	}
}

func (s MinionState) String() string {
	switch s {
	case Idle:
		return "idle"
	case Prepared:
		return "prepared"
	case Running:
		return "running"
	default:
		panic(fmt.Errorf("unknown minion state: %d", s))
	}
}
