package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var definitionFile string

var applyCmd = &cobra.Command{
	GroupID: cmdGroupResources.ID,
	Use:     "apply [resource type]",
	Short:   "Apply resources",
	Long:    "Apply (create/update) resources to your Tracetest server",
	PreRun:  setupCommand(),
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if definitionFile == "" {
			cliLogger.Error("file with definition must be specified")
			os.Exit(1)
			return
		}

		resourceType := args[0]
		ctx := context.Background()

		analytics.Track("Resource Apply", "cmd", map[string]string{
			resourceType: resourceType,
		})

		resourceActions, err := resourceRegistry.Get(resourceType)

		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to get resource instance for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}

		applyArgs := actions.ApplyArgs{
			File: definitionFile,
		}

		resource, err := resourceActions.Apply(ctx, applyArgs)
		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to apply definition for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}

		cmd.Println(pterm.FgGreen.Sprintf(fmt.Sprintf("✔  Definition applied successfully for resource type: %s", resourceType)))

		resourceFormatter := resourceActions.Formatter()
		formatter := formatters.BuildFormatter(output, resourceFormatter.ToTable, resourceFormatter.ToStruct)

		result, err := formatter.Format(resource)
		if err != nil {
			cliLogger.Error("failed to format resource", zap.Error(err))
			os.Exit(1)
			return
		}

		fmt.Println(result)
	},
	PostRun: teardownCommand,
}

func init() {
	applyCmd.Flags().StringVarP(&definitionFile, "file", "f", "", "file path with name where to export the resource")
	rootCmd.AddCommand(applyCmd)
}
