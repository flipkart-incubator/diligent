package main

import (
	"context"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	"google.golang.org/grpc"
	"net"
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

	bossInfoCmd := &grumble.Command{
		Name: "info",
		Help: "show boss info",
		Run:  bossInfo,
	}
	bossCmd.AddCommand(bossInfoCmd)
}

func bossInfo(c *grumble.Context) error {
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
