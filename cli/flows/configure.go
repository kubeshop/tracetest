package flows

import (
	"context"

	"github.com/kubeshop/tracetest/cli/analytics"
	flows_actions "github.com/kubeshop/tracetest/cli/flows/actions"
	flows_parameters "github.com/kubeshop/tracetest/cli/flows/parameters"
	"github.com/kubeshop/tracetest/cli/global"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

type Configure struct {
	args[*flows_parameters.Configure]
}

func NewConfigure(root global.Root) Configure {
	parameters := flows_parameters.NewConfigure()
	defaults := NewDefaults("Configure", root.Setup)

	configure := Configure{
		args: NewArgs(defaults, parameters),
	}

	configure.Cmd = &cobra.Command{
		GroupID: global.GroupConfig.ID,
		Use:     "configure",
		Short:   "Configure your tracetest CLI",
		Long:    "Configure your tracetest CLI",
		PreRun:  defaults.PreRun,
		Run: defaults.Run(func(cmd *cobra.Command, _ []string) (string, error) {
			analytics.Track("Configure", "cmd", map[string]string{})

			ctx := context.Background()
			client := utils.GetAPIClient(*configure.Setup.Config)
			action := flows_actions.NewConfigureAction(*configure.Setup.Config, *&configure.Setup.Logger, client)

			actionConfig := flows_actions.ConfigureConfig{
				Global:    configure.Parameters.Global,
				SetValues: flows_actions.ConfigureConfigSetValues{},
			}

			if flagProvided(cmd, "endpoint") {
				actionConfig.SetValues.Endpoint = &configure.Parameters.Endpoint
			}

			if flagProvided(cmd, "analytics") {
				actionConfig.SetValues.AnalyticsEnabled = &configure.Parameters.AnalyticsEnabled
			}

			err := action.Run(ctx, actionConfig)
			return "", err
		}),
		PostRun: defaults.PostRun,
	}

	configure.Cmd.PersistentFlags().BoolVarP(&configure.Parameters.Global, "global", "g", false, "configuration will be saved in your home dir")
	configure.Cmd.PersistentFlags().StringVarP(&configure.Parameters.Endpoint, "endpoint", "e", "", "set the value for the endpoint, so the CLI won't ask for this value")
	configure.Cmd.PersistentFlags().BoolVarP(&configure.Parameters.AnalyticsEnabled, "analytics", "a", true, "configure the analytic state, so the CLI won't ask for this value")
	root.Cmd.AddCommand(configure.Cmd)

	return configure
}

func flagProvided(cmd *cobra.Command, name string) bool {
	return cmd.Flags().Lookup(name).Changed
}

var configure Configure

func GetConfigure() Configure {
	if configure.Cmd == nil {
		panic("Configure not initialized")
	}

	return configure
}
