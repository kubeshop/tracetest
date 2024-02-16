package cmd

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/spf13/cobra"
)

var (
	applyParams = &applyParameters{}
	applyCmd    *cobra.Command
)

func init() {
	applyCmd = &cobra.Command{
		GroupID: cmdGroupResources.ID,
		Use:     "apply " + resourceList(),
		Short:   "Apply resources",
		Long:    "Apply (create/update) resources to your Tracetest server",
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

			inputFile, err := fileutil.Read(applyParams.DefinitionFile)
			if err != nil {
				return "", fmt.Errorf("cannot read file %s: %w", applyParams.DefinitionFile, err)
			}

			result, err := resourceClient.Apply(ctx, inputFile, resultFormat)
			if err != nil {
				return "", err
			}

			return result, nil
		}, applyParams),
		PostRun: teardownCommand,
	}

	applyCmd.Flags().StringVarP(&applyParams.DefinitionFile, "file", "f", "", "path to the definition file")
	rootCmd.AddCommand(applyCmd)
}

type applyParameters struct {
	DefinitionFile string
}

func (p applyParameters) Validate(cmd *cobra.Command, args []string) []error {
	errors := make([]error, 0)

	if p.DefinitionFile == "" {
		errors = append(errors, paramError{
			Parameter: "file",
			Message:   "Definition file must be provided",
		})
	}

	return errors
}
