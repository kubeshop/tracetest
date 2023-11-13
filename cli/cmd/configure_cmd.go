package cmd

import (
	"context"
	"net/url"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/spf13/cobra"
)

var configParams = &configureParameters{}

var (
	configurator = config.NewConfigurator(resources)
)

var configureCmd = &cobra.Command{
	GroupID: cmdGroupConfig.ID,
	Use:     "configure",
	Short:   "Configure your tracetest CLI",
	Long:    "Configure your tracetest CLI",
	PreRun:  setupLogger,
	Run: WithResultHandler(WithParamsHandler(configParams)(func(cmd *cobra.Command, _ []string) (string, error) {
		ctx := context.Background()
		flags := config.ConfigFlags{
			CI: configParams.CI,
		}
		config, err := config.LoadConfig("")
		if err != nil {
			return "", err
		}

		if flagProvided(cmd, "endpoint") {
			flags.Endpoint = configParams.Endpoint
		}

		if flagProvided(cmd, "token") {
			flags.Token = configParams.Token
		}

		if flagProvided(cmd, "environment") {
			flags.EnvironmentID = configParams.EnvironmentID
		}

		if flagProvided(cmd, "organization") {
			flags.OrganizationID = configParams.OrganizationID
		}

		err = configurator.Start(ctx, config, flags)
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

	configureCmd.PersistentFlags().StringVarP(&configParams.Token, "token", "t", "", "set authetication with token, so the CLI won't ask you for authentication")
	configureCmd.PersistentFlags().StringVarP(&configParams.EnvironmentID, "environment", "", "", "set environmentID, so the CLI won't ask you for it")
	configureCmd.PersistentFlags().StringVarP(&configParams.OrganizationID, "organization", "", "", "set organizationID, so the CLI won't ask you for it")

	configureCmd.PersistentFlags().BoolVarP(&configParams.CI, "ci", "", false, "if cloud is used, don't ask for authentication")
	rootCmd.AddCommand(configureCmd)
}

type configureParameters struct {
	Endpoint       string
	Global         bool
	CI             bool
	Token          string
	OrganizationID string
	EnvironmentID  string
}

func (p configureParameters) Validate(cmd *cobra.Command, args []string) []error {
	var errors []error

	if cmd.Flags().Lookup("endpoint").Changed {
		if p.Endpoint == "" {
			errors = append(errors, paramError{
				Parameter: "endpoint",
				Message:   "endpoint cannot be empty",
			})
		} else {
			_, err := url.Parse(p.Endpoint)
			if err != nil {
				errors = append(errors, paramError{
					Parameter: "endpoint",
					Message:   "endpoint is not a valid URL",
				})
			}
		}
	}

	return errors
}
