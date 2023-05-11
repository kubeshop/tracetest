package resources

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/global"
	global_formatters "github.com/kubeshop/tracetest/cli/global/formatters"
	resources_actions "github.com/kubeshop/tracetest/cli/resources/actions"
	resources_formatters "github.com/kubeshop/tracetest/cli/resources/formatters"
	resources_parameters "github.com/kubeshop/tracetest/cli/resources/parameters"
	"github.com/spf13/cobra"
)

type Apply struct {
	args[*resources_parameters.Apply]
}

func NewApply(root global.Root) Apply {
	parameters := resources_parameters.NewApply()
	defaults := NewDefaults("apply", root.Setup)

	apply := Apply{
		args: NewArgs(defaults, parameters),
	}

	apply.Cmd = &cobra.Command{
		GroupID: global.GroupResources.ID,
		Use:     "apply [resource type]",
		Short:   "Apply resources",
		Long:    "Apply (create/update) resources to your Tracetest server",
		Args:    cobra.MinimumNArgs(1),
		PreRun:  defaults.PreRun,
		Run:     defaults.Run(apply.Run),
		PostRun: defaults.PostRun,
	}

	apply.Cmd.Flags().StringVarP(&apply.Parameters.DefinitionFile, "file", "f", "", "file path with name where to export the resource")
	root.Cmd.AddCommand(apply.Cmd)

	return apply
}

func (a Apply) Run(cmd *cobra.Command, args []string) (string, error) {
	if a.Parameters.DefinitionFile == "" {
		return "", fmt.Errorf("file with definition must be specified")
	}

	ctx := context.Background()

	applyArgs := resources_actions.ApplyArgs{
		File: a.Parameters.DefinitionFile,
	}

	resource, _, err := a.Setup.ResourceActions.Apply(ctx, applyArgs)
	if err != nil {
		return "", err
	}

	resourceFormatter := a.Setup.ResourceActions.Formatter()
	formatter := resources_formatters.BuildFormatter(a.Setup.Parameters.Output, global_formatters.YAML, resourceFormatter)

	result, err := formatter.Format(resource)
	if err != nil {
		return "", err
	}

	return result, nil
}
