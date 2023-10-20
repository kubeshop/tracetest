package pollingprofile

import (
	"testing"

	"github.com/kubeshop/tracetest/testing/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestDeletePollingProfile(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to delete the polling profile
	// Then it should return a error message, showing that we cannot delete a prolling profile
	result := tracetestcli.Exec(t, "delete pollingprofile --id current", tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 1)
	require.Contains(result.StdErr, "resource PollingProfile does not support the action")
}
