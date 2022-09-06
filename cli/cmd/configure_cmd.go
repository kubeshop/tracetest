package cmd

import (
	"context"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var global bool

var configureCmd = &cobra.Command{
	Use:    "configure",
	Short:  "Configure your tracetest CLI",
	Long:   "Configure your tracetest CLI",
	PreRun: setupLogger,
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Track("Configure", "cmd", map[string]string{})

		ctx := context.Background()
		client := getAPIClient()
		action := actions.NewConfigureAction(cliConfig, cliLogger, client)

		actionConfig := actions.ConfigureConfig{
			Global: global,
		}
		err := action.Run(ctx, actionConfig)
		if err != nil {
			cliLogger.Error("could not get tests", zap.Error(err))
			return
		}
	},
	PostRun: teardownCommand,
}

func init() {
	configureCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "configuration will be saved in your home dir")
	rootCmd.AddCommand(configureCmd)
}
