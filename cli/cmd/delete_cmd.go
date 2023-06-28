package cmd

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/parameters"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/spf13/cobra"
)

var (
	deleteParams = &parameters.ResourceIdParams{}
	deleteCmd    *cobra.Command
)

func init() {
	deleteCmd = &cobra.Command{
		GroupID: cmdGroupResources.ID,
		Use:     "delete " + resourceList(),
		Short:   "Delete resources",
		Long:    "Delete resources from your Tracetest server",
		PreRun:  setupCommand(),
		Run: WithResourceMiddleware(func(_ *cobra.Command, args []string) (string, error) {
			resourceType := resourceParams.ResourceName
			ctx := context.Background()

			resourceClient, err := resources.Get(resourceType)
			if err != nil {
				return "", err
			}

			resultFormat, err := resourcemanager.Formats.Get(output, "yaml")
			if err != nil {
				return "", err
			}

			result, err := resourceClient.Delete(ctx, deleteParams.ResourceID, resultFormat)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("✔ %s", result), nil
		}, deleteParams),
		PostRun: teardownCommand,
	}

	deleteCmd.Flags().StringVar(&deleteParams.ResourceID, "id", "", "id of the resource to delete")
	rootCmd.AddCommand(deleteCmd)
}
