package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/cli/cmdutil"
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
	Run: WithResultHandler(WithParamsHandler(configParams)(func(ctx context.Context, cmd *cobra.Command, _ []string) (string, error) {
		flags := agentConfig.Flags{
			CI:         configParams.CI,
			SkipVerify: skipVerify,
		}

		cfg, err := config.LoadConfig("")
		if err != nil {
			return "", err
		}

		if flagProvided(cmd, "server-url") || flagProvided(cmd, "endpoint") {
			flags.ServerURL = configParams.ServerURL
		}

		if configParams.Token != "" {
			flags.Token = configParams.Token
		}

		if flagProvided(cmd, "environment") {
			flags.EnvironmentID = configParams.EnvironmentID
		}

		if flagProvided(cmd, "organization") {
			flags.OrganizationID = configParams.OrganizationID
		}

		configurator = configurator.WithLogger(cliLogger)

		// early exit if the versions are not compatible
		err = configurator.Start(ctx, &cfg, flags)
		if errors.Is(err, config.ErrVersionMismatch) {
			fmt.Println(err.Error())
			ExitCLI(1)
		}

		return "", err
	})),
	PostRun: teardownCommand,
}

func flagProvided(cmd *cobra.Command, name string) bool {
	return cmd.Flags().Lookup(name).Changed
}

var deprecatedEndpoint string

func init() {
	configureCmd.PersistentFlags().BoolVarP(&configParams.Global, "global", "g", false, "configuration will be saved in your home dir")

	configureCmd.PersistentFlags().StringVarP(&configParams.Token, "token", "t", defaultToken, "set authetication with token, so the CLI won't ask you for authentication")
	configureCmd.PersistentFlags().StringVarP(&configParams.EnvironmentID, "environment", "", "", "set environmentID, so the CLI won't ask you for it")
	configureCmd.PersistentFlags().StringVarP(&configParams.OrganizationID, "organization", "", "", "set organizationID, so the CLI won't ask you for it")

	configureCmd.PersistentFlags().BoolVarP(&configParams.CI, "ci", "", false, "if cloud is used, don't ask for authentication")

	configureCmd.PersistentFlags().StringVarP(&deprecatedEndpoint, "endpoint", "e", defaultEndpoint, "set the value for the endpoint, so the CLI won't ask for this value")
	configureCmd.PersistentFlags().MarkDeprecated("endpoint", "use --server-url instead")
	configureCmd.PersistentFlags().MarkShorthandDeprecated("e", "use --server-url instead")

	rootCmd.AddCommand(configureCmd)
}

type configureParameters struct {
	ServerURL      string
	Global         bool
	CI             bool
	Token          string
	OrganizationID string
	EnvironmentID  string
}

func (p *configureParameters) Validate(cmd *cobra.Command, args []string) []error {
	var errors []error

	p.ServerURL = overrideEndpoint
	flagUpdated := cmd.Flags().Lookup("server-url").Changed

	if deprecatedEndpoint != "" {
		p.ServerURL = deprecatedEndpoint
		flagUpdated = true
	}

	if flagUpdated {
		if p.ServerURL == "" {
			errors = append(errors, cmdutil.ParamError{
				Parameter: "server-url",
				Message:   "server-url cannot be empty",
			})
		} else {
			_, err := url.Parse(p.ServerURL)
			if err != nil {
				errors = append(errors, cmdutil.ParamError{
					Parameter: "server-url",
					Message:   "server-url is not a valid URL",
				})
			}
		}
	}

	return errors
}
