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

	shellVersionCmd := &grumble.Command{
		Name: "version",
		Help: "get version information about the shell",
		Run:  shellVersion,
	}
	shellCmd.AddCommand(shellVersionCmd)
}

func shellVersion(c *grumble.Context) error {
	c.App.Printf("go-version : %s\n", buildinfo.GoVersion)
	c.App.Printf("commit-hash: %s\n", buildinfo.CommitHash)
	c.App.Printf("build-time : %s\n", buildinfo.BuildTime)
	return nil
}
