package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/parameters"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

var listParams = &parameters.ListParams{}

var listCmd = &cobra.Command{
	GroupID: cmdGroupResources.ID,
	Use:     fmt.Sprintf("list %s", strings.Join(parameters.ValidResources, "|")),
	Short:   "List resources",
	Long:    "List resources from your Tracetest server",
	PreRun:  setupCommand(),
	Run: WithResourceMiddleware(func(_ *cobra.Command, args []string) (string, error) {
		resourceType := args[0]
		ctx := context.Background()

		resourceActions, err := resourceRegistry.Get(resourceType)
		if err != nil {
			return "", err
		}

		listArgs := utils.ListArgs{
			Take:          listParams.Take,
			Skip:          listParams.Skip,
			SortDirection: listParams.SortDirection,
			SortBy:        listParams.SortBy,
			All:           listParams.All,
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
	}, listParams),
	PostRun: teardownCommand,
}

func init() {
	listCmd.Flags().Int32Var(&listParams.Take, "take", 20, "Take number")
	listCmd.Flags().Int32Var(&listParams.Skip, "skip", 0, "Skip number")
	listCmd.Flags().StringVar(&listParams.SortBy, "sortBy", "", "Sort by")
	listCmd.Flags().StringVar(&listParams.SortDirection, "sortDirection", "desc", "Sort direction")
	listCmd.Flags().BoolVar(&listParams.All, "all", false, "All")

	rootCmd.AddCommand(listCmd)
}
