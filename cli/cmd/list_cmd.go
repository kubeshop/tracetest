package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/parameters"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/spf13/cobra"
)

var listParams = parameters.ListParams{}

var listCmd = &cobra.Command{
	GroupID: cmdGroupResources.ID,
	Use:     fmt.Sprintf("list %s", strings.Join(parameters.ValidResources, "|")),
	Short:   "List resources",
	Long:    "List resources from your Tracetest server",
	PreRun:  setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]
		ctx := context.Background()

		resourceClient, err := resources.Get(resourceType)
		if err != nil {
			panic(err)
		}

		resultFormat, err := resourcemanager.Formats.Get(output)
		if err != nil {
			return "", err
		}

		lp := resourcemanager.ListOption{
			Take:          listParams.Take,
			Skip:          listParams.Skip,
			SortBy:        listParams.SortBy,
			SortDirection: listParams.SortDirection,
			All:           listParams.All,
		}

		result, err := resourceClient.List(ctx, lp, resultFormat)
		if err != nil {
			panic(err)
		}

		fmt.Println(result)

	},
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
