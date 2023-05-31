package test

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestRunTestWithGrpcTrigger(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t, environment.WithDataStoreEnabled(), environment.WithPokeshop())
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created
	// And the datasource is already set

	// When I try to run a test with a gRPC trigger with embedded protobuf
	// Then it should pass
	testFile := env.GetTestResourcePath(t, "grpc-trigger-embedded-protobuf")

	command := fmt.Sprintf("test run -w -d %s", testFile)
	result := tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
	require.Contains(result.StdOut, "It calls Pokeshop correctly")

	// When I try to run a test with a gRPC trigger with a reference to a protobuf file
	// Then it should pass
	testFile = env.GetTestResourcePath(t, "grpc-trigger-reference-protobuf")

	command = fmt.Sprintf("test run -w -d %s", testFile)
	result = tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
	require.Contains(result.StdOut, "It calls Pokeshop correctly")
}
