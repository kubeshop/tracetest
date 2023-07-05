package cmd

import (
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	GroupID:    cmdGroupTests.ID,
	Use:        "test",
	Short:      "Manage your tracetest tests",
	Long:       "Manage your tracetest tests",
	Deprecated: "Please use `tracetest (apply|delete|export|get|list) test` commands instead.",
	PreRun:     setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PostRun: teardownCommand,
}

var testListCmd = &cobra.Command{
	Use:        "list",
	Short:      "List all tests",
	Long:       "List all tests",
	Deprecated: "Please use `tracetest list test` command instead.",
	PreRun:     setupCommand(),
	Run: func(_ *cobra.Command, _ []string) {
		listCmd.Run(listCmd, []string{"test"})
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.AddCommand(testListCmd)
}
