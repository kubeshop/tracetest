package cmd

import (
	"context"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/pkg/starter"
	"github.com/spf13/cobra"
)

var (
	start = starter.NewStarter(configurator, resources)
)

var startCmd = &cobra.Command{
	GroupID: cmdGroupCloud.ID,
	Use:     "start",
	Short:   "Start Tracetest",
	Long:    "Start using Tracetest",
	PreRun:  setupCommand(),
	Run: WithResultHandler((func(_ *cobra.Command, _ []string) (string, error) {
		ctx := context.Background()

		flags := config.ConfigFlags{
			OrganizationID: selectParams.organizationID,
			EnvironmentID:  selectParams.environmentID,
			Endpoint:       selectParams.endpoint,
		}

		err := start.Run(ctx, cliConfig, flags)
		return "", err
	})),
	PostRun: teardownCommand,
}

func init() {
	if isCloudEnabled {
		startCmd.Flags().StringVarP(&selectParams.organizationID, "organization", "", "", "organization id")
		startCmd.Flags().StringVarP(&selectParams.environmentID, "environment", "", "", "environment id")
		startCmd.Flags().StringVarP(&selectParams.endpoint, "endpoint", "e", "", "set the value for the endpoint, so the CLI won't ask for this value")
		rootCmd.AddCommand(startCmd)
	}
}
