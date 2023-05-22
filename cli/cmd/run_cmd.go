package cmd

import (
	"context"
	"fmt"

	cienvironment "github.com/cucumber/ci-environment/go"
	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/spf13/cobra"
)

var environmentID string

var runCmd = &cobra.Command{
	GroupID: cmdGroupResources.ID,
	Use:     "run [resource type]",
	Long:    "Run a resource in your Tracetest server",
	Short:   "Run resource",
	PreRun:  setupCommand(),
	Args:    cobra.MinimumNArgs(1),
	Run: WithResultHandler(func(cmd *cobra.Command, args []string) (string, error) {
		if resourceID == "" && definitionFile == "" {
			return "", fmt.Errorf("id or file must be specified")
		}

		resourceType := args[0]
		ctx := context.Background()

		analytics.Track("Resource Run", "cmd", map[string]string{
			resourceType: resourceType,
		})

		resourceActions, err := resourceRegistry.Get(resourceType)
		if err != nil {
			return "", err
		}

		runArguments := actions.RunArguments{
			EnvironmentID: environmentID,
			Metadata:      getRunMetadata(),
			Variables:     make(map[string]string),
		}

		var runResult any

		if resourceID != "" {
			runResult, err = resourceActions.RunByID(ctx, resourceID, runArguments)
			if err != nil {
				return "", fmt.Errorf("cannot run %s: %w", resourceType, err)
			}
		}

		if definitionFile != "" {
			file, err := file.Read(definitionFile)
			if err != nil {
				return "", fmt.Errorf("could not read file: %w", err)
			}

			runResult, err = resourceActions.Run(ctx, file, runArguments)
			if err != nil {
				return "", fmt.Errorf("cannot run %s: %w", resourceType, err)
			}
		}

		// resourceFormatter := resourceActions.Formatter()
		// formatter := formatters.BuildFormatter(output, formatters.Pretty, resourceFormatter)

		// result, err := formatter.FormatRunResult(runResult)
		// if err != nil {
		// 	return "", fmt.Errorf("cannot format output: %w", err)
		// }

		// return result, nil
		return fmt.Sprintf("%v", runResult), nil
	}),
	PostRun: teardownCommand,
}

func getRunMetadata() map[string]string {
	ci := cienvironment.DetectCIEnvironment()
	if ci == nil {
		return map[string]string{}
	}

	metadata := map[string]string{
		"name":        ci.Name,
		"url":         ci.URL,
		"buildNumber": ci.BuildNumber,
	}

	if ci.Git != nil {
		metadata["branch"] = ci.Git.Branch
		metadata["tag"] = ci.Git.Tag
		metadata["revision"] = ci.Git.Revision
	}

	return metadata
}

func init() {
	runCmd.Flags().StringVar(&resourceID, "id", "", "id of the resource to get")
	runCmd.Flags().StringVarP(&definitionFile, "file", "f", "", "path to the definition file containing the transaction")
	runCmd.Flags().StringVarP(&environmentID, "environment", "e", "", "id of the environment to be used in the run")
	rootCmd.AddCommand(runCmd)
}
