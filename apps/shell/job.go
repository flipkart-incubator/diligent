package main

import (
	"context"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/flipkart-incubator/diligent/pkg/datagen"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	"strings"
	"time"
)

func init() {
	jobCmd := &grumble.Command{
		Name:    "job",
		Help:    "work with the jobs",
		Aliases: []string{"jo"},
	}
	grumbleApp.AddCommand(jobCmd)

	jobPrepareCmd := &grumble.Command{
		Name: "prepare",
		Help: "Prepare to run a job",
		Flags: func(f *grumble.Flags) {
			f.String("n", "name", "", "name of the job")
			f.String("s", "dataspec-file", "", "name of the dataspec file")
			f.String("r", "db-driver", "", "db driver to use")
			f.String("d", "db-url", "", "db connection url")
			f.Int("t", "duration", 0, "duration after which workload is terminated (seconds). zero for no timout (default)")
			f.Int("c", "concurrency", 1, "number of concurrent workers")
			f.Int("k", "batch-size", 1, "number of statements in a transaction (for transaction based workloads)")
		},
		Args: func(a *grumble.Args) {
			a.String("workload", "name of workload to run [insert,insert-txn,select,select-txn,update,update-txn,delete,delete-txn")
			a.String("table", "name of the table to run the workload on")
		},
		Run: jobPrepare,
	}
	jobCmd.AddCommand(jobPrepareCmd)

	jobRunCmd := &grumble.Command{
		Name: "run",
		Help: "Start the execution of the current job",
		Run:  jobRun,
	}
	jobCmd.AddCommand(jobRunCmd)

	jobAbortCmd := &grumble.Command{
		Name: "abort",
		Help: "Abort the execution of the current job",
		Run:  jobAbort,
	}
	jobCmd.AddCommand(jobAbortCmd)

	jobInfoCmd := &grumble.Command{
		Name: "info",
		Help: "Show current job information",
		Run:  jobInfo,
	}
	jobCmd.AddCommand(jobInfoCmd)

	jobAwaitCmd := &grumble.Command{
		Name: "await-completion",
		Help: "wait for a job to complete",
		Flags: func(f *grumble.Flags) {
			f.Duration("t", "timeout", 1*time.Hour, "wait timeout")
		},
		Run: jobAwaitCompletion,
	}
	jobCmd.AddCommand(jobAwaitCmd)
}

func jobPrepare(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")

	// Job name param
	jobName := c.Flags.String("name")
	if jobName == "" {
		return fmt.Errorf("please specify a name for the job")
	}

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
	default:
		return fmt.Errorf("invalid driver: '%s'. Allowed values are 'mysql', 'pgx'", dbDriver)
	}

	// DB url param
	dbUrl := c.Flags.String("db-url")
	if dbUrl == "" {
		return fmt.Errorf("please specify the connection url")
	}

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

	// Load dataSpec
	c.App.Println("Loading data spec from file:", dataspecFileName)
	dataSpec, err := datagen.LoadSpecFromFile(dataspecFileName)
	if err != nil {
		return err
	}

	c.App.Println("Preparing job:")
	c.App.Println("DataSpec:")
	c.App.Println("\tfile:", dataspecFileName)
	c.App.Println("\tnumRows:", dataSpec.KeyGenSpec.NumKeys())
	c.App.Println("\trecordSize:", dataSpec.RecordSize)
	c.App.Println("DB:")
	c.App.Println("\tdriver:", dbDriver)
	c.App.Println("\turl:", dbUrl)
	c.App.Println("Workload:")
	c.App.Println("\tworkload:", workloadName)
	c.App.Println("\ttable:", table)
	c.App.Println("\tbatchSize:", batchSize)
	c.App.Println("\tconcurrency:", concurrency)
	c.App.Println("\tduration(s):", durationSec)

	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	res, err := bossClient.PrepareJob(grpcCtx, &proto.BossPrepareJobRequest{
		JobSpec: &proto.JobSpec{
			JobName:  jobName,
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
	return nil
}

func jobRun(c *grumble.Context) error {
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

func jobAbort(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	res, err := bossClient.AbortJob(grpcCtx, &proto.BossAbortJobRequest{})
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

func jobInfo(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	res, err := bossClient.GetJobInfo(grpcCtx, &proto.BossGetJobInfoRequest{})
	reqDuration := time.Since(reqStart)
	grpcCancel()

	if err != nil {
		c.App.Printf("Request failed [elapsed=%v]\n", reqDuration)
		return err
	}
	if !res.GetStatus().GetIsOk() {
		c.App.Printf("Request failed [elapsed=%v]\n", reqDuration)
		return fmt.Errorf(res.GetStatus().GetFailureReason())
	}

	c.App.Printf("OK [elapsed=%v]\n", reqDuration)
	showBossJobInfo(c, res.GetJobInfo())
	return nil
}

func showBossJobInfo(c *grumble.Context, ji *proto.BossJobInfo) {
	if ji == nil {
		c.App.Printf("No current job\n")
		return
	}
	ds := proto.DataSpecFromProto(ji.GetJobSpec().GetDataSpec())
	c.App.Printf("\n")
	c.App.Printf("job-name:  %s\n", ji.GetJobSpec().GetJobName())
	c.App.Printf("job-state: %s\n", ji.GetJobState())
	c.App.Printf("prepare-time: %s\n", time.UnixMilli(ji.GetPrepareTime()).Format(time.UnixDate))
	c.App.Printf("run-time:     %s\n", time.UnixMilli(ji.GetRunTime()).Format(time.UnixDate))
	c.App.Printf("end-time:     %s\n", time.UnixMilli(ji.GetEndTime()).Format(time.UnixDate))
	c.App.Printf("data-spec:\n")
	c.App.Printf("\tdata-num-recs: %d\n", ds.KeyGenSpec.NumKeys())
	c.App.Printf("\tdata-rec-size: %d\n", ji.GetJobSpec().GetDataSpec().GetRecordSize())
	c.App.Printf("db-spec:\n")
	c.App.Printf("\tdb-driver: %s\n", ji.GetJobSpec().GetDbSpec().GetDriver())
	c.App.Printf("\tdb-url:    %s\n", ji.GetJobSpec().GetDbSpec().GetUrl())
	c.App.Printf("workload-spec:\n")
	c.App.Printf("\tworkload-name:           %s\n", ji.GetJobSpec().GetWorkloadSpec().GetWorkloadName())
	c.App.Printf("\ttable-name:              %s\n", ji.GetJobSpec().GetWorkloadSpec().GetTableName())
	c.App.Printf("\tworkload-assigned-range: %s\n", ji.GetJobSpec().GetWorkloadSpec().GetAssignedRange())
	c.App.Printf("\tworkload-batch-size:     %d\n", ji.GetJobSpec().GetWorkloadSpec().GetBatchSize())
	c.App.Printf("\tworkload-concurrency:    %d\n", ji.GetJobSpec().GetWorkloadSpec().GetConcurrency())
	c.App.Printf("\tworkload-duration-sec:   %d\n", ji.GetJobSpec().GetWorkloadSpec().GetDurationSec())
	c.App.Printf("minions:\n")
	for _, ma := range ji.GetMinionAddrs() {
		c.App.Printf("\t%s\n", ma)
	}
}

func jobAwaitCompletion(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	timeout := c.Flags.Duration("timeout")

	c.App.Printf("Waiting for current job to end. Wait timeout=%s\n", timeout.String())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
		res, err := bossClient.GetJobInfo(grpcCtx, &proto.BossGetJobInfoRequest{})
		grpcCancel()
		if err != nil {
			c.App.Printf("Request to boss failed (%s)\n", err.Error())
			continue
		}
		switch res.GetJobInfo().GetJobState() {
		case proto.JobState_NEW, proto.JobState_PREPARED, proto.JobState_RUNNING:
			c.App.Printf(".")
		case proto.JobState_ENDED_SUCCESS:
			c.App.Printf("\nJob has ended successfully\n")
			return nil
		case proto.JobState_ENDED_FAILURE:
			c.App.Printf("\nJob has failed\n")
			return nil
		case proto.JobState_ENDED_ABORTED:
			c.App.Printf("\nJob was aborted\n")
			return nil
		}
		if ctx.Err() != nil {
			return fmt.Errorf("timeout occurred. job has not ended yet")
		}
		time.Sleep(1 * time.Second)
	}
}
