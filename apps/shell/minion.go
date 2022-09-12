package main

import (
	"context"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	"net"
	"time"
)

func init() {
	minionCmd := &grumble.Command{
		Name:    "minion",
		Help:    "work with the minions",
		Aliases: []string{"mi"},
	}
	grumbleApp.AddCommand(minionCmd)

	minionRegisterCmd := &grumble.Command{
		Name: "register",
		Help: "manually register a minion with the boss",
		Args: func(a *grumble.Args) {
			a.String("minion-addr", "host[:port] of minion server")
		},
		Run: minionRegister,
	}
	minionCmd.AddCommand(minionRegisterCmd)

	minionUnregisterCmd := &grumble.Command{
		Name: "unregister",
		Help: "manually unregister a minion with the boss",
		Args: func(a *grumble.Args) {
			a.String("minion-addr", "host[:port] of minion server")
		},
		Run: minionUnregister,
	}
	minionCmd.AddCommand(minionUnregisterCmd)

	minionInfoCmd := &grumble.Command{
		Name: "info",
		Help: "show info of minions registered with the boss",
		Flags: func(f *grumble.Flags) {
			f.Bool("b", "build-info", false, "show build information")
			f.Bool("p", "process-info", false, "show process information")
			f.Bool("j", "job-info", false, "show job information")
		},
		Run: minionInfo,
	}
	minionCmd.AddCommand(minionInfoCmd)

	minionAwaitCmd := &grumble.Command{
		Name: "await-count",
		Help: "wait until boss reports a desired number of live minions",
		Flags: func(f *grumble.Flags) {
			f.Duration("t", "timeout", 30*time.Minute, "wait timeout")
		},
		Args: func(a *grumble.Args) {
			a.Int("num-minions", "number of minions to wait for")
		},
		Run: minionAwaitCount,
	}
	minionCmd.AddCommand(minionAwaitCmd)
}

func minionRegister(c *grumble.Context) error {
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

func minionUnregister(c *grumble.Context) error {
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
	_, err = bossClient.UnregisterMinion(grpcCtx, &proto.BossUnregisterMinionRequest{Addr: minionAddr})
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

func minionInfo(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	res, err := bossClient.GetMinions(grpcCtx, &proto.BossGetMinionsRequest{})
	reqDuration := time.Since(reqStart)
	grpcCancel()
	if err != nil {
		c.App.Printf("Request failed [elapsed=%v]\n", reqDuration)
		return err
	} else {
		c.App.Printf("OK [elapsed=%v]\n", reqDuration)
	}

	for _, mi := range res.GetMinionInfos() {
		if !mi.GetReachability().GetIsOk() {
			c.App.Printf("%s : Not Reachable [%s]\n", mi.GetAddr(), mi.GetReachability().GetFailureReason())
			continue
		}
		c.App.Printf("%s : OK ", mi.GetAddr())
		switch {
		case c.Flags.Bool("build-info"):
			showMinionBuildInfo(c, mi.GetBuildInfo())
		case c.Flags.Bool("process-info"):
			showMinionProcessInfo(c, mi.GetProcessInfo())
		case c.Flags.Bool("job-info"):
			showMinionJobInfo(c, mi.GetJobInfo())
		default:
			showMinionSummaryInfo(c, mi)
		}
	}

	return nil
}

func showMinionBuildInfo(c *grumble.Context, bi *proto.BuildInfo) {
	if bi == nil {
		c.App.Printf("[No build info]\n")
		return
	}
	c.App.Printf("\n")
	c.App.Printf("\tapp-name: %s\n", bi.GetAppName())
	c.App.Printf("\tapp-version: %s\n", bi.GetAppVersion())
	c.App.Printf("\tcommit-hash: %s\n", bi.GetCommitHash())
	c.App.Printf("\tgo-version: %s\n", bi.GetGoVersion())
	c.App.Printf("\tbuild-time: %s\n", bi.GetBuildTime())
}

func showMinionProcessInfo(c *grumble.Context, pi *proto.ProcessInfo) {
	if pi == nil {
		c.App.Printf("[No process info]\n")
		return
	}

	startTime := time.UnixMilli(pi.GetStartTime())
	upTime := time.Since(startTime).Round(time.Second)

	c.App.Printf("\n")
	c.App.Printf("\tpid: %s\n", pi.GetPid())
	c.App.Printf("\tstart-time: %s (%s ago)\n", startTime.Format(time.UnixDate), upTime.String())
}

func showMinionJobInfo(c *grumble.Context, ji *proto.JobInfo) {
	if ji == nil {
		c.App.Printf("[No job info]\n")
		return
	}
	ds := proto.DataSpecFromProto(ji.GetJobSpec().GetDataSpec())
	c.App.Printf("\n")
	c.App.Printf("\tjob-name:  %s\n", ji.GetJobName())
	c.App.Printf("\tjob-state: %s\n", ji.GetJobState())
	c.App.Printf("\tprepare-time: %s\n", time.UnixMilli(ji.GetPrepareTime()).Format(time.UnixDate))
	c.App.Printf("\trun-time:     %s\n", time.UnixMilli(ji.GetRunTime()).Format(time.UnixDate))
	c.App.Printf("\tend-time:     %s\n", time.UnixMilli(ji.GetEndTime()).Format(time.UnixDate))
	c.App.Printf("\tfatal-errors:     %d\n", ji.GetFatalErrors())
	c.App.Printf("\tnon-fatal-errors: %d\n", ji.GetNonFatalErrors())
	c.App.Printf("\tdata-spec:\n")
	c.App.Printf("\t\tdata-num-recs: %d\n", ds.KeyGenSpec.NumKeys())
	c.App.Printf("\t\tdata-rec-size: %d\n", ji.GetJobSpec().GetDataSpec().GetRecordSize())
	c.App.Printf("\tdb-spec:\n")
	c.App.Printf("\t\tdb-driver: %s\n", ji.GetJobSpec().GetDbSpec().GetDriver())
	c.App.Printf("\t\tdb-url:    %s\n", ji.GetJobSpec().GetDbSpec().GetUrl())
	c.App.Printf("\tworkload-spec:\n")
	c.App.Printf("\t\tworkload-name:           %s\n", ji.GetJobSpec().GetWorkloadSpec().GetWorkloadName())
	c.App.Printf("\t\ttable-name:              %s\n", ji.GetJobSpec().GetWorkloadSpec().GetTableName())
	c.App.Printf("\t\tworkload-assigned-range: %s\n", ji.GetJobSpec().GetWorkloadSpec().GetAssignedRange())
	c.App.Printf("\t\tworkload-batch-size:     %d\n", ji.GetJobSpec().GetWorkloadSpec().GetBatchSize())
	c.App.Printf("\t\tworkload-concurrency:    %d\n", ji.GetJobSpec().GetWorkloadSpec().GetConcurrency())
	c.App.Printf("\t\tworkload-duration-sec:   %d\n", ji.GetJobSpec().GetWorkloadSpec().GetDurationSec())
}

func showMinionSummaryInfo(c *grumble.Context, mi *proto.MinionInfo) {
	startTime := time.UnixMilli(mi.GetProcessInfo().GetStartTime())
	upTime := time.Since(startTime).Round(time.Second)
	c.App.Printf("[version=%s, pid=%s, start-time=%s (%s ago), ",
		mi.GetBuildInfo().GetAppVersion(), mi.GetProcessInfo().GetPid(),
		startTime.Format(time.UnixDate), upTime.String())

	if mi.GetJobInfo() == nil {
		c.App.Printf("idle]\n")
	} else {
		c.App.Printf("job-name=%s, job-state=%s]\n",
			mi.GetJobInfo().GetJobName(), mi.GetJobInfo().GetJobState())
	}
}

func minionAwaitCount(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	desiredMinions := c.Args.Int("num-minions")
	timeout := c.Flags.Duration("timeout")

	c.App.Printf("Waiting till boss has %d live minions. Wait timeout=%s\n", desiredMinions, timeout.String())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
		res, err := bossClient.GetMinions(grpcCtx, &proto.BossGetMinionsRequest{})
		grpcCancel()
		if err != nil {
			c.App.Printf("Request to boss failed (%s)\n", err.Error())
		} else {
			count := 0
			for _, mi := range res.GetMinionInfos() {
				if mi.GetReachability().GetIsOk() {
					count++
				}
			}
			if count >= desiredMinions {
				c.App.Printf("Boss is reporting %d live minions\n", count)
				break
			}
		}
		if ctx.Err() != nil {
			return fmt.Errorf("timeout occurred. desired number of minions not found")
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
