package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var exportResourceID string

var exportCmd = &cobra.Command{
	Use:    "export [resource type]",
	Long:   "Export a resource from your Tracetest server",
	Short:  "Export resource",
	PreRun: setupCommand(),
	Args:   cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]
		ctx := context.Background()

		analytics.Track("Resource Export", "cmd", map[string]string{
			resourceType: resourceType,
		})

		resourceActions, err := resourceRegistry.Get(actions.SupportedResources(resourceType))
		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to export resource instance for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}

		err = resourceActions.Export(ctx, exportResourceID)
		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to export resource for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}

		cmd.Println(pterm.FgGreen.Sprintf(fmt.Sprintf("✔  Definition exported successfully for resource type: %s", resourceType)))
	},
	PostRun: teardownCommand,
}

func init() {
	exportCmd.PersistentFlags().StringVarP(&exportResourceID, "identifier", "i", "", "id of the resource to export")
	rootCmd.AddCommand(exportCmd)
}
