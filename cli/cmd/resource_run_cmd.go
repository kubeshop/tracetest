package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/cmdutil"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/runner"
	"github.com/spf13/cobra"

	cloudCmd "github.com/kubeshop/tracetest/cli/cloud/cmd"
)

var (
	runParams = &cmdutil.RunParameters{}
	runCmd    *cobra.Command
)

func init() {
	runCmd = &cobra.Command{
		GroupID: cmdGroupResources.ID,
		Use:     fmt.Sprintf("run [%s]", runnableResourceList()),
		Short:   "Run tests and test suites",
		Long:    "Run tests and test suites",
		PreRun:  setupCommand(WithOptionalResourceName()),
		Run: WithResourceMiddleware(func(ctx context.Context, _ *cobra.Command, args []string) (string, error) {
			if cliConfig.Jwt != "" {
				return cloudCmd.RunMultipleFiles(ctx, runParams, &cliConfig, runnerRegistry, output)
			}

			return runSingleFile(ctx)

		}, runParams),
		PostRun: teardownCommand,
	}

	runCmd.Flags().StringSliceVarP(&runParams.DefinitionFiles, "file", "f", []string{}, "path to the definition file (can be defined multiple times)")
	runCmd.Flags().StringVarP(&runParams.ID, "id", "", "", "id of the resource to run (can be defined multiple times)")
	runCmd.Flags().StringVarP(&runParams.VarsID, "vars", "", "", "variable set file or ID to be used")
	runCmd.Flags().BoolVarP(&runParams.SkipResultWait, "skip-result-wait", "W", false, "do not wait for results. exit immediately after test run started")
	runCmd.Flags().StringVarP(&runParams.JUnitOuptutFile, "junit", "j", "", "file path to save test results in junit format")
	runCmd.Flags().StringSliceVar(&runParams.RequiredGates, "required-gates", []string{}, "override default required gate. "+validRequiredGatesMsg())

	//deprecated
	runCmd.Flags().StringVarP(&runParams.EnvID, "environment", "e", "", "environment file or ID to be used")
	runCmd.Flags().MarkDeprecated("environment", "use --vars instead")
	runCmd.Flags().MarkShorthandDeprecated("e", "use --vars instead")

	rootCmd.AddCommand(runCmd)
}

func runSingleFile(ctx context.Context) (string, error) {
	orchestrator := runner.Orchestrator(
		cliLogger,
		config.GetAPIClient(cliConfig),
		variableSetClient,
		runnerRegistry,
	)

	if runParams.EnvID != "" {
		runParams.VarsID = runParams.EnvID
	}

	var definitionFile string
	if len(runParams.DefinitionFiles) > 0 {
		definitionFile = runParams.DefinitionFiles[0]
	}

	runParams := runner.RunOptions{
		ID:              runParams.ID,
		DefinitionFile:  definitionFile,
		VarsID:          runParams.VarsID,
		SkipResultWait:  runParams.SkipResultWait,
		JUnitOuptutFile: runParams.JUnitOuptutFile,
		RequiredGates:   runParams.RequiredGates,
	}

	exitCode, err := orchestrator.Run(ctx, runParams, output)
	if err != nil {
		return "", err
	}

	ExitCLI(exitCode)

	// ExitCLI will exit the process, so this return is just to satisfy the compiler
	return "", nil
}

func getResourceType(runParams *cmdutil.RunParameters, resourceParams *resourceParameters) (string, error) {
	if resourceParams.ResourceName != "" {
		return resourceParams.ResourceName, nil
	}

	filePath := ""
	if runParams.DefinitionFiles != nil && len(runParams.DefinitionFiles) > 0 {
		filePath = runParams.DefinitionFiles[0]
	}

	if filePath != "" {
		return cmdutil.GetResourceTypeFromFile(filePath)
	}

	return "", fmt.Errorf("resourceName is empty and no definition file provided")
}

func validRequiredGatesMsg() string {
	opts := make([]string, 0, len(openapi.AllowedSupportedGatesEnumValues))
	for _, v := range openapi.AllowedSupportedGatesEnumValues {
		opts = append(opts, string(v))
	}

	return "valid options: " + strings.Join(opts, ", ")
}
