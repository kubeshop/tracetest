package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

var resourceID string

var getCmd = &cobra.Command{
	GroupID:   cmdGroupResources.ID,
	Use:       fmt.Sprintf("get %s", strings.Join(validArgs, "|")),
	Short:     "Get resource",
	Long:      "Get a resource from your Tracetest server",
	PreRun:    setupCommand(),
	Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	ValidArgs: validArgs,
	Run: WithResultHandler(func(cmd *cobra.Command, args []string) (string, error) {
		if resourceID == "" {
			return "", fmt.Errorf("id of the resource to get must be specified")
		}

		resourceType := args[0]
		ctx := context.Background()

		analytics.Track("Resource Get", "cmd", map[string]string{
			resourceType: resourceType,
		})

		resourceActions, err := resourceRegistry.Get(resourceType)
		if err != nil {
			return "", err
		}

		if output == string(formatters.JSON) {
			ctx = context.WithValue(ctx, "X-Tracetest-Augmented", true)
		}

		resource, err := resourceActions.Get(ctx, resourceID)
		if err != nil && errors.Is(err, utils.ResourceNotFound) {
			return fmt.Sprintf("Resource %s with ID %s not found", resourceType, resourceID), nil
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
	}),
	PostRun: teardownCommand,
}

func init() {
	getCmd.Flags().StringVar(&resourceID, "id", "", "id of the resource to get")
	rootCmd.AddCommand(getCmd)
}
