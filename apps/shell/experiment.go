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

	expStartCmd := &grumble.Command{
		Name: "start",
		Help: "start an experiment",
		Args: func(a *grumble.Args) {
			a.String("name", "name of the experiment")
		},
		Run: expStart,
	}
	expCmd.AddCommand(expStartCmd)

	expStopCmd := &grumble.Command{
		Name: "stop",
		Help: "stop the current experiment",
		Run:  expStop,
	}
	expCmd.AddCommand(expStopCmd)

	expInfoCmd := &grumble.Command{
		Name: "info",
		Help: "Show current experiment information",
		Run:  expInfo,
	}
	expCmd.AddCommand(expInfoCmd)
}

func expStart(c *grumble.Context) error {
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
	res, err := bossClient.StartExperiment(grpcCtx, &proto.BossStartExperimentRequest{ExperimentName: experimentName})
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

func expStop(c *grumble.Context) error {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return err
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	reqStart := time.Now()
	res, err := bossClient.StopExperiment(grpcCtx, &proto.BossStopExperimentRequest{})
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
