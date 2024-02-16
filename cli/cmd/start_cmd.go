package cmd

import (
	"context"
	"os"

	"github.com/davecgh/go-spew/spew"
	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/runner"
	"github.com/kubeshop/tracetest/agent/ui"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/spf13/cobra"
)

var (
	agentRunner     = runner.NewRunner(configurator.WithErrorHandler(handleError), resources, ui.DefaultUI)
	defaultToken    = os.Getenv("TRACETEST_TOKEN")
	defaultEndpoint = os.Getenv("TRACETEST_SERVER_URL")
	defaultAPIKey   = os.Getenv("TRACETEST_API_KEY")
	saveParams      = &saveParameters{}
)

var startCmd = &cobra.Command{
	GroupID: cmdGroupConfig.ID,
	Use:     "start",
	Short:   "Start Tracetest",
	Long:    "Start using Tracetest",
	PreRun:  setupCommand(SkipConfigValidation(), SkipVersionMismatchCheck()),
	Run: WithResultHandler((func(ctx context.Context, _ *cobra.Command, _ []string) (string, error) {
		flags := agentConfig.Flags{
			OrganizationID: saveParams.organizationID,
			EnvironmentID:  saveParams.environmentID,
			ServerURL:      saveParams.endpoint,
			AgentApiKey:    saveParams.agentApiKey,
			Token:          saveParams.token,
			Mode:           agentConfig.Mode(saveParams.mode),
			LogLevel:       saveParams.logLevel,
		}

		// override organization and environment id from context.
		// this happens when auto rerunning the cmd after relogin
		if orgID := config.ContextGetOrganizationID(ctx); orgID != "" {
			flags.OrganizationID = orgID
		}
		if envID := config.ContextGetEnvironmentID(ctx); envID != "" {
			flags.EnvironmentID = envID
		}

		spew.Dump(flags)

		cfg, err := agentConfig.LoadConfig()
		if err != nil {
			return "", err
		}

		if cfg.APIKey != "" {
			flags.AgentApiKey = cfg.APIKey
		}

		err = agentRunner.Run(ctx, cliConfig, flags)
		return "", err
	})),
	PostRun: teardownCommand,
}

func init() {
	startCmd.Flags().StringVarP(&saveParams.organizationID, "organization", "", "", "organization id")
	startCmd.Flags().StringVarP(&saveParams.environmentID, "environment", "", "", "environment id")
	startCmd.Flags().StringVarP(&saveParams.agentApiKey, "api-key", "", defaultAPIKey, "agent api key")
	startCmd.Flags().StringVarP(&saveParams.token, "token", "", defaultToken, "token authentication key")
	startCmd.Flags().StringVarP(&saveParams.endpoint, "endpoint", "e", defaultEndpoint, "set the value for the endpoint, so the CLI won't ask for this value")
	startCmd.Flags().StringVarP(&saveParams.mode, "mode", "m", "desktop", "set how the agent will start")
	startCmd.Flags().StringVarP(&saveParams.logLevel, "log-level", "l", "debug", "set the agent log level")

	startCmd.Flags().MarkDeprecated("endpoint", "use --server-url instead")
	startCmd.Flags().MarkShorthandDeprecated("e", "use --server-url instead")

	rootCmd.AddCommand(startCmd)
}

type saveParameters struct {
	organizationID string
	environmentID  string
	endpoint       string
	agentApiKey    string
	token          string
	mode           string
	logLevel       string
}
