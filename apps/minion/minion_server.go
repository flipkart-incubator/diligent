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

	host        string
	grpcPort    string
	metricsPort string

	mut      *sync.Mutex
	data     *DataContext
	db       *DBContext
	workload *WorkloadContext
	metrics  *metrics.DiligentMetrics
}

func NewMinionServer(host string, grpcPort, metricsPort string) *MinionServer {
	return &MinionServer{
		host:        host,
		grpcPort:    grpcPort,
		metricsPort: metricsPort,

		mut:      &sync.Mutex{},
		data:     nil,
		db:       nil,
		workload: nil,
		metrics:  nil,
	}
}

func (s *MinionServer) Serve() error {
	// Start metrics serving
	s.metrics = metrics.NewDiligentMetrics(s.metricsPort)
	s.metrics.Register()

	// Create listening port
	address := s.host + s.grpcPort
	listener, err := net.Listen("tcp", address)
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

func (s *MinionServer) Ping(_ context.Context, in *proto.PingRequest) (*proto.PingResponse, error) {
	log.Infof("GRPC: Ping")
	s.mut.Lock()
	defer s.mut.Unlock()

	return &proto.PingResponse{Nonce: in.GetNonce()}, nil
}

func (s *MinionServer) LoadDataSpec(ctx context.Context, in *proto.LoadDataSpecRequest) (*proto.LoadDataSpecResponse, error) {
	log.Infof("GRPC: LoadDataSpec(name=%s)", in.GetSpecName())
	s.mut.Lock()
	defer s.mut.Unlock()

	// Don't process if workload is running
	if s.workload != nil && s.workload.workload.IsRunning() {
		return &proto.LoadDataSpecResponse{
			IsOk:          false,
			FailureReason: "Server is busy running workload",
		}, nil
	}

	dataSpec := proto.DataSpecFromProto(in.GetDataSpec())
	s.data = &DataContext{
		dataSpecName: in.GetSpecName(),
		dataSpec:     dataSpec,
		dataGen:      datagen.NewDataGen(dataSpec),
	}

	log.Tracef("%s\n", string(dataSpec.Json()))

	return &proto.LoadDataSpecResponse{
		IsOk:          true,
		FailureReason: "",
		DataSpecInfo: &proto.DataSpecInfo{
			SpecName:   in.GetSpecName(),
			SpecType:   dataSpec.SpecType,
			Version:    int32(dataSpec.Version),
			NumRecs:    int32(dataSpec.KeyGenSpec.NumKeys()),
			RecordSize: int32(dataSpec.RecordSize),
			Hash:       0,
		},
	}, nil
}

func (s *MinionServer) GetDataSpecInfo(context.Context, *proto.GetDataSpecInfoRequest) (*proto.GetDataSpecInfoResponse, error) {
	log.Infof("GRPC: GetDataSpecInfo")
	s.mut.Lock()
	defer s.mut.Unlock()

	var info *proto.GetDataSpecInfoResponse

	if s.data == nil {
		info = &proto.GetDataSpecInfoResponse{
			IsOk:          false,
			FailureReason: "No data spec loaded",
		}
	} else {
		info = &proto.GetDataSpecInfoResponse{
			IsOk:          true,
			FailureReason: "",
			DataSpecInfo: &proto.DataSpecInfo{
				SpecName:   s.data.dataSpecName,
				SpecType:   s.data.dataSpec.SpecType,
				Version:    int32(s.data.dataSpec.Version),
				NumRecs:    int32(s.data.dataSpec.KeyGenSpec.NumKeys()),
				RecordSize: int32(s.data.dataSpec.RecordSize),
				Hash:       0,
			},
		}
	}
	return info, nil
}

func (s *MinionServer) OpenDBConnection(_ context.Context, in *proto.OpenDBConnectionRequest) (*proto.OpenDBConnectionResponse, error) {
	log.Infof("GRPC: OpenDBConnection(%s, %s)", in.GetDriver(), in.GetUrl())
	s.mut.Lock()
	defer s.mut.Unlock()

	// Don't process if workload is running
	if s.workload != nil && s.workload.workload.IsRunning() {
		return &proto.OpenDBConnectionResponse{IsOk: false, FailureReason: "Server is busy running workload"}, nil
	}

	// Validate driver
	driver := in.GetDriver()
	switch driver {
	case "mysql", "pgx":
	default:
		errStr := fmt.Sprintf("invalid driver: '%s'. Allowed values are 'mysql', 'pgx'", driver)
		return &proto.OpenDBConnectionResponse{IsOk: false, FailureReason: errStr}, fmt.Errorf(errStr)
	}

	// Validate URL
	url := in.GetUrl()
	if url == "" {
		errStr := fmt.Sprintf("please specify the connection url")
		return &proto.OpenDBConnectionResponse{IsOk: false, FailureReason: errStr}, fmt.Errorf(errStr)
	}

	// Close any existing connection
	if s.db != nil {
		s.db.db.Close()
		s.db = nil
	}

	// Open new connection
	db, err := sql.Open(driver, url)
	if err != nil {
		return &proto.OpenDBConnectionResponse{IsOk: false, FailureReason: err.Error()}, nil
	}
	s.db = &DBContext{
		driver: driver,
		url:    url,
		db:     db,
	}

	err = work.ConnCheck(s.db.db)
	if err != nil {
		return &proto.OpenDBConnectionResponse{IsOk: false, FailureReason: err.Error()}, nil
	}
	log.Infof("DB Connection successful (driver=%s, url=%s)", driver, url)

	// Respond to RPC with success
	return &proto.OpenDBConnectionResponse{IsOk: true, FailureReason: ""}, nil
}

func (s *MinionServer) GetDBConnectionInfo(context.Context, *proto.GetDBConnectionInfoRequest) (*proto.GetDBConnectionInfoResponse, error) {
	log.Infof("GRPC: CheckDBConnection")
	s.mut.Lock()
	defer s.mut.Unlock()

	// Do we even have a connection?
	if s.db == nil {
		return &proto.GetDBConnectionInfoResponse{
			IsOk:          false,
			FailureReason: "no connection info",
			Driver:        "",
			Url:           "",
		}, nil
	}

	// Connection exists - check if it works
	err := work.ConnCheck(s.db.db)
	if err != nil {
		return &proto.GetDBConnectionInfoResponse{
			IsOk:          false,
			FailureReason: "connection check failed",
			Driver:        s.db.driver,
			Url:           s.db.url,
		}, nil
	}

	// Connection exists and check succeeded
	return &proto.GetDBConnectionInfoResponse{
		IsOk:          true,
		FailureReason: "",
		Driver:        s.db.driver,
		Url:           s.db.url,
	}, nil
}

func (s *MinionServer) CloseDBConnection(context.Context, *proto.CloseDBConnectionRequest) (*proto.CloseDBConnectionResponse, error) {
	log.Infof("GRPC: CheckDBConnection")
	s.mut.Lock()
	defer s.mut.Unlock()

	// Close any existing connection
	if s.db != nil {
		s.db.db.Close()
		s.db.db = nil
		s.db.driver = ""
		s.db.url = ""
	}

	return &proto.CloseDBConnectionResponse{}, nil
}

func (s *MinionServer) RunWorkload(_ context.Context, in *proto.RunWorkloadRequest) (*proto.RunWorkloadResponse, error) {
	log.Infof("GRPC: RunWorkload")
	s.mut.Lock()
	defer s.mut.Unlock()

	// Don't process if workload is running
	if s.workload != nil && s.workload.workload.IsRunning() {
		return &proto.RunWorkloadResponse{IsOk: false, FailureReason: "Server is busy running workload"}, nil
	}

	// Validate whether we have a dataspec and DB connection
	if s.data == nil {
		return &proto.RunWorkloadResponse{
			IsOk:          false,
			FailureReason: fmt.Sprintf("no dataspec loaded"),
		}, nil
	}
	if s.db == nil {
		return &proto.RunWorkloadResponse{
			IsOk:          false,
			FailureReason: fmt.Sprintf("no connection to database"),
		}, nil
	}

	// Validate request params
	durationSec := int(in.GetDurationSec())
	if durationSec < 0 {
		reason := fmt.Sprintf("invalid duration %d. Must be >= 0", durationSec)
		return &proto.RunWorkloadResponse{
			IsOk:          false,
			FailureReason: reason,
		}, fmt.Errorf(reason)
	}

	batchSize := int(in.GetBatchSize())
	if batchSize < 1 {
		reason := fmt.Sprintf("invalid batch size %d. Must be >= 1", batchSize)
		return &proto.RunWorkloadResponse{
			IsOk:          false,
			FailureReason: reason,
		}, fmt.Errorf(reason)
	}

	concurrency := int(in.GetConcurrency())
	if concurrency < 1 {
		reason := fmt.Sprintf("invalid concurrency %d. Must be >= 1", concurrency)
		return &proto.RunWorkloadResponse{
			IsOk:          false,
			FailureReason: reason,
		}, fmt.Errorf(reason)
	}

	table := in.GetTableName()
	if table == "" {
		reason := fmt.Sprintf("table name is missing")
		return &proto.RunWorkloadResponse{
			IsOk:          false,
			FailureReason: reason,
		}, fmt.Errorf(reason)
	}

	inputRange := in.GetAssignedRange()
	if inputRange == nil {
		reason := fmt.Sprintf("assigned range is missing")
		return &proto.RunWorkloadResponse{
			IsOk:          false,
			FailureReason: reason,
		}, fmt.Errorf(reason)
	}

	totalRows := s.data.dataSpec.KeyGenSpec.NumKeys()
	rangeStart := int(inputRange.GetStart())
	rangeLimit := int(inputRange.GetLimit())
	if rangeStart < 0 || rangeLimit < rangeStart {
		reason := fmt.Sprintf("invalid range [%d, %d)", rangeStart, rangeLimit)
		return &proto.RunWorkloadResponse{
			IsOk:          false,
			FailureReason: reason,
		}, fmt.Errorf(reason)
	}
	if rangeStart < 0 || rangeStart >= totalRows {
		reason := fmt.Sprintf("invalid range [%d, %d)", rangeStart, rangeLimit)
		return &proto.RunWorkloadResponse{
			IsOk:          false,
			FailureReason: reason,
		}, fmt.Errorf(reason)
	}
	if rangeLimit < 0 || rangeLimit > totalRows {
		reason := fmt.Sprintf("invalid range [%d, %d)", rangeStart, rangeLimit)
		return &proto.RunWorkloadResponse{
			IsOk:          false,
			FailureReason: reason,
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

	workloadName := in.GetWorkloadName()
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
		return &proto.RunWorkloadResponse{
			IsOk:          false,
			FailureReason: reason,
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

	return &proto.RunWorkloadResponse{
		IsOk:          true,
		FailureReason: "",
	}, nil
}

func (s *MinionServer) GetWorkloadInfo(_ context.Context, in *proto.GetWorkloadInfoRequest) (*proto.GetWorkloadInfoResponse, error) {
	log.Infof("GRPC: GetWorkloadInfo")
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.workload == nil {
		return &proto.GetWorkloadInfoResponse{
			IsOk:          false,
			FailureReason: "no workload info",
		}, nil
	} else {
		return &proto.GetWorkloadInfoResponse{
			IsOk:          true,
			FailureReason: "",
			WorkloadInfo: &proto.WorkloadInfo{
				WorkloadName:  s.workload.workloadName,
				AssignedRange: proto.RangeToProto(s.workload.assignedRange),
				TableName:     s.workload.tableName,
				DurationSec:   int32(s.workload.durationSec),
				Concurrency:   int32(s.workload.concurrency),
				BatchSize:     int32(s.workload.batchSize),
				IsRunning:     s.workload.workload.IsRunning(),
			},
		}, nil
	}
}

func (s *MinionServer) StopWorkload(_ context.Context, in *proto.StopWorkloadRequest) (*proto.StopWorkloadResponse, error) {
	log.Infof("GRPC: StopWorkload")
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.workload.workload.IsRunning() {
		s.workload.workload.Cancel()
		s.workload.workload.Wait()
	}
	return &proto.StopWorkloadResponse{}, nil
}
