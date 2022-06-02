package main
//
//import (
//	"context"
//	"fmt"
//	"github.com/desertbit/grumble"
//	"github.com/flipkart-incubator/diligent/pkg/proto"
//	"google.golang.org/grpc"
//	"math/rand"
//	"time"
//)
//
//func init() {
//	minionCmd := &grumble.Command{
//		Name:    "minion",
//		Help:    "manage diligent minions",
//		Aliases: []string{"min"},
//	}
//	grumbleShell.AddCommand(minionCmd)
//
//	minionAddCmd := &grumble.Command{
//		Name: "add",
//		Help: "register a diligent minion with this controller",
//		Args: func(a *grumble.Args) {
//			a.String("address", "host:port of minion to add", grumble.Default(""))
//		},
//		Run: minionAdd,
//	}
//	minionCmd.AddCommand(minionAddCmd)
//
//	minionShowCmd := &grumble.Command{
//		Name:    "show",
//		Help:    "show and ping-check all registered diligent minions",
//		Aliases: []string{"sh"},
//		Run:     minionShow,
//	}
//	minionCmd.AddCommand(minionShowCmd)
//
//	minionPingCmd := &grumble.Command{
//		Name: "list",
//		Help: "list all registered diligent minions without any ping check",
//		Aliases: []string{"ls"},
//		Run:  minionList,
//	}
//	minionCmd.AddCommand(minionPingCmd)
//
//	minionRemoveCmd := &grumble.Command{
//		Name: "rm",
//		Help: "remove a registered diligent minion from this controller",
//		Args: func(a *grumble.Args) {
//			a.String("address", "host:port of minion to remove", grumble.Default(""))
//		},
//		Run: minionRemove,
//	}
//	minionCmd.AddCommand(minionRemoveCmd)
//}
//
//func minionAdd(c *grumble.Context) error {
//	address := c.Args.String("address")
//	if address == "" {
//		return fmt.Errorf("please provide a valid address for the minion")
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
//	if err != nil {
//		return fmt.Errorf("failed to connect to %s (%v)", address, err)
//	}
//
//	// Create client instance
//	client := proto.NewMinionClient(conn)
//	controllerApp.minions.minions[address] = client
//	return nil
//}
//
//func minionList(c *grumble.Context) error {
//	c.App.Printf("%d Server(s) registered:\n", len(controllerApp.minions.minions))
//	for address, _ := range controllerApp.minions.minions {
//		c.App.Println(address)
//	}
//	return nil
//}
//
//func minionShow(c *grumble.Context) error {
//	nonce := int32(rand.Intn(1024))
//	success := 0
//	for minAddr, minClient := range controllerApp.minions.minions {
//		ctx, _ := context.WithTimeout(context.Background(), time.Second)
//		response, err := minClient.Ping(ctx, &proto.PingRequest{Nonce: nonce})
//		if err != nil {
//			c.App.Printf("%s: Ping failed! (%v)\n", minAddr, err)
//			continue
//		}
//		if response.GetNonce() != nonce {
//			c.App.Printf("%s: Ping failed! (Nonce mismatch. Expected %d, got %d)\n", minAddr, nonce, response.GetNonce())
//			continue
//		}
//		c.App.Printf("%s: OK\n", minAddr)
//		success++
//	}
//	c.App.Printf("Successful on %d/%d minions\n", success, len(controllerApp.minions.minions))
//	return nil
//}
//
//func minionRemove(c *grumble.Context) error {
//	address := c.Args.String("address")
//	if address == "" {
//		return fmt.Errorf("please provide a valid address for the minion")
//	}
//
//	_, ok := controllerApp.minions.minions[address]
//	if !ok {
//		c.App.Printf("Server %s not registered\n", address)
//	} else {
//		// TODO: Need handle to conn object and close it
//		delete(controllerApp.minions.minions, address)
//		c.App.Printf("Server %s has been removed\n", address)
//	}
//	return nil
//}
