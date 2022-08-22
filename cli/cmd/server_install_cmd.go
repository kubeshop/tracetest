package cmd

import (
	"github.com/kubeshop/tracetest/cli/installer"
	"github.com/spf13/cobra"
)

var serverInstallCmd = &cobra.Command{
	Use:    "install",
	Short:  "install a new server",
	Long:   "install a new server",
	PreRun: setupCommand,
	Run: func(cmd *cobra.Command, args []string) {
		installer.Start()
	},
	PostRun: teardownCommand,
}

func init() {
	serverCmd.AddCommand(serverInstallCmd)
}
