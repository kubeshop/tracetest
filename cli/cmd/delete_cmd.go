package cmd

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:    "delete [resource type]",
	Long:   "Delete resources from your Tracetest server",
	Short:  "Delete resources",
	PreRun: setupCommand(),
	Args:   cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
