package cmd

import (
	"context"

	"github.com/kubeshop/tracetest/cli/parameters"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/spf13/cobra"
)

var applyParams = &parameters.ApplyParams{}

var applyCmd = &cobra.Command{
	GroupID: cmdGroupResources.ID,
	Use:     "apply " + resourceList(),
	Short:   "Apply resources",
	Long:    "Apply (create/update) resources to your Tracetest server",
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

		result, err := resourceClient.Apply(ctx, applyParams.DefinitionFile, resultFormat)
		if err != nil {
			return "", err
		}

		return result, nil
	}, applyParams),
	PostRun: teardownCommand,
}

func init() {
	applyCmd.Flags().StringVarP(&applyParams.DefinitionFile, "file", "f", "", "path to the definition file")
	rootCmd.AddCommand(applyCmd)
}
