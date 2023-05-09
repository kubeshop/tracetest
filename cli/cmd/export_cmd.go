package cmd

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/spf13/cobra"
)

var (
	exportResourceID   string
	exportResourceFile string
)

var exportCmd = &cobra.Command{
	GroupID: cmdGroupResources.ID,
	Use:     "export [resource type]",
	Long:    "Export a resource from your Tracetest server",
	Short:   "Export resource",
	PreRun:  setupCommand(),
	Args:    cobra.MinimumNArgs(1),
	Run: WithResultHandler(func(cmd *cobra.Command, args []string) (string, error) {
		resourceType := args[0]
		ctx := context.Background()

		analytics.Track("Resource Export", "cmd", map[string]string{
			resourceType: resourceType,
		})

		resourceActions, err := resourceRegistry.Get(resourceType)
		if err != nil {
			return "", err
		}

		err = resourceActions.Export(ctx, exportResourceID, exportResourceFile)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("✔  Definition exported successfully for resource type: %s", resourceType), nil
	}),
	PostRun: teardownCommand,
}

func init() {
	exportCmd.Flags().StringVar(&exportResourceID, "id", "", "id of the resource to export")
	exportCmd.Flags().StringVarP(&exportResourceFile, "file", "f", "resource.yaml", "file path with name where to export the resource")
	rootCmd.AddCommand(exportCmd)
}
