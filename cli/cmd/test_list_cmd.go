package cmd

import (
	"github.com/spf13/cobra"
)

var testListCmd = &cobra.Command{
	Use:        "list",
	Short:      "List all tests",
	Long:       "List all tests",
	Deprecated: "Please use `tracetest list test` commands instead.",
	PreRun:     setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		listCmd.Run(listCmd, []string{"test"})
	},
	PostRun: teardownCommand,
}

func init() {
	testCmd.AddCommand(testListCmd)
}
