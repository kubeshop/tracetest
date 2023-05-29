package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

var (
	listTake          int32
	listSkip          int32
	listSortBy        string
	listSortDirection string
	listAll           bool
)

var listCmd = &cobra.Command{
	GroupID:   cmdGroupResources.ID,
	Use:       fmt.Sprintf("list %s", strings.Join(validArgs, "|")),
	Short:     "List resources",
	Long:      "List resources from your Tracetest server",
	PreRun:    setupCommand(),
	Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	ValidArgs: validArgs,
	Run: WithResultHandler(func(cmd *cobra.Command, args []string) (string, error) {
		resourceType := args[0]
		ctx := context.Background()

		analytics.Track("Resource List", "cmd", map[string]string{
			resourceType: resourceType,
		})

		resourceActions, err := resourceRegistry.Get(resourceType)
		if err != nil {
			return "", err
		}

		listArgs := utils.ListArgs{
			Take:          listTake,
			Skip:          listSkip,
			SortDirection: listSortDirection,
			SortBy:        listSortBy,
			All:           listAll,
		}

		resource, err := resourceActions.List(ctx, listArgs)
		if err != nil {
			return "", err
		}

		resourceFormatter := resourceActions.Formatter()
		formatter := formatters.BuildFormatter(output, formatters.Pretty, resourceFormatter)

		result, err := formatter.FormatList(resource)
		if err != nil {
			return "", err
		}

		return result, nil
	}),
	PostRun: teardownCommand,
}

func init() {
	listCmd.Flags().Int32Var(&listTake, "take", 20, "Take number")
	listCmd.Flags().Int32Var(&listSkip, "skip", 0, "Skip number")
	listCmd.Flags().StringVar(&listSortBy, "sortBy", "", "Sort by")
	listCmd.Flags().StringVar(&listSortDirection, "sortDirection", "desc", "Sort direction")
	listCmd.Flags().BoolVar(&listAll, "all", false, "All")

	rootCmd.AddCommand(listCmd)
}
