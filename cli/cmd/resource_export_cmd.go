package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/spf13/cobra"
)

var (
	exportParams = &exportParameters{}
	exportCmd    *cobra.Command
)

func init() {
	exportCmd = &cobra.Command{
		GroupID: cmdGroupResources.ID,
		Use:     "export " + resourceList(),
		Long:    "Export a resource from your Tracetest server",
		Short:   "Export resource",
		PreRun:  setupCommand(),
		Run: WithResourceMiddleware(func(_ *cobra.Command, args []string) (string, error) {
			resourceType := resourceParams.ResourceName
			ctx := context.Background()

			resourceClient, err := resources.Get(resourceType)
			if err != nil {
				return "", err
			}

			// export is ALWAYS yaml, so we can hardcode it here
			result, err := resourceClient.Get(ctx, exportParams.ResourceID, resourcemanager.Formats.Get("yaml"))
			if err != nil {
				return "", err
			}

			err = os.WriteFile(exportParams.OutputFile, []byte(result), 0644)
			if err != nil {
				return "", fmt.Errorf("could not write file: %w", err)
			}

			return fmt.Sprintf("✔  Definition exported successfully for resource type: %s", resourceType), nil
		}, exportParams),
		PostRun: teardownCommand,
	}

	exportCmd.Flags().StringVar(&exportParams.ResourceID, "id", "", "id of the resource to export")
	exportCmd.Flags().StringVarP(&exportParams.OutputFile, "file", "f", "resource.yaml", "file path with name where to export the resource")
	rootCmd.AddCommand(exportCmd)
}

type exportParameters struct {
	resourceIDParameters
	OutputFile string
}

func (p exportParameters) Validate(cmd *cobra.Command, args []string) []error {
	errors := p.resourceIDParameters.Validate(cmd, args)

	if p.OutputFile == "" {
		errors = append(errors, paramError{
			Parameter: "file",
			Message:   "output file must be provided",
		})
	}

	return errors
}
