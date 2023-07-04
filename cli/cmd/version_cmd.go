package cmd

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	GroupID: cmdGroupMisc.ID,
	Use:     "version",
	Short:   "Display this CLI tool version",
	Long:    "Display this CLI tool version",
	PreRun:  setupCommand(),
	Run: WithResultHandler(func(_ *cobra.Command, _ []string) (string, error) {
		return versionText, nil
	}),
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
