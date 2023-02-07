package cmd

import (
	"github.com/spf13/cobra"
)

var dataStoreCmd = &cobra.Command{
	Use:    "datastore",
	Short:  "Manage your tracetest data stores",
	Long:   "Manage your tracetest data stores",
	PreRun: setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(dataStoreCmd)
}
