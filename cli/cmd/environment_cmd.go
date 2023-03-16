package cmd

import (
	"github.com/spf13/cobra"
)

var environmentCmd = &cobra.Command{
	GroupID: cmdGroupConfig.ID,
	Use:     "environment",
	Short:   "Manage your tracetest environments",
	Long:    "Manage your tracetest environments",
	PreRun:  setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(environmentCmd)
}
