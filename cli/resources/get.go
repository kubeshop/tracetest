package resources

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/global"
	global_formatters "github.com/kubeshop/tracetest/cli/global/formatters"
	resources_formatters "github.com/kubeshop/tracetest/cli/resources/formatters"
	resource_parameters "github.com/kubeshop/tracetest/cli/resources/parameters"
	"github.com/spf13/cobra"
)

type Get struct {
	args[*resource_parameters.Id]
}

func NewGet(root global.Root) Get {
	parameters := resource_parameters.NewId()
	defaults := NewDefaults("get", root.Setup)

	get := Get{
		args: NewArgs(defaults, parameters),
	}

	get.Cmd = &cobra.Command{
		GroupID: global.GroupResources.ID,
		Use:     "get [resource type]",
		Long:    "Get a resource from your Tracetest server",
		Short:   "Get resource",
		Args:    cobra.MinimumNArgs(1),
		PreRun:  defaults.PreRun,
		Run:     defaults.Run(get.Run),
		PostRun: defaults.PostRun,
	}

	get.Cmd.Flags().StringVar(&get.Parameters.ResourceID, "id", "", "id of the resource to get")
	root.Cmd.AddCommand(get.Cmd)

	return get
}

func (g Get) Run(cmd *cobra.Command, args []string) (string, error) {
	if g.Parameters.ResourceID == "" {
		return "", fmt.Errorf("id of the resource to get must be specified")
	}

	ctx := context.Background()
	resource, err := g.Setup.ResourceActions.Get(ctx, g.Parameters.ResourceID)
	if err != nil {
		return "", err
	}

	resourceFormatter := g.Setup.ResourceActions.Formatter()
	formatter := resources_formatters.BuildFormatter(g.Setup.Parameters.Output, global_formatters.YAML, resourceFormatter)
	result, err := formatter.Format(resource)
	if err != nil {
		return "", err
	}

	return result, nil
}
