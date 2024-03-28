package cmd

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/cloud/runner"
	"github.com/kubeshop/tracetest/cli/cmdutil"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/preprocessor"

	cliRunner "github.com/kubeshop/tracetest/cli/runner"
)

func RunMultipleFiles(ctx context.Context, runParams *cmdutil.RunParameters, cliConfig *config.Config, runnerRegistry cliRunner.Registry) (string, error) {
	if cliConfig.Jwt == "" {
		return "", fmt.Errorf("you should be authenticated to run multiple files, please run 'tracetest configure'")
	}

	variableSetPreprocessor := preprocessor.VariableSet(cmdutil.GetLogger())

	orchestrator := runner.MultiFileOrchestrator(
		cmdutil.GetLogger(),
		config.GetAPIClient(*cliConfig),
		cmdutil.GetVariableSetClient(variableSetPreprocessor),
		runnerRegistry,
	)

	orchestrator.Run(ctx, runner.RunOptions{
		IDs:             runParams.ID,
		DefinitionFiles: runParams.DefinitionFiles,
		VarsID:          runParams.VarsID,
		SkipResultWait:  runParams.SkipResultWait,
		JUnitOuptutFile: runParams.JUnitOuptutFile,
		RequiredGates:   runParams.RequiredGates,
	})

	// create test / test suites if not defined

	// create run group

	// trigger each test run

	// return single formatting if skip is enabled

	// loop, validate if runGroup is finished

	// format run results

	return "", nil
}
