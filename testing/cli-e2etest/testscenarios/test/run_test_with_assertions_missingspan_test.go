package test

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestRunWithAssertionMissingSpan(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try run a test that has a selector that will find a span
	// And  a selector that will find no spans
	// Then it run, fail and return only the assertion to the selector with spans
	testWithAssertionWithNoSpan := env.GetTestResourcePath(t, "test-with-assertion-with-missing-span")

	command := fmt.Sprintf("test run -w -d %s", testWithAssertionWithNoSpan)
	result := tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 1) // test failed
	require.Contains(result.StdOut, "It should call /api/tests on HTTP")
	require.NotContains(result.StdOut, "It should not have grpc spans")
}
