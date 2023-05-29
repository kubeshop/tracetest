package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/parameters"
	"github.com/spf13/cobra"
)

var applyParams = &parameters.ApplyParams{}

var applyCmd = &cobra.Command{
	GroupID: cmdGroupResources.ID,
	Use:     fmt.Sprintf("apply %s", strings.Join(parameters.ValidResources, "|")),
	Short:   "Apply resources",
	Long:    "Apply (create/update) resources to your Tracetest server",
	PreRun:  setupCommand(),
	Run: WithResourceMiddleware(func(_ *cobra.Command, args []string) (string, error) {
		resourceType := args[0]
		ctx := context.Background()

		analytics.Track("Resource Apply", "cmd", map[string]string{
			resourceType: resourceType,
		})

		resourceActions, err := resourceRegistry.Get(resourceType)

		if err != nil {
			return "", err
		}

		applyArgs := actions.ApplyArgs{
			File: applyParams.DefinitionFile,
		}

		resource, _, err := resourceActions.Apply(ctx, applyArgs)
		if err != nil {
			return "", err
		}

		resourceFormatter := resourceActions.Formatter()
		formatter := formatters.BuildFormatter(output, formatters.YAML, resourceFormatter)

		result, err := formatter.Format(resource)
		if err != nil {
			return "", err
		}

		return result, nil
	}, applyParams),
	PostRun: teardownCommand,
}

func init() {
	applyCmd.Flags().StringVarP(&applyParams.DefinitionFile, "file", "f", "", "file path with name where to export the resource")
	rootCmd.AddCommand(applyCmd)
}
