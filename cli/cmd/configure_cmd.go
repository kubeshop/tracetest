package cmd

import (
	"context"
	"net/url"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/spf13/cobra"
)

var configParams = &configureParameters{}

var configureCmd = &cobra.Command{
	GroupID: cmdGroupConfig.ID,
	Use:     "configure",
	Short:   "Configure your tracetest CLI",
	Long:    "Configure your tracetest CLI",
	PreRun:  setupLogger,
	Run: WithResultHandler(WithParamsHandler(configParams)(func(cmd *cobra.Command, _ []string) (string, error) {
		ctx := context.Background()
		action := actions.NewConfigureAction(cliConfig)

		actionConfig := actions.ConfigureConfig{
			Global: configParams.Global,
		}

		if flagProvided(cmd, "endpoint") {
			actionConfig.SetValues.Endpoint = configParams.Endpoint
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

type configureParameters struct {
	Endpoint string
	Global   bool
}

func (p configureParameters) Validate(cmd *cobra.Command, args []string) []error {
	var errors []error

	if cmd.Flags().Lookup("endpoint").Changed {
		if p.Endpoint == "" {
			errors = append(errors, ParamError{
				Parameter: "endpoint",
				Message:   "endpoint cannot be empty",
			})
		} else {
			_, err := url.Parse(p.Endpoint)
			if err != nil {
				errors = append(errors, ParamError{
					Parameter: "endpoint",
					Message:   "endpoint is not a valid URL",
				})
			}
		}
	}

	return errors
}
