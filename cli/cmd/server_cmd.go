package cmd

import (
	"fmt"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:    "server",
	Short:  "Manage your tracetest server",
	Long:   "Manage your tracetest server",
	PreRun: setupCommand,
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Track("Server", "cmd", map[string]string{})
		fmt.Println("Manage your server")
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
