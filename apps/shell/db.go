package main

//
//import (
//	"context"
//	"database/sql"
//	"fmt"
//	"github.com/desertbit/grumble"
//	"github.com/flipkart-incubator/diligent/pkg/proto"
//	"github.com/flipkart-incubator/diligent/pkg/work"
//
//	_ "github.com/go-sql-driver/mysql"
//	_ "github.com/jackc/pgx/stdlib"
//)
//
//func init() {
//	dbCmd := &grumble.Command{
//		Name:    "db",
//		Help:    "manage database connections",
//	}
//	grumbleApp.AddCommand(dbCmd)
//
//	dbConnCmd := &grumble.Command{
//		Name: "conn",
//		Help: "connect to database",
//		Flags: func(f *grumble.Flags) {
//			f.String("d", "driver", "mysql", "Driver to use: [mysql,pgx]")
//		},
//		Args: func(a *grumble.Args) {
//			a.String("url", "Connection URL including DB (example: 'user:password@tcp(host:port)/db')", grumble.Default(""))
//		},
//		Run: dbConn,
//	}
//	dbCmd.AddCommand(dbConnCmd)
//
//	dbShowCmd := &grumble.Command{
//		Name: "show",
//		Help: "show database connection",
//		Aliases: []string{"sh"},
//		Run:  dbShow,
//	}
//	dbCmd.AddCommand(dbShowCmd)
//
//	dbCloseCmd := &grumble.Command{
//		Name: "close",
//		Help: "close database connection",
//		Run: dbClose,
//	}
//	dbCmd.AddCommand(dbCloseCmd)
//}
//
//func dbConn(c *grumble.Context) error {
//	// Validate driver
//	driver := c.Flags.String("driver")
//	switch driver {
//	case "mysql", "pgx":
//		controllerApp.db.driver = driver
//	default:
//		return fmt.Errorf("invalid driver: '%s'. Allowed values are 'mysql', 'pgx'", driver)
//	}
//
//	// Validate URL
//	url := c.Args.String("url")
//	if url == "" {
//		return fmt.Errorf("please specify the connection url")
//	}
//	controllerApp.db.url = url
//
//	// Close any existing connection
//	if controllerApp.db.db != nil {
//		c.App.Println("Closing existing local connection")
//		controllerApp.db.db.Close()
//		controllerApp.db.db = nil
//		controllerApp.db.driver = ""
//		controllerApp.db.url = ""
//	}
//
//	// Connect locally and check
//	c.App.Println("Connecting to database locally...")
//	db, err := sql.Open(controllerApp.db.driver, controllerApp.db.url)
//	if err != nil {
//		return err
//	}
//	controllerApp.db.db = db
//	err = work.ConnCheck(controllerApp.db.db)
//	if err != nil {
//		return err
//	}
//	c.App.Println("Local connection successful!")
//
//	// Ask minions to connect
//	c.App.Println("Asking minions to connect...")
//	success := 0
//	for minAddr, minClient := range controllerApp.minions.minions {
//		request := &proto.OpenDBConnectionRequest{Driver: driver, Url: url}
//		response, err := minClient.OpenDBConnection(context.Background(), request)
//		if err != nil {
//			c.App.Printf("%s: Request failed (%v)\n", minAddr, err)
//			continue
//		}
//		if !response.GetIsOk() {
//			c.App.Printf("%s: Failed! (%s)\n", minAddr, response.GetFailureReason())
//			continue
//		}
//		c.App.Printf("%s: OK\n", minAddr)
//		success++
//	}
//	c.App.Printf("Successful on %d/%d minions\n", success, len(controllerApp.minions.minions))
//
//	return nil
//}
//
//func dbShow(c *grumble.Context) error {
//	if controllerApp.db.driver == "" || controllerApp.db.url == "" || controllerApp.db.db == nil {
//		return fmt.Errorf("no connection defined. Please use 'db conn' command to connect")
//	}
//
//	c.App.Println("Checking connection locally...")
//	err := work.ConnCheck(controllerApp.db.db)
//	if err != nil {
//		return err
//	}
//	c.App.Println("Local connection successful!")
//
//	c.App.Println("Checking connection from minions...")
//	success := 0
//	for minAddr, minClient := range controllerApp.minions.minions {
//		request := &proto.GetDBConnectionInfoRequest{}
//		response, err := minClient.GetDBConnectionInfo(context.Background(), request)
//		if err != nil {
//			c.App.Printf("%s: Request failed (%v)\n", minAddr, err)
//			continue
//		}
//		if !response.GetIsOk() {
//			c.App.Printf("%s: Failed (%s)\n", minAddr, response.GetFailureReason())
//			continue
//		}
//		c.App.Printf("%s: OK [driver=%s, url=%s]\n", minAddr, response.GetDriver(), response.GetUrl())
//		success++
//	}
//	c.App.Printf("Successful on %d/%d minions\n", success, len(controllerApp.minions.minions))
//
//	return nil
//}
//
//func dbClose(c *grumble.Context) error {
//	if controllerApp.db.db != nil {
//		c.App.Println("Closing existing local connection")
//		controllerApp.db.db.Close()
//		controllerApp.db.db = nil
//		controllerApp.db.driver = ""
//		controllerApp.db.url = ""
//	}
//
//	c.App.Println("Asking minions to close connection...")
//	success := 0
//	for minAdddr, minClient := range controllerApp.minions.minions {
//		request := &proto.CloseDBConnectionRequest{}
//		_, err := minClient.CloseDBConnection(context.Background(), request)
//		if err != nil {
//			c.App.Printf("%s: Request failed (%v)\n", minAdddr, err)
//		} else {
//			success++
//		}
//	}
//	c.App.Printf("Successful on %d/%d minions\n", success, len(controllerApp.minions.minions))
//
//	return nil
//}
