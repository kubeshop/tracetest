package cmd

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/cmdutil"
	"github.com/kubeshop/tracetest/cli/config"
)

func RunMultipleFiles(ctx context.Context, runParams *cmdutil.RunParameters, cliConfig *config.Config) (string, error) {
	if cliConfig.Jwt == "" {
		return "", fmt.Errorf("you should be authenticated to run multiple files, please run 'tracetest configure'")
	}

	// create test / test suites if not defined

	// create run group

	// trigger each test run

	// return single formatting if skip is enabled

	// loop, validate if runGroup is finished

	// format run results

	return "", nil
}
