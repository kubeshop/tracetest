package cmd

import (
	"context"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var analyticsEnabled bool
var endpoint string
var global bool

var configureCmd = &cobra.Command{
	GroupID: cmdGroupConfig.ID,
	Use:     "configure",
	Short:   "Configure your tracetest CLI",
	Long:    "Configure your tracetest CLI",
	PreRun:  setupLogger,
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Track("Configure", "cmd", map[string]string{})

		ctx := context.Background()
		client := utils.GetAPIClient(cliConfig)
		action := actions.NewConfigureAction(cliConfig, cliLogger, client)

		actionConfig := actions.ConfigureConfig{
			Global:    global,
			SetValues: actions.ConfigureConfigSetValues{},
		}

		if flagProvided(cmd, "endpoint") {
			actionConfig.SetValues.Endpoint = &endpoint
		}

		if flagProvided(cmd, "analytics") {
			actionConfig.SetValues.AnalyticsEnabled = &analyticsEnabled
		}

		err := action.Run(ctx, actionConfig)
		if err != nil {
			cliLogger.Error("could not get tests", zap.Error(err))
			return
		}
	},
	PostRun: teardownCommand,
}

func flagProvided(cmd *cobra.Command, name string) bool {
	return cmd.Flags().Lookup(name).Changed
}

func init() {
	configureCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "configuration will be saved in your home dir")
	configureCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "", "set the value for the endpoint, so the CLI won't ask for this value")
	configureCmd.PersistentFlags().BoolVarP(&analyticsEnabled, "analytics", "a", true, "configure the analytic state, so the CLI won't ask for this value")
	rootCmd.AddCommand(configureCmd)
}
