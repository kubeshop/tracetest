package cmd

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/spf13/cobra"
)

var deletedResourceID string

var deleteCmd = &cobra.Command{
	GroupID: cmdGroupResources.ID,
	Use:     "delete [resource type]",
	Long:    "Delete resources from your Tracetest server",
	Short:   "Delete resources",
	PreRun:  setupCommand(),
	Args:    cobra.MinimumNArgs(1),
	Run: WithResultHandler(func(_ *cobra.Command, args []string) (string, error) {
		if deletedResourceID == "" {
			return "", fmt.Errorf("id of the resource to delete must be specified")
		}

		resourceType := args[0]
		ctx := context.Background()

		analytics.Track("Resource Delete", "cmd", map[string]string{
			resourceType: resourceType,
		})

		resourceActions, err := resourceRegistry.Get(resourceType)
		if err != nil {
			return "", err
		}

		message, err := resourceActions.Delete(ctx, deletedResourceID)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("✔ %s", message), nil
	}),
	PostRun: teardownCommand,
}

func init() {
	deleteCmd.Flags().StringVar(&deletedResourceID, "id", "", "id of the resource to delete")
	rootCmd.AddCommand(deleteCmd)
}
