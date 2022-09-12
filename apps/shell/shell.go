package main

import (
	"github.com/desertbit/grumble"
	"github.com/flipkart-incubator/diligent/pkg/buildinfo"
)

func init() {
	shellCmd := &grumble.Command{
		Name:    "shell",
		Help:    "work with the shell",
		Aliases: []string{"sh"},
	}
	grumbleApp.AddCommand(shellCmd)

	shellInfoCmd := &grumble.Command{
		Name: "info",
		Help: "get information about the shell",
		Run:  shellInfo,
	}
	shellCmd.AddCommand(shellInfoCmd)
}

func shellInfo(c *grumble.Context) error {
	c.App.Printf("go-version : %s\n", buildinfo.GoVersion)
	c.App.Printf("commit-hash: %s\n", buildinfo.CommitHash)
	c.App.Printf("build-time : %s\n", buildinfo.BuildTime)
	return nil
}
