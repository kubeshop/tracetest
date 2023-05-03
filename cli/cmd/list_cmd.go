package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	listTake          int32
	listSkip          int32
	listSortBy        string
	listSortDirection string
)

var listCmd = &cobra.Command{
	GroupID: cmdGroupResources.ID,
	Use:     "list [resource type]",
	Long:    "List resources from your Tracetest server",
	Short:   "List resources",
	PreRun:  setupCommand(),
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]
		ctx := context.Background()

		analytics.Track("Resource List", "cmd", map[string]string{
			resourceType: resourceType,
		})

		resourceActions, err := resourceRegistry.Get(resourceType)

		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to get resource instance for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}

		listArgs := utils.ListArgs{
			Take:          listTake,
			Skip:          listSkip,
			SortDirection: listSortDirection,
			SortBy:        listSortBy,
		}

		resource, err := resourceActions.List(ctx, listArgs)
		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to list for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}

		resourceFormatter := resourceActions.Formatter()
		formatter := formatters.BuildFormatter(output, formatters.Pretty, resourceFormatter.ToListTable, resourceFormatter.ToListStruct)

		result, err := formatter.Format(resource)
		if err != nil {
			cliLogger.Error("failed to format resource", zap.Error(err))
			os.Exit(1)
			return
		}

		cmd.Println(result)
	},
	PostRun: teardownCommand,
}

func init() {
	listCmd.Flags().Int32Var(&listTake, "take", 20, "Take number")
	listCmd.Flags().Int32Var(&listSkip, "skip", 0, "Skip number")
	listCmd.Flags().StringVar(&listSortBy, "sortBy", "", "Sort by")
	listCmd.Flags().StringVar(&listSortDirection, "sortDirection", "desc", "Sort direction")

	rootCmd.AddCommand(listCmd)
}
