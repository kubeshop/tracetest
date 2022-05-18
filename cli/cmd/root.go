package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var verbose bool
var configFile string

var rootCmd = &cobra.Command{
	Use:     "tracetest",
	Short:   "tracetest CLI is a tool to interact with a tracetest server",
	Long:    `tracetest CLI is a tool to interact with a tracetest server`,
	PreRun:  setupCommand,
	PostRun: teardownCommand,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file will be used by the CLI")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "display debug information")
}
