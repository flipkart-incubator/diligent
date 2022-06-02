package main

import (
	"context"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	"google.golang.org/grpc"
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
		Name: "minion-register",
		Help: "register a minion with the boss",
		Flags: func(f *grumble.Flags) {
			f.String("b", "boss-url", "", "URL of boss server")
			f.String("m", "minion-url", "", "URL of minion server")
		},
		Run: bossMinionRegister,
	}
	bossCmd.AddCommand(bossMinionRegisterCmd)

	bossMinionUnregisterCmd := &grumble.Command{
		Name: "minion-unregister",
		Help: "unregister a minion with the boss",
		Flags: func(f *grumble.Flags) {
			f.String("b", "boss-url", "", "URL of boss server")
			f.String("m", "minion-url", "", "URL of minion server")
		},
		Run: bossMinionUnregister,
	}
	bossCmd.AddCommand(bossMinionUnregisterCmd)

	bossMinionShowCmd := &grumble.Command{
		Name: "minion-show",
		Help: "show minions registered with the boss",
		Flags: func(f *grumble.Flags) {
			f.String("b", "boss-url", "", "URL of boss server")
		},
		Run: bossMinionShow,
	}
	bossCmd.AddCommand(bossMinionShowCmd)
}

func bossPing(c *grumble.Context) error {
	bossUrl := c.Flags.String("boss-url")
	if bossUrl == "" {
		return fmt.Errorf("please provide a valid bossUrl for the boss server")
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, bossUrl, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("failed to connect to %s (%v)", bossUrl, err)
	}
	bossClient := proto.NewBossClient(conn)

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	_, err = bossClient.Ping(ctx, &proto.BossPingRequest{})
	if err != nil {
		c.App.Printf("%s: Ping failed! (%v)\n", bossUrl, err)
	} else {
		c.App.Printf("OK\n")
	}
	return nil
}


func bossMinionRegister(c *grumble.Context) error {
	bossUrl := c.Flags.String("boss-url")
	if bossUrl == "" {
		return fmt.Errorf("please provide a valid bossUrl for the boss server")
	}

	minionUrl := c.Flags.String("minion-url")
	if minionUrl == "" {
		return fmt.Errorf("please provide a valid minionUrl for the minion server")
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, bossUrl, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("failed to connect to %s (%v)", bossUrl, err)
	}
	bossClient := proto.NewBossClient(conn)

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	_, err = bossClient.MinionRegister(ctx, &proto.MinionRegisterRequest{Url: minionUrl})
	if err != nil {
		c.App.Printf("Registration failed: (%v)\n", err)
	} else {
		c.App.Printf("OK\n")
	}
	return nil
}


func bossMinionUnregister(c *grumble.Context) error {
	bossUrl := c.Flags.String("boss-url")
	if bossUrl == "" {
		return fmt.Errorf("please provide a valid bossUrl for the boss server")
	}

	minionUrl := c.Flags.String("minion-url")
	if minionUrl == "" {
		return fmt.Errorf("please provide a valid minionUrl for the minion server")
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, bossUrl, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("failed to connect to %s (%v)", bossUrl, err)
	}
	bossClient := proto.NewBossClient(conn)

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	_, err = bossClient.MinionUnregister(ctx, &proto.MinionUnregisterRequest{Url: minionUrl})
	if err != nil {
		c.App.Printf("Unregistration failed: (%v)\n", err)
	} else {
		c.App.Printf("OK\n")
	}
	return nil
}


func bossMinionShow(c *grumble.Context) error {
	bossUrl := c.Flags.String("boss-url")
	if bossUrl == "" {
		return fmt.Errorf("please provide a valid bossUrl for the boss server")
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, bossUrl, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("failed to connect to %s (%v)", bossUrl, err)
	}
	bossClient := proto.NewBossClient(conn)

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	res, err := bossClient.MinionShow(ctx, &proto.MinionShowRequest{})
	if err != nil {
		c.App.Printf("Request failed: (%v)\n", err)
		return err
	}

	for _, minionStatus := range res.Minions {
		c.App.Printf("%s : %v\n", minionStatus.GetUrl(), minionStatus.GetIsReachable())
	}
	return nil
}
