package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/parameters"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

var getParams = &parameters.ResourceIdParams{}

var getCmd = &cobra.Command{
	GroupID: cmdGroupResources.ID,
	Use:     fmt.Sprintf("get %s", strings.Join(parameters.ValidResources, "|")),
	Short:   "Get resource",
	Long:    "Get a resource from your Tracetest server",
	PreRun:  setupCommand(),
	Run: WithResourceMiddleware(func(_ *cobra.Command, args []string) (string, error) {
		resourceType := args[0]
		ctx := context.Background()

		resourceActions, err := resourceRegistry.Get(resourceType)
		if err != nil {
			return "", err
		}

		if output == string(formatters.JSON) {
			ctx = context.WithValue(ctx, "X-Tracetest-Augmented", true)
		}

		resource, err := resourceActions.Get(ctx, getParams.ResourceId)
		if err != nil && errors.Is(err, utils.ErrResourceNotFound) {
			return fmt.Sprintf("Resource %s with ID %s not found", resourceType, getParams.ResourceId), nil
		} else if err != nil {
			return "", err
		}

		resourceFormatter := resourceActions.Formatter()
		formatter := formatters.BuildFormatter(output, formatters.YAML, resourceFormatter)

		result, err := formatter.Format(resource)
		if err != nil {
			return "", err
		}

		return result, nil
	}, getParams),
	PostRun: teardownCommand,
}

func init() {
	getCmd.Flags().StringVar(&getParams.ResourceId, "id", "", "id of the resource to get")
	rootCmd.AddCommand(getCmd)
}
