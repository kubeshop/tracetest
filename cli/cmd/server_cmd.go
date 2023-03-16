package cmd

import (
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	GroupID: cmdGroupConfig.ID,
	Use:     "server",
	Short:   "Manage your tracetest server",
	Long:    "Manage your tracetest server",
	PreRun:  setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
