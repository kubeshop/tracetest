package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/spf13/cobra"
)

var definitionFile string

var applyCmd = &cobra.Command{
	GroupID:   cmdGroupResources.ID,
	Use:       fmt.Sprintf("apply %s", strings.Join(validArgs, "|")),
	Short:     "Apply resources",
	Long:      "Apply (create/update) resources to your Tracetest server",
	PreRun:    setupCommand(),
	Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	ValidArgs: validArgs,
	Run: WithResultHandler(func(cmd *cobra.Command, args []string) (string, error) {
		if definitionFile == "" {
			return "", fmt.Errorf("file with definition must be specified")
		}

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
			File: definitionFile,
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
	}),
	PostRun: teardownCommand,
}

func init() {
	applyCmd.Flags().StringVarP(&definitionFile, "file", "f", "", "file path with name where to export the resource")
	rootCmd.AddCommand(applyCmd)
}
