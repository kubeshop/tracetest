package cmd

import (
	"fmt"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:    "test",
	Short:  "Manage your tracetest tests",
	Long:   "Manage your tracetest tests",
	PreRun: setupCommand,
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Track("Test", "cmd", map[string]string{})

		fmt.Println("Manage your tests")
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(testCmd)
}
