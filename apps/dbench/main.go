package main

import (
	"fmt"
	"os"
	shell "github.com/brianstrauch/cobra-shell"

)

func main() {
	rootCmd.AddCommand(shell.New(rootCmd, nil))
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
