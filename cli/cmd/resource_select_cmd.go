package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/spf13/cobra"
)

type selectableFn func(ctx context.Context, cfg config.Config, flags config.ConfigFlags)

var (
	selectParams  = &selectParameters{}
	selectCmd     *cobra.Command
	selectable    = strings.Join([]string{"organization"}, "|")
	selectableMap = map[string]selectableFn{
		"organization": func(ctx context.Context, cfg config.Config, flags config.ConfigFlags) {
			configurator.ShowOrganizationSelector(ctx, cfg, flags)
		}}
)

func init() {
	selectCmd = &cobra.Command{
		GroupID: cmdGroupCloud.ID,
		Use:     "select " + selectable,
		Short:   "Select resources",
		Long:    "Select resources to your Tracetest CLI config",
		PreRun:  setupCommand(),
		Run: WithResourceMiddleware(func(_ *cobra.Command, args []string) (string, error) {
			resourceType := resourceParams.ResourceName
			ctx := context.Background()

			selectableFn, ok := selectableMap[resourceType]
			if !ok {
				return "", fmt.Errorf("resource type %s not selectable. Selectable resources are %s", resourceType, selectable)
			}

			flags := config.ConfigFlags{
				OrganizationID: selectParams.organizationID,
				EnvironmentID:  selectParams.environmentID,
			}

			selectableFn(ctx, cliConfig, flags)
			return "", nil
		}),
		PostRun: teardownCommand,
	}

	if isCloudEnabled {
		selectCmd.Flags().StringVarP(&selectParams.organizationID, "organization", "", "", "organization id")
		selectCmd.Flags().StringVarP(&selectParams.environmentID, "environment", "", "", "environment id")
		rootCmd.AddCommand(selectCmd)
	}
}

type selectParameters struct {
	organizationID string
	environmentID  string
	endpoint       string
	agentApiKey    string
}
