package cmd

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/cloud/runner"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/formatters"
	cliRunner "github.com/kubeshop/tracetest/cli/runner"
	"go.uber.org/zap"
)

func Wait(ctx context.Context, logger *zap.Logger, cliConfig *config.Config, runGroupID, format string) (int, error) {
	rungroupWaiter := runner.RunGroup(logger, config.GetAPIClient(*cliConfig))
	runGroup, err := rungroupWaiter.WaitForCompletion(ctx, runGroupID)
	if err != nil {
		return runner.ExitCodeGeneralError, err
	}

	formatter := formatters.MultipleRun[cliRunner.RunResult](func() string { return cliConfig.UI() }, true)
	runnerGetter := func(resource any) (formatters.Runner[cliRunner.RunResult], error) {
		return nil, nil
	}

	output := formatters.MultipleRunOutput[cliRunner.RunResult]{
		Runs:         []cliRunner.RunResult{},
		Resources:    []any{},
		RunGroup:     runGroup,
		RunnerGetter: runnerGetter,
		HasResults:   true,
	}

	fmt.Println(formatter.Format(output, formatters.Output(format)))

	exitCode := runner.ExitCodeSuccess
	if runGroup.GetStatus() == "failed" {
		exitCode = runner.ExitCodeTestNotPassed
	}

	return exitCode, nil
}
