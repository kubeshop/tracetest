package cmd

import (
	"context"
	"os"

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
	startParams     = &startParameters{}
)

var startCmd = &cobra.Command{
	GroupID: cmdGroupConfig.ID,
	Use:     "start",
	Short:   "Start Tracetest",
	Long:    "Start using Tracetest",
	PreRun:  setupCommand(SkipConfigValidation(), SkipVersionMismatchCheck()),
	Run: WithResultHandler((func(ctx context.Context, _ *cobra.Command, _ []string) (string, error) {
		flags := agentConfig.Flags{
			OrganizationID:    startParams.organizationID,
			EnvironmentID:     startParams.environmentID,
			ServerURL:         overrideEndpoint,
			AgentApiKey:       startParams.agentApiKey,
			Token:             startParams.token,
			Mode:              agentConfig.Mode(startParams.mode),
			LogLevel:          startParams.logLevel,
			CollectorEndpoint: startParams.collectorEndpoint,
		}

		// override organization and environment id from context.
		// this happens when auto rerunning the cmd after relogin
		if orgID := config.ContextGetOrganizationID(ctx); orgID != "" {
			flags.OrganizationID = orgID
		}
		if envID := config.ContextGetEnvironmentID(ctx); envID != "" {
			flags.EnvironmentID = envID
		}
		if serverURL := config.ContextGetServerURL(ctx); serverURL != "" {
			flags.ServerURL = serverURL
		}

		cfg, err := agentConfig.LoadConfig()
		if err != nil {
			return "", err
		}

		if cfg.APIKey != "" {
			flags.AgentApiKey = cfg.APIKey
		}

		// if there is a config to the collector, it should preceed the flags
		// otherwise, sync flags and config
		if cfg.CollectorEndpoint != "" {
			flags.CollectorEndpoint = cfg.CollectorEndpoint
		} else {
			cfg.CollectorEndpoint = flags.CollectorEndpoint
		}

		if cfg.Mode != "" {
			flags.Mode = agentConfig.Mode(cfg.Mode)
		}

		err = agentRunner.Run(ctx, cliConfig, flags, verbose)
		return "", err
	})),
	PostRun: teardownCommand,
}

func init() {
	startCmd.Flags().StringVarP(&startParams.organizationID, "organization", "", "", "organization id")
	startCmd.Flags().StringVarP(&startParams.environmentID, "environment", "", "", "environment id")
	startCmd.Flags().StringVarP(&startParams.agentApiKey, "api-key", "", defaultAPIKey, "agent api key")
	startCmd.Flags().StringVarP(&startParams.token, "token", "", defaultToken, "token authentication key")
	startCmd.Flags().StringVarP(&overrideEndpoint, "endpoint", "e", defaultEndpoint, "set the value for the endpoint, so the CLI won't ask for this value")
	startCmd.Flags().StringVarP(&startParams.mode, "mode", "m", "desktop", "set how the agent will start")
	startCmd.Flags().StringVarP(&startParams.logLevel, "log-level", "l", "debug", "set the agent log level")
	startCmd.Flags().StringVarP(&startParams.collectorEndpoint, "collector-endpoint", "", "", "address of the OTel Collector endpoint")

	startCmd.Flags().MarkDeprecated("endpoint", "use --server-url instead")
	startCmd.Flags().MarkShorthandDeprecated("e", "use --server-url instead")

	rootCmd.AddCommand(startCmd)
}

type startParameters struct {
	organizationID    string
	environmentID     string
	agentApiKey       string
	token             string
	mode              string
	logLevel          string
	collectorEndpoint string
}
