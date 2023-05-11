package resources

import (
	"context"

	"github.com/kubeshop/tracetest/cli/global"
	global_formatters "github.com/kubeshop/tracetest/cli/global/formatters"
	resources_formatters "github.com/kubeshop/tracetest/cli/resources/formatters"
	resource_parameters "github.com/kubeshop/tracetest/cli/resources/parameters"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

type List struct {
	args[*resource_parameters.List]
}

func NewList(root global.Root) List {
	parameters := resource_parameters.NewList()
	defaults := NewDefaults("list", root.Setup)

	list := List{
		args: NewArgs(defaults, parameters),
	}

	list.Cmd = &cobra.Command{
		GroupID: global.GroupResources.ID,
		Use:     "list [resource type]",
		Long:    "List resources from your Tracetest server",
		Short:   "List resources",
		Args:    cobra.MinimumNArgs(1),
		PreRun:  defaults.PreRun,
		Run:     defaults.Run(list.Run),
		PostRun: defaults.PostRun,
	}

	list.Cmd.Flags().Int32Var(&list.Parameters.Take, "take", 20, "Take number")
	list.Cmd.Flags().Int32Var(&list.Parameters.Skip, "skip", 0, "Skip number")
	list.Cmd.Flags().StringVar(&list.Parameters.SortBy, "sortBy", "", "Sort by")
	list.Cmd.Flags().StringVar(&list.Parameters.SortDirection, "sortDirection", "desc", "Sort direction")
	list.Cmd.Flags().BoolVar(&list.Parameters.All, "all", false, "All")

	root.Cmd.AddCommand(list.Cmd)

	return list
}

func (l List) Run(cmd *cobra.Command, args []string) (string, error) {
	ctx := context.Background()

	listArgs := utils.ListArgs{
		Take:          l.Parameters.Take,
		Skip:          l.Parameters.Skip,
		SortDirection: l.Parameters.SortDirection,
		SortBy:        l.Parameters.SortBy,
		All:           l.Parameters.All,
	}

	resource, err := l.Setup.ResourceActions.List(ctx, listArgs)
	if err != nil {
		return "", err
	}

	resourceFormatter := l.Setup.ResourceActions.Formatter()
	formatter := resources_formatters.BuildFormatter(l.Setup.Parameters.Output, global_formatters.Pretty, resourceFormatter)

	result, err := formatter.FormatList(resource)
	if err != nil {
		return "", err
	}

	return result, nil
}
