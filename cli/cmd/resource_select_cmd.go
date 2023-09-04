package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/spf13/cobra"
)

type selectableFn func(ctx context.Context, cfg config.Config)

var (
	selectParams  = &resourceIDParameters{}
	selectCmd     *cobra.Command
	selectable    = strings.Join([]string{"organization"}, "|")
	selectableMap = map[string]selectableFn{
		"organization": func(ctx context.Context, cfg config.Config) {
			configurator.ShowOrganizationSelector(ctx, cfg)
		}}
)

func init() {
	selectCmd = &cobra.Command{
		GroupID: cmdGroupCloud.ID,
		Use:     "select " + selectable,
		Short:   "select resources",
		Long:    "Select resources to your Tracetest CLI config",
		PreRun:  setupCommand(),
		Run: WithResourceMiddleware(func(_ *cobra.Command, args []string) (string, error) {
			resourceType := resourceParams.ResourceName
			ctx := context.Background()

			selectableFn, ok := selectableMap[resourceType]
			if !ok {
				return "", fmt.Errorf("resource type %s not selectable. Selectable resources are %s", resourceType, selectable)
			}

			resourceClient, err := resources.Get(resourceType)
			if err != nil {
				return "", err
			}

			resultFormat, err := resourcemanager.Formats.GetWithFallback(output, "yaml")
			if err != nil {
				return "", err
			}

			result, err := resourceClient.Get(ctx, selectParams.ResourceID, resultFormat)
			if errors.Is(err, resourcemanager.ErrNotFound) {
				return result, nil
			}
			if err != nil {
				return "", err
			}

			selectableFn(ctx, cliConfig)
			return "", nil
		}),
		PostRun: teardownCommand,
	}

	rootCmd.AddCommand(selectCmd)
}
