package resources

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/global"
	resource_parameters "github.com/kubeshop/tracetest/cli/resources/parameters"
	"github.com/spf13/cobra"
)

type Export struct {
	args[*resource_parameters.Export]
}

func NewExport(root global.Root) Export {
	parameters := resource_parameters.NewExport()
	defaults := NewDefaults("export", root.Setup)

	export := Export{
		args: NewArgs(defaults, parameters),
	}

	export.Cmd = &cobra.Command{
		GroupID: global.GroupResources.ID,
		Use:     "export [resource type]",
		Long:    "Export a resource from your Tracetest server",
		Short:   "Export resource",
		Args:    cobra.MinimumNArgs(1),
		PreRun:  defaults.PreRun,
		Run:     defaults.Run(export.Run),
		PostRun: defaults.PostRun,
	}

	export.Cmd.Flags().StringVar(&export.Parameters.ResourceID, "id", "", "id of the resource to export")
	export.Cmd.Flags().StringVarP(&export.Parameters.ResourceFile, "file", "f", "resource.yaml", "file path with name where to export the resource")
	root.Cmd.AddCommand(export.Cmd)

	return export
}

var export Export

func GetExport() Export {
	if export.Cmd == nil {
		panic("export command not initialized")
	}

	return export
}

func (e Export) Run(cmd *cobra.Command, args []string) (string, error) {
	if e.Parameters.ResourceID == "" {
		return "", fmt.Errorf("id of the resource to export must be specified")
	}

	ctx := context.Background()
	err := e.Setup.ResourceActions.Export(ctx, e.Parameters.ResourceID, e.Parameters.ResourceFile)
	if err != nil {
		return "", err
	}

	return "✔  Definition exported successfully", nil
}
