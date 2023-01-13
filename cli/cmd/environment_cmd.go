package cmd

import (
	"fmt"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/spf13/cobra"
)

var environmentCmd = &cobra.Command{
	Use:    "environment",
	Short:  "Manage your tracetest environments",
	Long:   "Manage your tracetest environments",
	PreRun: setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Track("Environment", "cmd", map[string]string{})

		fmt.Println("Manage your environments")
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(environmentCmd)
}
