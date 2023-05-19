package resources

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/global"
	global_decorators "github.com/kubeshop/tracetest/cli/global/decorators"
	global_formatters "github.com/kubeshop/tracetest/cli/global/formatters"
	resources_actions "github.com/kubeshop/tracetest/cli/resources/actions"
	resources_decorators "github.com/kubeshop/tracetest/cli/resources/decorators"
	resources_formatters "github.com/kubeshop/tracetest/cli/resources/formatters"
	resources_parameters "github.com/kubeshop/tracetest/cli/resources/parameters"
	"github.com/spf13/cobra"
)

type apply struct {
	resources_decorators.Resources
	Parameters *resources_parameters.Apply
}

type Apply interface {
	Resource
}

var _ Apply = &apply{}

func NewApply() Apply {
	parameters := resources_parameters.NewApply()
	apply := &apply{
		Parameters: parameters,
	}

	cmd := &cobra.Command{
		GroupID: global.GroupResources.ID,
		Use:     "apply [resource type]",
		Short:   "Apply resources",
		Long:    "Apply (create/update) resources to your Tracetest server",
		Args:    cobra.MinimumNArgs(1),
		PreRun:  global_decorators.NoopRun,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("apply")
		},
		PostRun: global_decorators.NoopRun,
	}
	cmd.Flags().StringVarP(&apply.Parameters.DefinitionFile, "file", "f", "", "file path with name where to export the resource")

	return global_decorators.Decorate[Apply](
		apply,
		resources_decorators.WithResources,
	)
}

func (a *apply) Run(cmd *cobra.Command, args []string) {
	if a.Parameters.DefinitionFile == "" {
		a.Error(fmt.Errorf("file with definition must be specified"))
		return
	}

	ctx := context.Background()

	applyArgs := resources_actions.ApplyArgs{
		File: a.Parameters.DefinitionFile,
	}

	resourceActions := a.GetResourceActions()

	resource, _, err := resourceActions.Apply(ctx, applyArgs)
	if err != nil {
		a.Error(err)
		return
	}

	resourceFormatter := resourceActions.Formatter()
	formatter := resources_formatters.BuildFormatter(a.Parameters.Output, global_formatters.YAML, resourceFormatter)

	result, err := formatter.Format(resource)
	if err != nil {
		a.Error(err)
		return
	}

	a.Response(result)
}
