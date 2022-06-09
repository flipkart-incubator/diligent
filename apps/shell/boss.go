package main

import (
	"context"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/flipkart-incubator/diligent/pkg/datagen"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	"google.golang.org/grpc"
	"strings"
	"time"
)

func init() {
	bossCmd := &grumble.Command{
		Name:    "boss",
		Help:    "work with the boss",
		Aliases: []string{"bs"},
	}
	grumbleShell.AddCommand(bossCmd)

	bossPingCmd := &grumble.Command{
		Name: "ping",
		Help: "ping the boss",
		Flags: func(f *grumble.Flags) {
			f.String("b", "boss-url", "", "URL of boss server")
		},
		Run: bossPing,
	}
	bossCmd.AddCommand(bossPingCmd)

	bossMinionRegisterCmd := &grumble.Command{
		Name: "register-minion",
		Help: "register a minion with the boss",
		Flags: func(f *grumble.Flags) {
			f.String("b", "boss-url", "", "URL of boss server")
			f.String("m", "minion-url", "", "URL of minion server")
		},
		Run: bossRegisterMinion,
	}
	bossCmd.AddCommand(bossMinionRegisterCmd)

	bossMinionUnregisterCmd := &grumble.Command{
		Name: "unregister-minion",
		Help: "unregister a minion with the boss",
		Flags: func(f *grumble.Flags) {
			f.String("b", "boss-url", "", "URL of boss server")
			f.String("m", "minion-url", "", "URL of minion server")
		},
		Run: bossUnregisterMinion,
	}
	bossCmd.AddCommand(bossMinionUnregisterCmd)

	bossMinionShowCmd := &grumble.Command{
		Name: "show-minions",
		Help: "show minions registered with the boss",
		Flags: func(f *grumble.Flags) {
			f.String("b", "boss-url", "", "URL of boss server")
		},
		Run: bossShowMinions,
	}
	bossCmd.AddCommand(bossMinionShowCmd)

	bossRunWorkloadCmd := &grumble.Command{
		Name: "run-workload",
		Help: "Run a workload",
		Flags: func(f *grumble.Flags) {
			f.String("b", "boss-url", "", "URL of boss server")
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
		Flags: func(f *grumble.Flags) {
			f.String("b", "boss-url", "", "URL of boss server")
		},
		Run: bossStopWorkload,
	}
	bossCmd.AddCommand(bossStopWorkloadCmd)
}

func bossPing(c *grumble.Context) error {
	bossUrl := c.Flags.String("boss-url")
	if bossUrl == "" {
		return fmt.Errorf("please provide a valid bossUrl for the boss server")
	}

	bossClient, err := getBossClient(bossUrl)
	if err != nil {
		c.App.Printf("Request failed: (%v)\n", err)
		return err
	}

	grpcCtx, grpcCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = bossClient.Ping(grpcCtx, &proto.BossPingRequest{})
	grpcCtxCancel()
	if err != nil {
		c.App.Printf("%s: Ping failed! (%v)\n", bossUrl, err)
	} else {
		c.App.Printf("OK\n")
	}
	return nil
}

func bossRegisterMinion(c *grumble.Context) error {
	bossUrl := c.Flags.String("boss-url")
	if bossUrl == "" {
		return fmt.Errorf("please provide a valid bossUrl for the boss server")
	}

	minionUrl := c.Flags.String("minion-url")
	if minionUrl == "" {
		return fmt.Errorf("please provide a valid minionUrl for the minion server")
	}

	bossClient, err := getBossClient(bossUrl)
	if err != nil {
		c.App.Printf("Request failed: (%v)\n", err)
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = bossClient.RegisterMinion(grpcCtx, &proto.BossRegisterMinionRequest{Url: minionUrl})
	grpcCancel()
	if err != nil {
		c.App.Printf("Registration failed: (%v)\n", err)
	} else {
		c.App.Printf("OK\n")
	}
	return nil
}

func bossUnregisterMinion(c *grumble.Context) error {
	bossUrl := c.Flags.String("boss-url")
	if bossUrl == "" {
		return fmt.Errorf("please provide a valid bossUrl for the boss server")
	}

	minionUrl := c.Flags.String("minion-url")
	if minionUrl == "" {
		return fmt.Errorf("please provide a valid minionUrl for the minion server")
	}

	bossClient, err := getBossClient(bossUrl)
	if err != nil {
		c.App.Printf("Request failed: (%v)\n", err)
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = bossClient.UnregisterMinion(grpcCtx, &proto.BossUnregisterMinonRequest{Url: minionUrl})
	grpcCancel()
	if err != nil {
		c.App.Printf("Unregistration failed: (%v)\n", err)
	} else {
		c.App.Printf("OK\n")
	}
	return nil
}

func bossShowMinions(c *grumble.Context) error {
	bossUrl := c.Flags.String("boss-url")
	if bossUrl == "" {
		return fmt.Errorf("please provide a valid bossUrl for the boss server")
	}

	bossClient, err := getBossClient(bossUrl)
	if err != nil {
		c.App.Printf("Request failed: (%v)\n", err)
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), 20*time.Second)
	res, err := bossClient.ShowMinions(grpcCtx, &proto.BossShowMinionRequest{})
	grpcCancel()
	if err != nil {
		c.App.Printf("Request failed: (%v)\n", err)
		return err
	}

	for _, minionStatus := range res.Minions {
		c.App.Printf("%s : reachable=%v [reason=%s]\n",
			minionStatus.GetUrl(), minionStatus.GetStatus().GetIsOk(), minionStatus.GetStatus().GetFailureReason())
	}
	return nil
}

func bossRunWorkload(c *grumble.Context) error {
	// Boss URL param
	bossUrl := c.Flags.String("boss-url")
	if bossUrl == "" {
		return fmt.Errorf("please provide a valid bossUrl for the boss server")
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

	bossClient, err := getBossClient(bossUrl)
	if err != nil {
		c.App.Printf("Request failed: (%v)\n", err)
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), 15*time.Second)
	res, err := bossClient.RunWorkload(grpcCtx, &proto.BossRunWorkloadRequest{
		DataSpec: proto.DataSpecToProto(dataSpec),
		DbSpec:   &proto.DBSpec{
			Driver: dbDriver,
			Url:    dbUrl,
		},
		WlSpec:   &proto.WorkloadSpec{
			WorkloadName:  workloadName,
			AssignedRange: nil,
			TableName:     table,
			DurationSec:   int32(durationSec),
			Concurrency:   int32(concurrency),
			BatchSize:     int32(batchSize),
		},
	})
	grpcCancel()
	if err != nil {
		c.App.Printf("gRPC Request failed: (%v)\n", err)
		return err
	}
	if res.GetStatus().GetIsOk() != true {
		c.App.Printf("Request failed: (%v)\n", res.GetStatus().GetFailureReason())
		return err
	}
	c.App.Printf("OK\n")
	return nil
}

func bossStopWorkload(c *grumble.Context) error {
	bossUrl := c.Flags.String("boss-url")
	if bossUrl == "" {
		return fmt.Errorf("please provide a valid bossUrl for the boss server")
	}

	bossClient, err := getBossClient(bossUrl)
	if err != nil {
		c.App.Printf("Request failed: (%v)\n", err)
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := bossClient.StopWorkload(grpcCtx, &proto.BossStopWorkloadRequest{})
	grpcCancel()
	if err != nil {
		c.App.Printf("Request failed: (%v)\n", err)
		return err
	}
	if res.GetStatus().GetIsOk() != true {
		c.App.Printf("Request failed: (%v)\n", res.GetStatus().GetFailureReason())
		return err
	}
	c.App.Printf("OK")
	return nil
}

func getBossClient(bossUrl string) (proto.BossClient, error) {
	connCtx, connCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(connCtx, bossUrl, grpc.WithInsecure(), grpc.WithBlock())
	connCtxCancel()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s (%v)", bossUrl, err)
	}
	bossClient := proto.NewBossClient(conn)
	return bossClient, nil
}