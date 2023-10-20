package testrunner

import (
	"testing"

	"github.com/kubeshop/tracetest/testing/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestDeleteTestRunner(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to delete the testrunner
	// Then it should return a error message, showing that we cannot delete a testrunner
	result := tracetestcli.Exec(t, "delete testrunner --id current", tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 1)
	require.Contains(result.StdErr, "resource TestRunner does not support the action")
}
