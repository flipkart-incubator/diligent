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
		Run:  bossShowMinions,
	}
	bossCmd.AddCommand(bossMinionShowCmd)

	bossRunWorkloadCmd := &grumble.Command{
		Name: "run-workload",
		Help: "Run a workload",
		Flags: func(f *grumble.Flags) {
			f.String("s", "dataspec-file", "", "name of the dataspec file")
			f.String("r", "db-driver", "", "db driver to use")
			f.String("d", "db-url", "", "db connection url")
			f.Int("t", "duration", 0, "duration after which workload is terminated (seconds). zero for no timout (default)")
			f.Int("c", "concurrency", 1, "number of concurrent workers")
			f.Int("k", "batch-size", 10, "number of statements in a transaction (for transaction based workloads)")
		},
		Args: func(a *grumble.Args) {
			a.String("workload", "name of workload to run [insert,insert-txn,select,select-txn,update,update-txn,delete,delete-txn", grumble.Default(""))
			a.String("table", "name of the table to run the workload on", grumble.Default(""))
		},
		Run: bossRunWorkload,
	}
	bossCmd.AddCommand(bossRunWorkloadCmd)

	bossStopWorkloadCmd := &grumble.Command{
		Name: "stop-workload",
		Help: "Stop any running workload",
		Run:  bossStopWorkload,
	}
	bossCmd.AddCommand(bossStopWorkloadCmd)
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
		c.App.Printf("Version   : %s\n", res.GetBuildInfo().GetVersion())
		c.App.Printf("CommitHash: %s\n", res.GetBuildInfo().GetCommitHash())
		c.App.Printf("GoVersion : %s\n", res.GetBuildInfo().GetGoVersion())
		c.App.Printf("BuildTime : %s\n", res.GetBuildInfo().GetBuildTime())
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

	c.App.Printf("Reachability status:\n")
	for _, ms := range res.GetMinionStatuses() {
		if ms.GetStatus().GetIsOk() {
			c.App.Printf("\t%s : Reachable\n", ms.GetAddr())
		} else {
			c.App.Printf("\t%s : Not Reachable [reason=%s]\n", ms.GetAddr(), ms.GetStatus().GetFailureReason())
		}
	}

	c.App.Printf("Build info:\n")
	for _, bi := range res.GetMinionBuildInfos() {
		c.App.Printf("%s : version=%s, commit-hash=%s, go-version=%s, build-time=%s\n",
			bi.GetAddr(), bi.GetBuildInfo().GetVersion(), bi.GetBuildInfo().GetCommitHash(),
			bi.GetBuildInfo().GoVersion, bi.GetBuildInfo().GetBuildTime())
	}
	return nil
}

func bossRunWorkload(c *grumble.Context) error {
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

	c.App.Println("Running workload:")
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
	res, err := bossClient.RunWorkload(grpcCtx, &proto.BossRunWorkloadRequest{
		DataSpec: proto.DataSpecToProto(dataSpec),
		DbSpec: &proto.DBSpec{
			Driver: dbDriver,
			Url:    dbUrl,
		},
		WlSpec: &proto.WorkloadSpec{
			WorkloadName:  workloadName,
			AssignedRange: nil,
			TableName:     table,
			DurationSec:   int32(durationSec),
			Concurrency:   int32(concurrency),
			BatchSize:     int32(batchSize),
		},
	})
	reqDuration := time.Since(reqStart)
	grpcCancel()
	if err != nil {
		c.App.Printf("Request failed [elapsed=%v]\n", reqDuration)
		return err
	}
	if res.GetOverallStatus().GetIsOk() != true {
		c.App.Printf("Request failed [elapsed=%v]\n", reqDuration)
		for _, ms := range res.GetMinionStatuses() {
			if ms.GetStatus().GetIsOk() {
				c.App.Printf("%s : OK\n", ms.GetAddr())
			} else {
				c.App.Printf("%s : Failed [reason=%s]\n", ms.GetAddr(), ms.GetStatus().GetFailureReason())
			}
		}
		return fmt.Errorf(res.GetOverallStatus().GetFailureReason())
	}
	c.App.Printf("OK [elapsed=%v]\n", reqDuration)
	return nil
}

func bossStopWorkload(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	res, err := bossClient.StopWorkload(grpcCtx, &proto.BossStopWorkloadRequest{})
	reqDuration := time.Since(reqStart)
	grpcCancel()
	if err != nil {
		c.App.Printf("Request failed [elapsed=%v]\n", reqDuration)
		return err
	}
	if res.GetOverallStatus().GetIsOk() != true {
		c.App.Printf("Request failed [elapsed=%v]\n", reqDuration)
		for _, ms := range res.GetMinionStatuses() {
			if ms.GetStatus().GetIsOk() {
				c.App.Printf("%s : OK\n", ms.GetAddr())
			} else {
				c.App.Printf("%s : Failed [reason=%s]\n", ms.GetAddr(), ms.GetStatus().GetFailureReason())
			}
		}
		return fmt.Errorf(res.GetOverallStatus().GetFailureReason())
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
