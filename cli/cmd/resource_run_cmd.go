package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/cmdutil"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/runner"
	"github.com/spf13/cobra"

	cloudCmd "github.com/kubeshop/tracetest/cli/cloud/cmd"
	cliRunner "github.com/kubeshop/tracetest/cli/runner"
)

var (
	runParams = &cmdutil.RunParameters{}
	runCmd    *cobra.Command
)

const ExitCodeResourceNotFound = 3

func init() {
	runCmd = &cobra.Command{
		GroupID: cmdGroupResources.ID,
		Use:     fmt.Sprintf("run [%s]", runnableResourceList()),
		Short:   "Run tests and test suites",
		Long:    "Run tests and test suites",
		PreRun:  setupCommand(WithOptionalResourceName()),
		Run: WithResourceMiddleware(func(ctx context.Context, _ *cobra.Command, args []string) (string, error) {
			runParams.ResourceName = resourceParams.ResourceName
			if cliConfig.Jwt != "" {
				return runMultipleFiles(ctx)
			}

			return runSingleFile(ctx)

		}, runParams),
		PostRun: teardownCommand,
	}

	runCmd.Flags().StringSliceVarP(&runParams.DefinitionFiles, "file", "f", []string{}, "path to the definition file (you can define multiple paths by repeating this option), or alternatively, you can pass a directory.")
	runCmd.Flags().StringSliceVarP(&runParams.IDs, "id", "", []string{}, "id of the resource to run (can be defined multiple times)")
	runCmd.Flags().StringVarP(&runParams.VarsID, "vars", "", "", "variable set file or ID to be used")
	runCmd.Flags().StringVarP(&runParams.RunGroupID, "group", "g", "", "Sets the Run Group ID for the run. This is used to group multiple runs together.")
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

		if !hasValidResourceFiles() {
			cliLogger.Warn("Invalid definition file found, stopping execution")
			ExitCLI(ExitCodeResourceNotFound)
			return "", nil
		}
	}

	ID := ""
	if len(runParams.IDs) > 0 {
		ID = runParams.IDs[0]
	}

	runParams := runner.RunOptions{
		ID:              ID,
		DefinitionFile:  definitionFile,
		VarsID:          runParams.VarsID,
		SkipResultWait:  runParams.SkipResultWait,
		JUnitOuptutFile: runParams.JUnitOuptutFile,
		RequiredGates:   runParams.RequiredGates,
	}

	exitCode, err := orchestrator.Run(ctx, runParams, output)
	ExitCLI(exitCode)

	// ExitCLI will exit the process, so this return is just to satisfy the compiler
	return "", err
}

func runMultipleFiles(ctx context.Context) (string, error) {
	if !hasValidResourceFiles() {
		cliLogger.Warn("Invalid definition files found, stopping execution")
		ExitCLI(ExitCodeResourceNotFound)
		return "", nil
	}

	exitCode, err := cloudCmd.RunMultipleFiles(ctx, cliLogger, httpClient, runParams, &cliConfig, runnerRegistry, output)
	// General Error is 1, which is the default for errors. if this is the case,
	// let the CLI handle the error and exit.
	// otherwise exit with the exit code.
	if exitCode > cliRunner.ExitCodeGeneralError {
		ExitCLI(exitCode)
	}
	return "", err
}

func validRequiredGatesMsg() string {
	opts := make([]string, 0, len(openapi.AllowedSupportedGatesEnumValues))
	for _, v := range openapi.AllowedSupportedGatesEnumValues {
		opts = append(opts, string(v))
	}

	return "valid options: " + strings.Join(opts, ", ")
}

func hasValidResourceFiles() bool {
	if len(runParams.DefinitionFiles) == 0 {
		return true
	}

	validResourceFiles := true

	for _, file := range runParams.DefinitionFiles {
		if !fileutil.IsExistingFile(file) {
			cliLogger.Warn("Definition file does not exist: " + file)
			validResourceFiles = false
		}
	}

	return validResourceFiles
}
