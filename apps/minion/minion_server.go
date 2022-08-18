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

	grpcAddr      string
	metricsAddr   string
	advertiseAddr string
	bossAddr      string

	mut      *sync.Mutex
	data     *DataContext
	db       *DBContext
	workload *WorkloadContext
	metrics  *metrics.DiligentMetrics
}

func NewMinionServer(grpcAddr, metricsAddr, advertiseAddr, bossAddr string) *MinionServer {
	return &MinionServer{
		grpcAddr:      grpcAddr,
		metricsAddr:   metricsAddr,
		advertiseAddr: advertiseAddr,
		bossAddr:      bossAddr,

		mut:      &sync.Mutex{},
		data:     nil,
		db:       nil,
		workload: nil,
		metrics:  nil,
	}
}

func (s *MinionServer) RegisterWithBoss() error {
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

	return nil
}

func (s *MinionServer) Serve() error {
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

	return &proto.MinionPingResponse{}, nil
}

func (s *MinionServer) LoadDataSpec(ctx context.Context, in *proto.MinionLoadDataSpecRequest) (*proto.MinionLoadDataSpecResponse, error) {
	log.Infof("GRPC: LoadDataSpec(name=%s)", in.GetSpecName())
	s.mut.Lock()
	defer s.mut.Unlock()

	// Don't process if workload is running
	if s.workload != nil && s.workload.workload.IsRunning() {
		return &proto.MinionLoadDataSpecResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "Server is busy running workload",
			},
		}, nil
	}

	dataSpec := proto.DataSpecFromProto(in.GetDataSpec())
	s.data = &DataContext{
		dataSpecName: in.GetSpecName(),
		dataSpec:     dataSpec,
		dataGen:      datagen.NewDataGen(dataSpec),
	}

	log.Tracef("%s\n", string(dataSpec.Json()))

	return &proto.MinionLoadDataSpecResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
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

func (s *MinionServer) OpenDBConnection(_ context.Context, in *proto.MinionOpenDBConnectionRequest) (*proto.MinionOpenDBConnectionResponse, error) {
	log.Infof("GRPC: OpenDBConnection(%s, %s)", in.GetDbSpec().GetDriver(), in.GetDbSpec().GetUrl())
	s.mut.Lock()
	defer s.mut.Unlock()

	// Don't process if workload is running
	if s.workload != nil && s.workload.workload.IsRunning() {
		return &proto.MinionOpenDBConnectionResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "Server is busy running workload",
			},
		}, nil
	}

	// Validate driver
	driver := in.GetDbSpec().GetDriver()
	switch driver {
	case "mysql", "pgx":
	default:
		errStr := fmt.Sprintf("invalid driver: '%s'. Allowed values are 'mysql', 'pgx'", driver)
		return &proto.MinionOpenDBConnectionResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: errStr,
			},
		}, fmt.Errorf(errStr)
	}

	// Validate URL
	url := in.GetDbSpec().GetUrl()
	if url == "" {
		errStr := fmt.Sprintf("please specify the connection url")
		return &proto.MinionOpenDBConnectionResponse{
			Status: &proto.GeneralStatus{
				IsOk: false, FailureReason: errStr,
			},
		}, fmt.Errorf(errStr)
	}

	// Close any existing connection
	if s.db != nil {
		s.db.db.Close()
		s.db = nil
	}

	// Open new connection
	db, err := sql.Open(driver, url)
	if err != nil {
		return &proto.MinionOpenDBConnectionResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}, nil
	}
	s.db = &DBContext{
		driver: driver,
		url:    url,
		db:     db,
	}

	err = work.ConnCheck(s.db.db)
	if err != nil {
		return &proto.MinionOpenDBConnectionResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: err.Error(),
			},
		}, nil
	}
	log.Infof("DB Connection successful (driver=%s, url=%s)", driver, url)

	// Respond to RPC with success
	return &proto.MinionOpenDBConnectionResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
		},
	}, nil
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

func (s *MinionServer) RunWorkload(_ context.Context, in *proto.MinionRunWorkloadRequest) (*proto.MinionRunWorkloadResponse, error) {
	log.Infof("GRPC: RunWorkload")
	s.mut.Lock()
	defer s.mut.Unlock()

	// Don't process if workload is running
	if s.workload != nil && s.workload.workload.IsRunning() {
		return &proto.MinionRunWorkloadResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: "Server is busy running workload",
			},
		}, nil
	}

	// Validate whether we have a dataspec and DB connection
	if s.data == nil {
		return &proto.MinionRunWorkloadResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: fmt.Sprintf("no dataspec loaded"),
			},
		}, nil
	}
	if s.db == nil {
		return &proto.MinionRunWorkloadResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: fmt.Sprintf("no connection to database"),
			},
		}, nil
	}

	// Validate request params
	durationSec := int(in.GetWorkloadSpec().GetDurationSec())
	if durationSec < 0 {
		reason := fmt.Sprintf("invalid duration %d. Must be >= 0", durationSec)
		return &proto.MinionRunWorkloadResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}, fmt.Errorf(reason)
	}

	batchSize := int(in.GetWorkloadSpec().GetBatchSize())
	if batchSize < 1 {
		reason := fmt.Sprintf("invalid batch size %d. Must be >= 1", batchSize)
		return &proto.MinionRunWorkloadResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}, fmt.Errorf(reason)
	}

	concurrency := int(in.GetWorkloadSpec().GetConcurrency())
	if concurrency < 1 {
		reason := fmt.Sprintf("invalid concurrency %d. Must be >= 1", concurrency)
		return &proto.MinionRunWorkloadResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}, fmt.Errorf(reason)
	}

	table := in.GetWorkloadSpec().GetTableName()
	if table == "" {
		reason := fmt.Sprintf("table name is missing")
		return &proto.MinionRunWorkloadResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}, fmt.Errorf(reason)
	}

	inputRange := in.GetWorkloadSpec().GetAssignedRange()
	if inputRange == nil {
		reason := fmt.Sprintf("assigned range is missing")
		return &proto.MinionRunWorkloadResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}, fmt.Errorf(reason)
	}

	totalRows := s.data.dataSpec.KeyGenSpec.NumKeys()
	rangeStart := int(inputRange.GetStart())
	rangeLimit := int(inputRange.GetLimit())
	if rangeStart < 0 || rangeLimit < rangeStart {
		reason := fmt.Sprintf("invalid range [%d, %d)", rangeStart, rangeLimit)
		return &proto.MinionRunWorkloadResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}, fmt.Errorf(reason)
	}
	if rangeStart < 0 || rangeStart >= totalRows {
		reason := fmt.Sprintf("invalid range [%d, %d)", rangeStart, rangeLimit)
		return &proto.MinionRunWorkloadResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}, fmt.Errorf(reason)
	}
	if rangeLimit < 0 || rangeLimit > totalRows {
		reason := fmt.Sprintf("invalid range [%d, %d)", rangeStart, rangeLimit)
		return &proto.MinionRunWorkloadResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}, fmt.Errorf(reason)
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

	workloadName := in.GetWorkloadSpec().GetWorkloadName()
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
		reason := fmt.Sprintf("invalid workload '%s'", workloadName)
		return &proto.MinionRunWorkloadResponse{
			Status: &proto.GeneralStatus{
				IsOk:          false,
				FailureReason: reason,
			},
		}, fmt.Errorf(reason)
	}

	log.Infof("Running workload: %s", workloadName)
	log.Infof("    dataspec: %s", s.data.dataSpecName)
	log.Infof("    totalRows: %d", s.data.dataSpec.KeyGenSpec.NumKeys())
	log.Infof("    assignedRange: %s", assignedRange)
	log.Infof("    recordSize: %d", s.data.dataSpec.RecordSize)
	log.Infof("    driver: %s", s.db.driver)
	log.Infof("    url: %s", s.db.url)
	log.Infof("    table: %s", table)
	log.Infof("    batchSize: %d", batchSize)
	log.Infof("    concurrency: %d", concurrency)
	log.Infof("    duration(s): %d", durationSec)

	s.workload = &WorkloadContext{
		workloadName:  workloadName,
		assignedRange: assignedRange,
		tableName:     table,
		durationSec:   durationSec,
		concurrency:   concurrency,
		batchSize:     batchSize,
		workload:      workload,
	}

	workload.Start(time.Duration(durationSec) * time.Second)

	return &proto.MinionRunWorkloadResponse{
		Status: &proto.GeneralStatus{
			IsOk:          true,
			FailureReason: "",
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

func (s *MinionServer) StopWorkload(_ context.Context, in *proto.MinionStopWorkloadRequest) (*proto.MinionStopWorkloadResponse, error) {
	log.Infof("GRPC: StopWorkload")
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.workload.workload.IsRunning() {
		s.workload.workload.Cancel()
		s.workload.workload.Wait()
	}

	return &proto.MinionStopWorkloadResponse{
		Status: &proto.GeneralStatus{
			IsOk: true,
		},
	}, nil
}
