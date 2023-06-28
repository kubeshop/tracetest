package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/parameters"
	"github.com/spf13/cobra"
)

var (
	exportResourceID   string
	exportResourceFile string
)

var exportCmd = &cobra.Command{
	GroupID: cmdGroupResources.ID,
	Use:     fmt.Sprintf("export %s", strings.Join(parameters.ValidResources, "|")),
	Long:    "Export a resource from your Tracetest server",
	Short:   "Export resource",
	PreRun:  setupCommand(),
	Run: WithResultHandler(func(_ *cobra.Command, args []string) (string, error) {
		resourceType := resourceParams.ResourceName
		ctx := context.Background()

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
