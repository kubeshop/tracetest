package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/parameters"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/spf13/cobra"
)

var getParams = &parameters.ResourceIdParams{}

var getCmd = &cobra.Command{
	GroupID: cmdGroupResources.ID,
	Use:     fmt.Sprintf("get %s", strings.Join(parameters.ValidResources, "|")),
	Short:   "Get resource",
	Long:    "Get a resource from your Tracetest server",
	PreRun:  setupCommand(),
	Run: WithResourceMiddleware(func(_ *cobra.Command, args []string) (string, error) {
		resourceType := args[0]
		ctx := context.Background()

		resourceClient, err := resources.Get(resourceType)
		if err != nil {
			return "", err
		}

		resultFormat, err := resourcemanager.Formats.Get(output, "yaml")
		if err != nil {
			return "", err
		}

		result, err := resourceClient.Get(ctx, getParams.ResourceID, resultFormat)
		if err != nil {
			return "", err
		}

		return result, nil
	}, getParams),
	PostRun: teardownCommand,
}

func init() {
	getCmd.Flags().StringVar(&getParams.ResourceID, "id", "", "id of the resource to get")
	rootCmd.AddCommand(getCmd)
}
