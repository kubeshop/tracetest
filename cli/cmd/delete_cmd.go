package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/parameters"
	"github.com/spf13/cobra"
)

var deleteParams = &parameters.ResourceIdParams{}

var deleteCmd = &cobra.Command{
	GroupID: cmdGroupResources.ID,
	Use:     fmt.Sprintf("delete %s", strings.Join(parameters.ValidResources, "|")),
	Short:   "Delete resources",
	Long:    "Delete resources from your Tracetest server",
	PreRun:  setupCommand(),
	Run: WithResourceMiddleware(func(_ *cobra.Command, args []string) (string, error) {
		resourceType := args[0]
		ctx := context.Background()

		resourceActions, err := resourceRegistry.Get(resourceType)
		if err != nil {
			return "", err
		}

		message, err := resourceActions.Delete(ctx, deleteParams.ResourceID)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("✔ %s", message), nil
	}, deleteParams),
	PostRun: teardownCommand,
}

func init() {
	deleteCmd.Flags().StringVar(&deleteParams.ResourceID, "id", "", "id of the resource to delete")
	rootCmd.AddCommand(deleteCmd)
}
