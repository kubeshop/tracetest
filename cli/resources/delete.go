package resources

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/global"
	resources_parameters "github.com/kubeshop/tracetest/cli/resources/parameters"
	"github.com/spf13/cobra"
)

type Delete struct {
	args[*resources_parameters.Id]
}

func NewDelete(root global.Root) Delete {
	parameters := resources_parameters.NewId()
	defaults := NewDefaults("delete", root.Setup)
	delete := Delete{
		args: NewArgs(defaults, parameters),
	}

	delete.Cmd = &cobra.Command{
		GroupID: global.GroupResources.ID,
		Use:     "delete [resource type]",
		Long:    "Delete resources from your Tracetest server",
		Short:   "Delete resources",
		Args:    cobra.MinimumNArgs(1),
		PreRun:  defaults.PreRun,
		Run:     defaults.Run(delete.Run),
		PostRun: defaults.PostRun,
	}

	delete.Cmd.Flags().StringVar(&delete.Parameters.ResourceID, "id", "", "id of the resource to delete")
	root.Cmd.AddCommand(delete.Cmd)

	return delete
}

func (d Delete) Run(cmd *cobra.Command, args []string) (string, error) {
	if d.Parameters.ResourceID == "" {
		return "", fmt.Errorf("id of the resource to delete must be specified")
	}

	ctx := context.Background()

	message, err := d.Setup.ResourceActions.Delete(ctx, d.Parameters.ResourceID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("✔ %s", message), nil
}
