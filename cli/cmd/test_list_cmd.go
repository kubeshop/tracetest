package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var testListCmd = &cobra.Command{
	Use:    "list",
	Short:  "list all test",
	Long:   "list all test",
	PreRun: setupCommand,
	Run: func(cmd *cobra.Command, args []string) {
		cliLogger.Debug("Retrieving list of tests", zap.String("endpoint", cliConfig.Endpoint))
	},
	PostRun: teardownCommand,
}

func init() {
	testCmd.AddCommand(testListCmd)
}
