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

var (
	listTake          int32
	listSkip          int32
	listSortBy        string
	listSortDirection string
)

var listCmd = &cobra.Command{
	Use:    "list [resource type]",
	Long:   "List resources from your Tracetest server",
	Short:  "List resources",
	PreRun: setupCommand(),
	Args:   cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]
		ctx := context.Background()

		analytics.Track("Resource List", "cmd", map[string]string{
			resourceType: resourceType,
		})

		resourceActions, err := resourceRegistry.Get(actions.SupportedResources(resourceType))

		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to get resource instance for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}

		listArgs := actions.ListArgs{
			Take:          listTake,
			Skip:          listSkip,
			SortDirection: listSortDirection,
			SortBy:        listSortBy,
		}

		err = resourceActions.List(ctx, listArgs)

		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to list for type: %s", resourceType), zap.Error(err))
			os.Exit(1)
			return
		}
	},
	PostRun: teardownCommand,
}

func init() {
	listCmd.PersistentFlags().Int32VarP(&listTake, "take", "t", 20, "Take number")
	listCmd.PersistentFlags().Int32VarP(&listSkip, "skip", "", 0, "Skip number")
	listCmd.PersistentFlags().StringVarP(&listSortBy, "sortBy", "b", "", "Sort by")
	listCmd.PersistentFlags().StringVarP(&listSortDirection, "sortDirection", "d", "desc", "Sort direction")
	rootCmd.AddCommand(listCmd)
}
