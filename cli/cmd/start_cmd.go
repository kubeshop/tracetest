package cmd

import (
	"context"

	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/pkg/starter"
	"github.com/spf13/cobra"
)

var (
	start      = starter.NewStarter(configurator, resources)
	saveParams = &saveParameters{}
)

var startCmd = &cobra.Command{
	GroupID: cmdGroupConfig.ID,
	Use:     "start",
	Short:   "Start Tracetest",
	Long:    "Start using Tracetest",
	PreRun:  setupCommand(SkipVersionMismatchCheck()),
	Run: WithResultHandler((func(_ *cobra.Command, _ []string) (string, error) {
		ctx := context.Background()

		flags := config.ConfigFlags{
			OrganizationID: saveParams.organizationID,
			EnvironmentID:  saveParams.environmentID,
			Endpoint:       saveParams.endpoint,
			AgentApiKey:    saveParams.agentApiKey,
		}

		cfg, err := agentConfig.LoadConfig()
		if err != nil {
			return "", err
		}

		if cfg.APIKey != "" {
			flags.AgentApiKey = cfg.APIKey
		}

		err = start.Run(ctx, cliConfig, flags)
		return "", err
	})),
	PostRun: teardownCommand,
}

func init() {
	startCmd.Flags().StringVarP(&saveParams.organizationID, "organization", "", "", "organization id")
	startCmd.Flags().StringVarP(&saveParams.environmentID, "environment", "", "", "environment id")
	startCmd.Flags().StringVarP(&saveParams.agentApiKey, "api-key", "", "", "agent api key")
	startCmd.Flags().StringVarP(&saveParams.endpoint, "endpoint", "e", config.DefaultCloudEndpoint, "set the value for the endpoint, so the CLI won't ask for this value")
	rootCmd.AddCommand(startCmd)
}

type saveParameters struct {
	organizationID string
	environmentID  string
	endpoint       string
	agentApiKey    string
}
