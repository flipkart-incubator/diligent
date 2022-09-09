package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dbench",
	Short: "dbench is the diligent benchmark runner",
}

func init() {
	var err error
	runBenchmarkCmd := &cobra.Command{
		Use:       "run-benchmark name",
		Short:     "run a standard benchmark",
		Long:      "Runs one of the standard benchmarks defined by diligent.\nRequires environment and override parameters to be provided via a yaml file",
		ValidArgs: []string{"b00-simple-dml"},
		Args:      cobra.OnlyValidArgs,
		Run:       runBenchmark,
	}
	runBenchmarkCmd.Flags().StringP("values-file", "f", "", "yaml file containing values")
	err = runBenchmarkCmd.MarkFlagRequired("values-file")
	if err != nil {
		panic(err)
	}
	rootCmd.AddCommand(runBenchmarkCmd)

	runScriptCmd := &cobra.Command{
		Use:   "run-script script",
		Short: "run a custom benchmark script",
		Args:  cobra.ExactArgs(1),
		RunE:  runScript,
	}
	runScriptCmd.Flags().StringP("values-file", "f", "", "yaml file containing values")
	err = runScriptCmd.MarkFlagRequired("values-file")
	if err != nil {
		panic(err)
	}
	rootCmd.AddCommand(runScriptCmd)
}

func runBenchmark(cmd *cobra.Command, args []string) {
	fmt.Println("run-benchmark", args)
}

func runScript(cmd *cobra.Command, args []string) error {
	fmt.Println("run-script", args)

	scriptFileName := args[0]
	valuesFileName, err := cmd.Flags().GetString("values-file")
	if err != nil {
		panic(err)
	}

	bms, err := LoadScript(scriptFileName)
	if err != nil {
		return err
	}

	bmv, err := LoadValues(valuesFileName)
	if err != nil {
		return err
	}

	err = bms.Execute(bmv)
	if err != nil {
		return err
	}
	return nil
}

