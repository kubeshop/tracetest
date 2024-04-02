package cmd

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/cloud/runner"
	"github.com/kubeshop/tracetest/cli/cmdutil"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/preprocessor"

	cliRunner "github.com/kubeshop/tracetest/cli/runner"
)

func RunMultipleFiles(ctx context.Context, runParams *cmdutil.RunParameters, cliConfig *config.Config, runnerRegistry cliRunner.Registry, format string) (string, error) {
	if cliConfig.Jwt == "" {
		return "", fmt.Errorf("you should be authenticated to run multiple files, please run 'tracetest configure'")
	}

	variableSetPreprocessor := preprocessor.VariableSet(cmdutil.GetLogger())

	formatter := formatters.MultipleRun(func() string { return cliConfig.UI() }, true)

	orchestrator := runner.MultiFileOrchestrator(
		cmdutil.GetLogger(),
		config.GetAPIClient(*cliConfig),
		cmdutil.GetVariableSetClient(variableSetPreprocessor),
		runnerRegistry,
		formatter,
	)

	orchestrator.Run(ctx, runner.RunOptions{
		DefinitionFiles: runParams.DefinitionFiles,
		VarsID:          runParams.VarsID,
		SkipResultWait:  runParams.SkipResultWait,
		JUnitOuptutFile: runParams.JUnitOuptutFile,
		RequiredGates:   runParams.RequiredGates,
	}, format)

	return "", nil
}
