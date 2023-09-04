package cmd

import (
	"context"

	"github.com/kubeshop/tracetest/cli/pkg/starter"
	"github.com/spf13/cobra"
)

var (
	start = starter.NewStarter(configurator, resources)
)

var startCmd = &cobra.Command{
	GroupID: cmdGroupCloud.ID,
	Use:     "start",
	Short:   "Start Tracetest",
	Long:    "Start using Tracetest",
	PreRun:  setupCommand(),
	Run: WithResultHandler((func(_ *cobra.Command, _ []string) (string, error) {
		ctx := context.Background()

		err := start.Run(ctx, cliConfig)
		return "", err
	})),
	PostRun: teardownCommand,
}

func init() {
	if isCloudEnabled {
		rootCmd.AddCommand(startCmd)
	}
}
