package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var resourceID string

var getCmd = &cobra.Command{
	GroupID:   cmdGroupResources.ID,
	Use:       fmt.Sprintf("get %s", strings.Join(validArgs, "|")),
	Short:     "Get resource",
	Long:      "Get a resource from your Tracetest server",
	PreRun:    setupCommand(),
	Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	ValidArgs: validArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if resourceID == "" {
			cliLogger.Error("id of the resource to get must be specified")
			os.Exit(1)
			return
		}

		resourceType := args[0]
		ctx := context.Background()

		analytics.Track("Resource Get", "cmd", map[string]string{
			resourceType: resourceType,
		})

		resourceActions, err := resourceRegistry.Get(resourceType)

		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to get resource instance for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}

		resource, err := resourceActions.Get(ctx, resourceID)
		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to get resource for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}

		resourceFormatter := resourceActions.Formatter()
		formatter := formatters.BuildFormatter(output, formatters.YAML, resourceFormatter)

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
	getCmd.Flags().StringVar(&resourceID, "id", "", "id of the resource to get")
	rootCmd.AddCommand(getCmd)
}
