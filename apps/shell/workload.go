package main
//
//import (
//	"context"
//	"fmt"
//	"github.com/desertbit/grumble"
//	"github.com/flipkart-incubator/diligent/pkg/intgen"
//	"github.com/flipkart-incubator/diligent/pkg/proto"
//	"time"
//)
//
//func init() {
//	wlCmd := &grumble.Command{
//		Name:    "workload",
//		Help:    "run workloads",
//		Aliases: []string{"wl"},
//	}
//	grumbleShell.AddCommand(wlCmd)
//
//	wlRunCmd := &grumble.Command{
//		Name: "run",
//		Help: "start running a workload",
//		Flags: func(f *grumble.Flags) {
//			f.Int("d", "duration", 0, "duration after which workload is terminated (seconds)")
//			f.Int("c", "concurrency", 1, "number of concurrent workers")
//			f.Int("b", "batch-size", 10, "number of statements in a transaction (for transaction based workloads)")
//		},
//		Args: func(a *grumble.Args) {
//			a.String("workload", "name of workload to run [insert,insert-txn,select,select-txn,update,update-txn,delete,delete-txn", grumble.Default(""))
//			a.String("table", "name of the table to run the workload on", grumble.Default(""))
//		},
//		Run: wlRun,
//	}
//	wlCmd.AddCommand(wlRunCmd)
//
//	wlShowCmd := &grumble.Command{
//		Name:    "show",
//		Help:    "get info on running workloads",
//		Aliases: []string{"sh"},
//		Run:     wlShow,
//	}
//	wlCmd.AddCommand(wlShowCmd)
//
//	wlStopCmd := &grumble.Command{
//		Name: "stop",
//		Help: "stop any running workloads",
//		Run: wlStop,
//	}
//	wlCmd.AddCommand(wlStopCmd)
//}
//
//func wlRun(c *grumble.Context) error {
//	if controllerApp.data.dataSpec == nil {
//		return fmt.Errorf("no dataspec loaded. Use the 'dataspec load' command to load one")
//	}
//	if controllerApp.db.driver == "" {
//		return fmt.Errorf("no driver specified. Use the 'db conn' command to specify")
//	}
//	if controllerApp.db.url == "" {
//		return fmt.Errorf("no url specified. Use the 'db conn' command to specify")
//	}
//	durationSec := c.Flags.Int("duration")
//	if durationSec < 0 {
//		return fmt.Errorf("invalid duration %d. Must be >= 0", durationSec)
//	}
//	batchSize := c.Flags.Int("batch-size")
//	if batchSize < 1 {
//		return fmt.Errorf("invalid batch size %d. Must be >= 1", batchSize)
//	}
//	concurrency := c.Flags.Int("concurrency")
//	if concurrency < 1 {
//		return fmt.Errorf("invalid concurrency %d. Must be >= 1", concurrency)
//	}
//	table := c.Args.String("table")
//	if table == "" {
//		return fmt.Errorf("please specify the table to run the workload on")
//	}
//
//	workloadName := c.Args.String("workload")
//	switch workloadName {
//	case "insert", "insert-txn",
//	     "select-pk", "select-pk-txn", "select-uk", "select-uk-txn",
//	     "update", "update-txn",
//	     "delete", "delete-txn":
//	default:
//		return fmt.Errorf("invalid workload '%s'", workloadName)
//	}
//
//	c.App.Println("Running workload:", workloadName)
//	c.App.Println("    dataspec:", controllerApp.data.dataSpecName)
//	c.App.Println("    numRows:", controllerApp.data.dataSpec.KeyGenSpec.NumKeys())
//	c.App.Println("    recordSize:", controllerApp.data.dataSpec.RecordSize)
//	c.App.Println("    driver:", controllerApp.db.driver)
//	c.App.Println("    url:", controllerApp.db.url)
//	c.App.Println("    table:", table)
//	c.App.Println("    batchSize:", batchSize)
//	c.App.Println("    concurrency:", concurrency)
//	c.App.Println("    duration(s):", durationSec)
//
//
//	// Partition the data among the number of minions
//	// TODO: We can improve this by checking which minions are live, has DB connection and dataspec loaded
//	// 	and then run the workload only on those minions
//	numRecs := controllerApp.data.dataSpec.KeyGenSpec.NumKeys()
//	fullRange := intgen.NewRange(0, numRecs)
//	numServers := len(controllerApp.minions.minions)
//	var assignedRanges []*intgen.Range
//	switch workloadName {
//	case "insert", "insert-txn", "delete", "delete-txn":
//		assignedRanges = fullRange.Partition(numServers)
//	case "select-pk", "select-pk-txn", "select-uk", "select-uk-txn", "update", "update-txn":
//		assignedRanges = fullRange.Duplicate(numServers)
//	default:
//		return fmt.Errorf("invalid workload '%s'", workloadName)
//	}
//
//	// Run the workload on all minions
//	c.App.Println("Running workload on all minions...")
//	success := 0
//	server := 0
//	for minAddr, minClient := range controllerApp.minions.minions {
//		ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
//		request := &proto.RunWorkloadRequest{
//			WorkloadName: workloadName,
//			AssignedRange: proto.RangeToProto(assignedRanges[server]),
//			TableName: table,
//			DurationSec: int32(durationSec),
//			Concurrency: int32(concurrency),
//			BatchSize: int32(batchSize),
//		}
//		server++
//		response, err := minClient.RunWorkload(ctx, request)
//		if err != nil {
//			c.App.Printf("%s: Request failed (%v)\n", minAddr, err)
//			continue
//		}
//		if !response.GetIsOk() {
//			c.App.Printf("%s: Failed! (%s)\n", minAddr, response.GetFailureReason())
//			continue
//		}
//		success++
//		c.App.Printf("%s: OK\n", minAddr)
//	}
//	c.App.Printf("Successful on %d/%d minions\n", success, numServers)
//	return nil
//}
//
//func wlShow(c *grumble.Context) error {
//	success := 0
//	for minAddr, minClient := range controllerApp.minions.minions {
//		ctx, _ := context.WithTimeout(context.Background(), time.Second)
//		request := &proto.GetWorkloadInfoRequest{}
//		response, err := minClient.GetWorkloadInfo(ctx, request)
//		if err != nil {
//			c.App.Printf("%s: Request failed (%v)\n", minAddr, err)
//			continue
//		}
//		if !response.GetIsOk() {
//			c.App.Printf("%s: Failed (%s)\n", minAddr, response.GetFailureReason())
//			continue
//		}
//		success++
//		running := response.GetWorkloadInfo().GetIsRunning()
//		workloadName := response.GetWorkloadInfo().GetWorkloadName()
//		tableName := response.GetWorkloadInfo().GetTableName()
//		if running {
//			c.App.Printf("%s: Currently running workload=%s table=%s\n", minAddr, workloadName, tableName)
//		} else {
//			c.App.Printf("%s: Finished running workload=%s table=%s\n", minAddr, workloadName, tableName)
//		}
//	}
//	c.App.Printf("Successful on %d/%d minions\n", success, len(controllerApp.minions.minions))
//	return nil
//}
//
//func wlStop(c *grumble.Context) error {
//	success := 0
//	for minAddr, minClient := range controllerApp.minions.minions {
//		ctx, _ := context.WithTimeout(context.Background(), time.Second)
//		request := &proto.StopWorkloadRequest{}
//		_, err := minClient.StopWorkload(ctx, request)
//		if err != nil {
//			c.App.Printf("%s: Request failed (%v)\n", minAddr, err)
//			continue
//		} else {
//			success++
//		}
//		c.App.Printf("%s: OK\n", minAddr)
//	}
//	c.App.Printf("Successful on %d/%d minions\n", success, len(controllerApp.minions.minions))
//	return nil
//}
