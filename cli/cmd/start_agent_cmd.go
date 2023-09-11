package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var (
	startAgentParams = &startAgentParameters{}
)

var startAgentCmd = &cobra.Command{
	Use:    "agent",
	Short:  "Start an Agent",
	Long:   "Start an Environment Agent",
	PreRun: setupCommand(),
	Run: WithResultHandler(WithParamsHandler(startAgentParams)(func(_ *cobra.Command, _ []string) (string, error) {
		ctx := context.Background()

		endpoint := cliConfig.AgentEndpoint
		if startAgentParams.endpoint != "" {
			endpoint = startAgentParams.endpoint
		}

		err := start.StartAgent(ctx, endpoint, startAgentParams.name, startAgentParams.apiKey, cliConfig.UIEndpoint)
		return "", err
	})),
	PostRun: teardownCommand,
}

func init() {
	if isCloudEnabled {
		startAgentCmd.Flags().StringVarP(&startAgentParams.apiKey, "api-key", "", "", "set the agent api key")
		startAgentCmd.Flags().StringVarP(&startAgentParams.endpoint, "endpoint", "", "", "set a custom server endpoint")
		startAgentCmd.Flags().StringVarP(&startAgentParams.name, "name", "", "default", "set the agent name")
		startCmd.AddCommand(startAgentCmd)
	}
}

type startAgentParameters struct {
	apiKey   string
	endpoint string
	name     string
}

func (p startAgentParameters) Validate(cmd *cobra.Command, args []string) []error {
	var errors []error

	if p.apiKey == "" {
		errors = append(errors, paramError{
			Parameter: "api-key",
			Message:   "The Agent API Key is required.",
		})
	}

	return errors
}
