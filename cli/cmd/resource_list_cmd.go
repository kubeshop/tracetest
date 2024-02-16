package cmd

import (
	"context"

	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/spf13/cobra"
)

var (
	listParams = &listParameters{}
	listCmd    *cobra.Command
)

func init() {
	listCmd = &cobra.Command{
		GroupID: cmdGroupResources.ID,
		Use:     "list " + resourceList(),
		Short:   "List resources",
		Long:    "List resources from your Tracetest server",
		PreRun:  setupCommand(),
		Run: WithResourceMiddleware(func(ctx context.Context, _ *cobra.Command, args []string) (string, error) {
			resourceType := resourceParams.ResourceName

			resourceClient, err := resources.Get(resourceType)
			if err != nil {
				return "", err
			}

			resultFormat, err := resourcemanager.Formats.GetWithFallback(output, "pretty")
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
				return "", err
			}

			return result, nil
		}, listParams),
		PostRun: teardownCommand,
	}

	listCmd.Flags().Int32Var(&listParams.Take, "take", 20, "Take number")
	listCmd.Flags().Int32Var(&listParams.Skip, "skip", 0, "Skip number")
	listCmd.Flags().StringVar(&listParams.SortBy, "sortBy", "", "Sort by")
	listCmd.Flags().StringVar(&listParams.SortDirection, "sortDirection", "desc", "Sort direction")
	listCmd.Flags().BoolVar(&listParams.All, "all", false, "All")

	rootCmd.AddCommand(listCmd)
}

type listParameters struct {
	Take          int32
	Skip          int32
	SortBy        string
	SortDirection string
	All           bool
}

func (p listParameters) Validate(cmd *cobra.Command, args []string) []error {
	errors := make([]error, 0)

	if p.Take < 0 {
		errors = append(errors, paramError{
			Parameter: "take",
			Message:   "take parameter must be greater than 0",
		})
	}

	if p.Skip < 0 {
		errors = append(errors, paramError{
			Parameter: "skip",
			Message:   "skip parameter must be greater than 0",
		})
	}

	if p.SortDirection != "" && p.SortDirection != "asc" && p.SortDirection != "desc" {
		errors = append(errors, paramError{
			Parameter: "sortDirection",
			Message:   "sortDirection parameter must be either asc or desc",
		})
	}

	return errors
}
