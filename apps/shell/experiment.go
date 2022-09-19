package main

import (
	"context"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	"time"
)

func init() {
	expCmd := &grumble.Command{
		Name:    "experiment",
		Help:    "work with experiments",
		Aliases: []string{"ex"},
	}
	grumbleApp.AddCommand(expCmd)

	expBeginCmd := &grumble.Command{
		Name: "begin",
		Help: "begin an experiment",
		Args: func(a *grumble.Args) {
			a.String("name", "name of the experiment")
		},
		Run: expBegin,
	}
	expCmd.AddCommand(expBeginCmd)

	expEndCmd := &grumble.Command{
		Name: "end",
		Help: "end the current experiment",
		Run:  expEnd,
	}
	expCmd.AddCommand(expEndCmd)

	expInfoCmd := &grumble.Command{
		Name: "info",
		Help: "Show current experiment information",
		Run:  expInfo,
	}
	expCmd.AddCommand(expInfoCmd)

	expRunScript := &grumble.Command{
		Name: "run-script",
		Help: "run a scripted experiment",
		Flags: func(f *grumble.Flags) {
			f.String("v", "values-file", "", "YAML file for values")
			f.Bool("d", "dry-run", false, "Perform a dry run")
		},
		Args: func(a *grumble.Args) {
			a.String("script-file", "Script File")
		},
		Run: expRunScript,
	}
	expCmd.AddCommand(expRunScript)
}

func expBegin(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")

	// Experiment name param
	experimentName := c.Args.String("name")
	if experimentName == "" {
		return fmt.Errorf("please specify a name for the experiment")
	}

	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	res, err := bossClient.BeginExperiment(grpcCtx, &proto.BossBeginExperimentRequest{ExperimentName: experimentName})
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
	return nil
}

func expEnd(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	res, err := bossClient.EndExperiment(grpcCtx, &proto.BossEndExperimentRequest{})
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
	return nil
}

func expInfo(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	res, err := bossClient.GetExperimentInfo(grpcCtx, &proto.BossGetExperimentInfoRequest{})
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

	if res.GetExperimentInfo() == nil {
		c.App.Printf("No current experiment\n")
	} else {
		c.App.Printf("name      : %s\n", res.GetExperimentInfo().GetName())
		c.App.Printf("state     : %s\n", res.GetExperimentInfo().GetState())
		c.App.Printf("start-time: %s\n", time.UnixMilli(res.GetExperimentInfo().GetStartTime()).Format(time.UnixDate))
		c.App.Printf("stop-time : %s\n", time.UnixMilli(res.GetExperimentInfo().GetStopTime()).Format(time.UnixDate))
	}
	return nil
}

func expRunScript(c *grumble.Context) error {
	scriptFileName := c.Args.String("script-file")
	valuesFileName := c.Flags.String("values-file")
	dryRun := c.Flags.Bool("dry-run")

	script, err := LoadScript(scriptFileName)
	if err != nil {
		return err
	}

	values, err := LoadValues(valuesFileName)
	if err != nil {
		return err
	}

	executor, err := NewExecutor(script, values, dryRun)
	if err != nil {
		return err
	}

	err = executor.Execute()
	if err != nil {
		return err
	}
	return nil
}
