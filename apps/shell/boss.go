package main

import (
	"context"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/flipkart-incubator/diligent/pkg/datagen"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	"google.golang.org/grpc"
	"net"
	"strings"
	"time"
)

const (
	defaultBossPort           = "5710"
	defaultMinionPort         = "5711"
	bossConnectionTimeoutSecs = 5
	bossRequestTimeoutSecs    = 5
)

func init() {
	bossCmd := &grumble.Command{
		Name:    "boss",
		Help:    "work with the boss",
		Aliases: []string{"bs"},
	}
	grumbleApp.AddCommand(bossCmd)

	bossPingCmd := &grumble.Command{
		Name: "ping",
		Help: "ping the boss",
		Run:  bossPing,
	}
	bossCmd.AddCommand(bossPingCmd)

	bossMinionRegisterCmd := &grumble.Command{
		Name: "register-minion",
		Help: "manually register a minion with the boss",
		Args: func(a *grumble.Args) {
			a.String("minion-addr", "host[:port] of minion server")
		},
		Run: bossRegisterMinion,
	}
	bossCmd.AddCommand(bossMinionRegisterCmd)

	bossMinionUnregisterCmd := &grumble.Command{
		Name: "unregister-minion",
		Help: "manually unregister a minion with the boss",
		Args: func(a *grumble.Args) {
			a.String("minion-addr", "host[:port] of minion server")
		},
		Run: bossUnregisterMinion,
	}
	bossCmd.AddCommand(bossMinionUnregisterCmd)

	bossMinionShowCmd := &grumble.Command{
		Name: "show-minions",
		Help: "show minions registered with the boss",
		Flags: func(f *grumble.Flags) {
			f.Bool("b", "build-info", false, "show build information")
			f.Bool("p", "process-info", false, "show process information")
			f.Bool("j", "job-info", false, "show job information")
		},
		Run: bossShowMinions,
	}
	bossCmd.AddCommand(bossMinionShowCmd)

	bossPrepareJobCmd := &grumble.Command{
		Name: "prepare-job",
		Help: "Prepare to run a job",
		Flags: func(f *grumble.Flags) {
			f.String("s", "dataspec-file", "", "name of the dataspec file")
			f.String("r", "db-driver", "", "db driver to use")
			f.String("d", "db-url", "", "db connection url")
			f.Int("t", "duration", 0, "duration after which workload is terminated (seconds). zero for no timout (default)")
			f.Int("c", "concurrency", 1, "number of concurrent workers")
			f.Int("k", "batch-size", 1, "number of statements in a transaction (for transaction based workloads)")
			f.String("m", "description", "", "a description of the job")
		},
		Args: func(a *grumble.Args) {
			a.String("workload", "name of workload to run [insert,insert-txn,select,select-txn,update,update-txn,delete,delete-txn", grumble.Default(""))
			a.String("table", "name of the table to run the workload on", grumble.Default(""))
		},
		Run: bossPrepareJob,
	}
	bossCmd.AddCommand(bossPrepareJobCmd)

	bossRunJobCmd := &grumble.Command{
		Name: "run-job",
		Help: "Start the execution of the current job",
		Run:  bossRunJob,
	}
	bossCmd.AddCommand(bossRunJobCmd)

	bossStopJobCmd := &grumble.Command{
		Name: "stop-job",
		Help: "Stop the execution of the current job",
		Run:  bossStopJob,
	}
	bossCmd.AddCommand(bossStopJobCmd)
}

func bossPing(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	grpcCtx, grpcCtxCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	res, err := bossClient.Ping(grpcCtx, &proto.BossPingRequest{})
	reqDuration := time.Since(reqStart)
	grpcCtxCancel()
	if err != nil {
		c.App.Printf("Ping failed [elapsed=%v, bossAddr=%s]\n", reqDuration, bossAddr)
		return err
	} else {
		c.App.Printf("OK [elapsed=%v]\n", reqDuration)
		c.App.Printf("AppName   : %s\n", res.GetBuildInfo().GetAppName())
		c.App.Printf("AppVersion: %s\n", res.GetBuildInfo().GetAppVersion())
		c.App.Printf("CommitHash: %s\n", res.GetBuildInfo().GetCommitHash())
		c.App.Printf("GoVersion : %s\n", res.GetBuildInfo().GetGoVersion())
		c.App.Printf("BuildTime : %s\n", res.GetBuildInfo().GetBuildTime())
		c.App.Printf("Pid       : %s\n", res.GetProcessInfo().GetPid())
		c.App.Printf("StartTime : %s\n", res.GetProcessInfo().GetStartTime())
		c.App.Printf("Uptime    : %s\n", res.GetProcessInfo().GetUptime())
	}
	return nil
}

func bossRegisterMinion(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	minionAddr := c.Args.String("minion-addr")
	if minionAddr == "" {
		return fmt.Errorf("please provide a valid minion-addr for the minion server")
	}
	if _, _, err := net.SplitHostPort(minionAddr); err != nil {
		minionAddr = net.JoinHostPort(minionAddr, defaultMinionPort)
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	_, err = bossClient.RegisterMinion(grpcCtx, &proto.BossRegisterMinionRequest{Addr: minionAddr})
	reqDuration := time.Since(reqStart)
	grpcCancel()
	if err != nil {
		c.App.Printf("Request failed [elapsed=%v, minionAddr=%s]\n", reqDuration, minionAddr)
		return err
	} else {
		c.App.Printf("OK [elapsed=%v]\n", reqDuration)
	}
	return nil
}

func bossUnregisterMinion(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	minionAddr := c.Args.String("minion-addr")
	if minionAddr == "" {
		return fmt.Errorf("please provide a valid minion-addr for the minion server")
	}
	if _, _, err := net.SplitHostPort(minionAddr); err != nil {
		minionAddr = net.JoinHostPort(minionAddr, defaultMinionPort)
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	_, err = bossClient.UnregisterMinion(grpcCtx, &proto.BossUnregisterMinonRequest{Addr: minionAddr})
	reqDuration := time.Since(reqStart)
	grpcCancel()
	if err != nil {
		c.App.Printf("Request failed [elapsed=%v, minionAddr=%s]\n", reqDuration, minionAddr)
		return err
	} else {
		c.App.Printf("OK [elapsed=%v]\n", reqDuration)
	}
	return nil
}

func bossShowMinions(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	res, err := bossClient.ShowMinions(grpcCtx, &proto.BossShowMinionRequest{})
	reqDuration := time.Since(reqStart)
	grpcCancel()
	if err != nil {
		c.App.Printf("Request failed [elapsed=%v]\n", reqDuration)
		return err
	} else {
		c.App.Printf("OK [elapsed=%v]\n", reqDuration)
	}

	for _, mi := range res.GetMinionInfos() {
		if mi.GetReachability().GetIsOk() {
			c.App.Printf("%s : OK ", mi.GetAddr())
		} else {
			c.App.Printf("%s : Not Reachable [%s]\n", mi.GetAddr(), mi.GetReachability().GetFailureReason())
			continue
		}
		switch {
		case c.Flags.Bool("build-info"):
			showMinionBuildInfo(c, mi)
		case c.Flags.Bool("process-info"):
			showMinionProcessInfo(c, mi)
		case c.Flags.Bool("job-info"):
			showMinionJobInfo(c, mi)
		default:
			showMinionSummaryInfo(c, mi)
		}
	}

	return nil
}

func showMinionBuildInfo(c *grumble.Context, mi *proto.MinionInfo) {
	c.App.Printf("\n")
	c.App.Printf("\tapp-name: %s\n", mi.GetBuildInfo().GetAppName())
	c.App.Printf("\tapp-version: %s\n", mi.GetBuildInfo().GetAppVersion())
	c.App.Printf("\tcommit-hash: %s\n", mi.GetBuildInfo().GetCommitHash())
	c.App.Printf("\tgo-version: %s\n", mi.GetBuildInfo().GetGoVersion())
	c.App.Printf("\tbuild-time: %s\n", mi.GetBuildInfo().GetBuildTime())
}

func showMinionProcessInfo(c *grumble.Context, mi *proto.MinionInfo) {
	c.App.Printf("\n")
	c.App.Printf("\tpid: %s\n", mi.GetProcessInfo().GetPid())
	c.App.Printf("\tstart-time: %s\n", mi.GetProcessInfo().GetStartTime())
	c.App.Printf("\tup-time: %s\n", mi.GetProcessInfo().GetUptime())
}

func showMinionJobInfo(c *grumble.Context, mi *proto.MinionInfo) {
	if mi.GetJobInfo() == nil {
		c.App.Printf("[No job info]\n")
		return
	}
	ds := proto.DataSpecFromProto(mi.GetJobInfo().GetJobSpec().GetDataSpec())
	c.App.Printf("\n")
	c.App.Printf("\tjob-id:    %s\n", mi.GetJobInfo().GetJobId())
	c.App.Printf("\tjob-state: %s\n", mi.GetJobInfo().GetJobState())
	c.App.Printf("\tprepare-time: %s\n", mi.GetJobInfo().GetPrepareTime())
	c.App.Printf("\trun-time:     %s\n", mi.GetJobInfo().GetRunTime())
	c.App.Printf("\tend-time:     %s\n", mi.GetJobInfo().GetEndTime())
	c.App.Printf("\tdata-spec:\n")
	c.App.Printf("\t\tdata-num-recs: %d\n", ds.KeyGenSpec.NumKeys())
	c.App.Printf("\t\tdata-rec-size: %d\n", mi.GetJobInfo().GetJobSpec().GetDataSpec().GetRecordSize())
	c.App.Printf("\tdb-spec:\n")
	c.App.Printf("\t\tdb-driver: %s\n", mi.GetJobInfo().GetJobSpec().GetDbSpec().GetDriver())
	c.App.Printf("\t\tdb-url:    %s\n", mi.GetJobInfo().GetJobSpec().GetDbSpec().GetUrl())
	c.App.Printf("\tworkload-spec:\n")
	c.App.Printf("\t\tworkload-name:           %s\n", mi.GetJobInfo().GetJobSpec().GetWorkloadSpec().GetWorkloadName())
	c.App.Printf("\t\ttable-name:              %s\n", mi.GetJobInfo().GetJobSpec().GetWorkloadSpec().GetTableName())
	c.App.Printf("\t\tworkload-assigned-range: %s\n", mi.GetJobInfo().GetJobSpec().GetWorkloadSpec().GetAssignedRange())
	c.App.Printf("\t\tworkload-batch-size:     %d\n", mi.GetJobInfo().GetJobSpec().GetWorkloadSpec().GetBatchSize())
	c.App.Printf("\t\tworkload-concurrency:    %d\n", mi.GetJobInfo().GetJobSpec().GetWorkloadSpec().GetConcurrency())
	c.App.Printf("\t\tworkload-duration-sec:   %d\n", mi.GetJobInfo().GetJobSpec().GetWorkloadSpec().GetDurationSec())
}

func showMinionSummaryInfo(c *grumble.Context, mi *proto.MinionInfo) {
	jobId := "none"
	jobState := "none"
	if mi.GetJobInfo() != nil {
		jobId = mi.GetJobInfo().GetJobId()
		jobState = mi.GetJobInfo().GetJobState().String()
	}
	c.App.Printf("[version=%s, pid=%s, uptime=%s, jobId=%s, jobState=%s]\n",
		mi.GetBuildInfo().GetAppVersion(), mi.GetProcessInfo().GetPid(),
		mi.GetProcessInfo().GetUptime(), jobId, jobState)
}

func bossPrepareJob(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")

	// Dataspec param
	dataspecFileName := c.Flags.String("dataspec-file")
	if dataspecFileName == "" {
		return fmt.Errorf("please specify a name for the dataspec")
	}
	if !strings.HasSuffix(dataspecFileName, ".json") {
		dataspecFileName = dataspecFileName + ".json"
	}

	// DB driver param
	dbDriver := c.Flags.String("db-driver")
	switch dbDriver {
	case "mysql", "pgx":
		controllerApp.db.driver = dbDriver
	default:
		return fmt.Errorf("invalid driver: '%s'. Allowed values are 'mysql', 'pgx'", dbDriver)
	}

	// DB url param
	dbUrl := c.Flags.String("db-url")
	if dbUrl == "" {
		return fmt.Errorf("please specify the connection url")
	}
	controllerApp.db.url = dbUrl

	durationSec := c.Flags.Int("duration")
	if durationSec < 0 {
		return fmt.Errorf("invalid duration %d. Must be >= 0", durationSec)
	}
	batchSize := c.Flags.Int("batch-size")
	if batchSize < 1 {
		return fmt.Errorf("invalid batch size %d. Must be >= 1", batchSize)
	}
	concurrency := c.Flags.Int("concurrency")
	if concurrency < 1 {
		return fmt.Errorf("invalid concurrency %d. Must be >= 1", concurrency)
	}

	workloadName := c.Args.String("workload")
	switch workloadName {
	case "insert", "insert-txn",
		"select-pk", "select-pk-txn", "select-uk", "select-uk-txn",
		"update", "update-txn",
		"delete", "delete-txn":
	default:
		return fmt.Errorf("invalid workload '%s'", workloadName)
	}

	table := c.Args.String("table")
	if table == "" {
		return fmt.Errorf("please specify the table to run the workload on")
	}

	// Load dataSpec
	c.App.Println("Loading data spec from file:", dataspecFileName)
	dataSpec, err := datagen.LoadSpecFromFile(dataspecFileName)
	if err != nil {
		return err
	}

	// Description
	jobDesc := c.Flags.String("description")

	c.App.Println("Preparing job:")
	c.App.Println("    Description:", jobDesc)
	c.App.Println("    DataSpec:")
	c.App.Println("    	file:", dataspecFileName)
	c.App.Println("    	numRows:", dataSpec.KeyGenSpec.NumKeys())
	c.App.Println("    	recordSize:", dataSpec.RecordSize)
	c.App.Println("    DB:")
	c.App.Println("    	driver:", dbDriver)
	c.App.Println("    	url:", dbUrl)
	c.App.Println("    Workload:")
	c.App.Println("    	workload:", workloadName)
	c.App.Println("    	table:", table)
	c.App.Println("    	batchSize:", batchSize)
	c.App.Println("    	concurrency:", concurrency)
	c.App.Println("    	duration(s):", durationSec)

	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	res, err := bossClient.PrepareJob(grpcCtx, &proto.BossPrepareJobRequest{
		JobDesc: jobDesc,
		JobSpec: &proto.JobSpec{
			DataSpec: proto.DataSpecToProto(dataSpec),
			DbSpec: &proto.DBSpec{
				Driver: dbDriver,
				Url:    dbUrl,
			},
			WorkloadSpec: &proto.WorkloadSpec{
				WorkloadName:  workloadName,
				AssignedRange: nil,
				TableName:     table,
				DurationSec:   int32(durationSec),
				Concurrency:   int32(concurrency),
				BatchSize:     int32(batchSize),
			},
		},
	})
	reqDuration := time.Since(reqStart)
	grpcCancel()
	if err != nil {
		c.App.Printf("Request failed [elapsed=%v]\n", reqDuration)
		return err
	}
	if res.GetStatus().GetIsOk() != true {
		c.App.Printf("Request failed [elapsed=%v]\n", reqDuration)
		for _, ms := range res.GetMinionStatuses() {
			if ms.GetStatus().GetIsOk() {
				c.App.Printf("%s : OK\n", ms.GetAddr())
			} else {
				c.App.Printf("%s : Failed [reason=%s]\n", ms.GetAddr(), ms.GetStatus().GetFailureReason())
			}
		}
		return fmt.Errorf(res.GetStatus().GetFailureReason())
	}
	c.App.Printf("OK [elapsed=%v]\n", reqDuration)
	c.App.Printf("JobId=%s\n", res.GetJobId())
	return nil
}

func bossRunJob(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	res, err := bossClient.RunJob(grpcCtx, &proto.BossRunJobRequest{})
	reqDuration := time.Since(reqStart)
	grpcCancel()
	if err != nil {
		c.App.Printf("Request failed [elapsed=%v]\n", reqDuration)
		return err
	}
	if res.GetStatus().GetIsOk() != true {
		c.App.Printf("Request failed [elapsed=%v]\n", reqDuration)
		for _, ms := range res.GetMinionStatuses() {
			if ms.GetStatus().GetIsOk() {
				c.App.Printf("%s : OK\n", ms.GetAddr())
			} else {
				c.App.Printf("%s : Failed [reason=%s]\n", ms.GetAddr(), ms.GetStatus().GetFailureReason())
			}
		}
		return fmt.Errorf(res.GetStatus().GetFailureReason())
	}
	c.App.Printf("OK [elapsed=%v]\n", reqDuration)
	return nil
}

func bossStopJob(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	res, err := bossClient.StopJob(grpcCtx, &proto.BossStopJobRequest{})
	reqDuration := time.Since(reqStart)
	grpcCancel()
	if err != nil {
		c.App.Printf("Request failed [elapsed=%v]\n", reqDuration)
		return err
	}
	if res.GetStatus().GetIsOk() != true {
		c.App.Printf("Request failed [elapsed=%v]\n", reqDuration)
		for _, ms := range res.GetMinionStatuses() {
			if ms.GetStatus().GetIsOk() {
				c.App.Printf("%s : OK\n", ms.GetAddr())
			} else {
				c.App.Printf("%s : Failed [reason=%s]\n", ms.GetAddr(), ms.GetStatus().GetFailureReason())
			}
		}
		return fmt.Errorf(res.GetStatus().GetFailureReason())
	}
	c.App.Printf("OK [elapsed=%v]\n", reqDuration)
	return nil
}

func getBossClient(bossAddr string) (proto.BossClient, error) {
	if _, _, err := net.SplitHostPort(bossAddr); err != nil {
		bossAddr = net.JoinHostPort(bossAddr, defaultBossPort)
	}

	connCtx, connCtxCancel := context.WithTimeout(context.Background(), bossConnectionTimeoutSecs*time.Second)
	conn, err := grpc.DialContext(connCtx, bossAddr, grpc.WithInsecure(), grpc.WithBlock())
	connCtxCancel()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s (%v)", bossAddr, err)
	}
	bossClient := proto.NewBossClient(conn)
	return bossClient, nil
}
