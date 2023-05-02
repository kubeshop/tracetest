package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var deletedResourceID string

var deleteCmd = &cobra.Command{
	GroupID: cmdGroupResources.ID,
	Use:     "delete [resource type]",
	Long:    "Delete resources from your Tracetest server",
	Short:   "Delete resources",
	PreRun:  setupCommand(),
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if deletedResourceID == "" {
			cliLogger.Error("id of the resource to delete must be specified")
			os.Exit(1)
			return
		}

		resourceType := args[0]
		ctx := context.Background()

		analytics.Track("Resource Delete", "cmd", map[string]string{
			resourceType: resourceType,
		})

		resourceActions, err := resourceRegistry.Get(resourceType)
		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to get resource instance for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}

		message, err := resourceActions.Delete(ctx, deletedResourceID)
		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to apply definition for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}

		cmd.Println(pterm.FgGreen.Sprintf(fmt.Sprintf("✔ %s", message)))
	},
	PostRun: teardownCommand,
}

func init() {
	deleteCmd.Flags().StringVar(&deletedResourceID, "id", "", "id of the resource to delete")
	rootCmd.AddCommand(deleteCmd)
}
