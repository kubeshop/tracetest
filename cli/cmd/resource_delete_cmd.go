package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/spf13/cobra"
)

var (
	deleteParams = &resourceIDParameters{}
	deleteCmd    *cobra.Command
)

func init() {
	deleteCmd = &cobra.Command{
		GroupID: cmdGroupResources.ID,
		Use:     "delete " + resourceList(),
		Short:   "Delete resources",
		Long:    "Delete resources from your Tracetest server",
		PreRun:  setupCommand(),
		Run: WithResourceMiddleware(func(ctx context.Context, _ *cobra.Command, args []string) (string, error) {
			resourceType := resourceParams.ResourceName

			resourceClient, err := resources.Get(resourceType)
			if err != nil {
				return "", err
			}

			resultFormat, err := resourcemanager.Formats.GetWithFallback(output, "yaml")
			if err != nil {
				return "", err
			}

			result, err := resourceClient.Delete(ctx, deleteParams.ResourceID, resultFormat)
			if errors.Is(err, resourcemanager.ErrNotFound) {
				return "", errors.New(result)
			}
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
