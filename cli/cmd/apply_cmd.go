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

var definitionFile string

var applyCmd = &cobra.Command{
	Use:    "apply [resource type]",
	Short:  "Apply resources",
	Long:   "Apply (create/update) resources to your Tracetest server",
	PreRun: setupCommand(),
	Args:   cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]
		ctx := context.Background()

		analytics.Track("Resource Apply", "cmd", map[string]string{
			resourceType: resourceType,
		})

		resourceActions, err := resourceRegistry.Get(actions.SupportedResources(resourceType))

		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to get resource instance for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}

		applyArgs := actions.ApplyArgs{
			File: definitionFile,
		}

		err = resourceActions.Apply(ctx, applyArgs)

		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to apply definition for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}

		cmd.Println(pterm.FgGreen.Sprintf(fmt.Sprintf("✔  Definition applied successfully for resource type: %s", resourceType)))
	},
	PostRun: teardownCommand,
}

func init() {
	applyCmd.Flags().StringVar(&definitionFile, "file", "", "path to the definition file")
	rootCmd.AddCommand(applyCmd)
}
