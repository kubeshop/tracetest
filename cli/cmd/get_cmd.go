package cmd

import (
	"context"

	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/spf13/cobra"
)

var (
	getParams = ResourceIDParameters{}
	getCmd    *cobra.Command
)

func init() {
	getCmd = &cobra.Command{
		GroupID: cmdGroupResources.ID,
		Use:     "get " + resourceList(),
		Short:   "Get resource",
		Long:    "Get a resource from your Tracetest server",
		PreRun:  setupCommand(),
		Run: WithResourceMiddleware(func(_ *cobra.Command, args []string) (string, error) {
			resourceType := resourceParams.ResourceName
			ctx := context.Background()

			resourceClient, err := resources.Get(resourceType)
			if err != nil {
				return "", err
			}

			resultFormat, err := resourcemanager.Formats.GetWithFallback(output, "yaml")
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

	getCmd.Flags().StringVar(&getParams.ResourceID, "id", "", "id of the resource to get")
	rootCmd.AddCommand(getCmd)
}

type ResourceIDParameters struct {
	ResourceID string
}

func (p ResourceIDParameters) Validate(cmd *cobra.Command, args []string) []error {
	errors := make([]error, 0)

	if p.ResourceID == "" {
		errors = append(errors, ParamError{
			Parameter: "id",
			Message:   "resource id must be provided",
		})
	}

	return errors
}
