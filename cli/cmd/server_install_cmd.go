package cmd

import (
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/installer"
	"github.com/spf13/cobra"
)

var (
	force = false
)

var serverInstallCmd = &cobra.Command{
	Use:    "install",
	Short:  "Install a new Tracetest server",
	Long:   "Install a new Tracetest server",
	PreRun: setupCommand(SkipConfigValidation()),
	Run: func(cmd *cobra.Command, args []string) {
		installer.Force = force
		analytics.Track("Server Install", "cmd", map[string]string{})
		installer.Start()
	},
	PostRun: teardownCommand,
}

func init() {
	serverInstallCmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing files")
	serverCmd.AddCommand(serverInstallCmd)
}
