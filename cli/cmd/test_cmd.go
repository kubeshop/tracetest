package cmd

import (
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	GroupID: cmdGroupTests.ID,
	Use:     "test",
	Short:   "Manage your tracetest tests",
	Long:    "Manage your tracetest tests",
	PreRun:  setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(testCmd)
}
