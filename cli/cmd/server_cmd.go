package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:    "server",
	Short:  "Manage your tracetest server",
	Long:   "Manage your tracetest server",
	PreRun: setupCommand,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Manage your server")
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
