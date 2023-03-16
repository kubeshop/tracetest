package cmd

import (
	"fmt"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	GroupID: cmdGroupMisc.ID,
	Use:     "version",
	Short:   "Display this CLI tool version",
	Long:    "Display this CLI tool version",
	PreRun:  setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Track("Version", "cmd", map[string]string{})

		fmt.Println(config.Version)
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
