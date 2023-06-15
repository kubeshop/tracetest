package cmd

import (
	"context"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/parameters"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

var configParams = &parameters.ConfigureParams{}

var configureCmd = &cobra.Command{
	GroupID: cmdGroupConfig.ID,
	Use:     "configure",
	Short:   "Configure your tracetest CLI",
	Long:    "Configure your tracetest CLI",
	PreRun:  setupLogger,
	Run: WithResultHandler(WithParamsHandler(configParams)(func(cmd *cobra.Command, _ []string) (string, error) {
		ctx := context.Background()
		client := utils.GetAPIClient(cliConfig)
		action := actions.NewConfigureAction(cliConfig, cliLogger, client)

		actionConfig := actions.ConfigureConfig{
			Global:    configParams.Global,
			SetValues: actions.ConfigureConfigSetValues{},
		}

		if flagProvided(cmd, "endpoint") {
			actionConfig.SetValues.Endpoint = &configParams.Endpoint
		}

		err := action.Run(ctx, actionConfig)
		return "", err
	})),
	PostRun: teardownCommand,
}

func flagProvided(cmd *cobra.Command, name string) bool {
	return cmd.Flags().Lookup(name).Changed
}

func init() {
	configureCmd.PersistentFlags().BoolVarP(&configParams.Global, "global", "g", false, "configuration will be saved in your home dir")
	configureCmd.PersistentFlags().StringVarP(&configParams.Endpoint, "endpoint", "e", "", "set the value for the endpoint, so the CLI won't ask for this value")
	rootCmd.AddCommand(configureCmd)
}
