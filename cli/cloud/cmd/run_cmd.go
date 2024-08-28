package cmd

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/cloud/runner"
	"github.com/kubeshop/tracetest/cli/cmdutil"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/kubeshop/tracetest/cli/processor"
	"go.uber.org/zap"

	cliRunner "github.com/kubeshop/tracetest/cli/runner"
)

func RunMultipleFiles(ctx context.Context, logger *zap.Logger, httpClient *resourcemanager.HTTPClient, runParams *cmdutil.RunParameters, cliConfig *config.Config, runnerRegistry cliRunner.Registry, format string) (int, error) {
	if cliConfig.Jwt == "" {
		return cliRunner.ExitCodeGeneralError, fmt.Errorf("you should be authenticated to run multiple files, please run 'tracetest configure'")
	}

	variableSetPreprocessor := processor.VariableSet(logger)

	runGroup := runner.RunGroup(logger, config.GetAPIClient(*cliConfig))
	formatter := formatters.MultipleRun[cliRunner.RunResult](func() string { return cliConfig.UI() }, true)

	orchestrator := runner.MultiFileOrchestrator(
		logger,
		runGroup,
		GetVariableSetClient(logger, httpClient, variableSetPreprocessor),
		runnerRegistry,
		formatter,
	)

	return orchestrator.Run(ctx, runner.RunOptions{
		IDs:             runParams.IDs,
		ResourceName:    runParams.ResourceName,
		DefinitionFiles: runParams.DefinitionFiles,
		VarsID:          runParams.VarsID,
		SkipResultWait:  runParams.SkipResultWait,
		JUnitOuptutFile: runParams.JUnitOuptutFile,
		RequiredGates:   runParams.RequiredGates,
		RunGroupID:      runParams.RunGroupID,
	}, format)
}

func GetVariableSetClient(logger *zap.Logger, httpClient *resourcemanager.HTTPClient, preprocessor processor.Preprocessor) resourcemanager.Client {
	variableSetClient := resourcemanager.NewClient(
		httpClient, logger,
		"variableset", "variablesets",
		resourcemanager.WithTableConfig(resourcemanager.TableConfig{
			Cells: []resourcemanager.TableCellConfig{
				{Header: "ID", Path: "spec.id"},
				{Header: "NAME", Path: "spec.name"},
				{Header: "DESCRIPTION", Path: "spec.description"},
			},
		}),
		resourcemanager.WithResourceType("VariableSet"),
		resourcemanager.WithApplyPreProcessor(preprocessor.Preprocess),
		resourcemanager.WithDeprecatedAlias("Environment"),
	)

	return variableSetClient
}
