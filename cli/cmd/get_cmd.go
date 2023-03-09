package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var resourceID string

var getCmd = &cobra.Command{
	Use:    "get [resource type]",
	Long:   "Get a resource from your Tracetest server",
	Short:  "Get resource",
	PreRun: setupCommand(),
	Args:   cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]
		ctx := context.Background()

		analytics.Track("Resource Get", "cmd", map[string]string{
			resourceType: resourceType,
		})

		resourceActions, err := resourceRegistry.Get(actions.SupportedResources(resourceType))

		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to get resource instance for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}

		err = resourceActions.Get(ctx, resourceID)

		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to get resource for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}
	},
	PostRun: teardownCommand,
}

func init() {
	getCmd.PersistentFlags().StringVarP(&resourceID, "identifier", "i", "", "id of the resource to get")
	rootCmd.AddCommand(getCmd)
}
