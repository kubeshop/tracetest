package cmd

import (
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	GroupID: cmdGroupMisc.ID,
	Use:     "version",
	Short:   "Display this CLI tool version",
	Long:    "Display this CLI tool version",
	PreRun:  setupCommand(),
	Run: WithResultHandler(func(cmd *cobra.Command, args []string) (string, error) {
		analytics.Track("Version", "cmd", map[string]string{})

		return versionText, nil
	}),
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
