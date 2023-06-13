package demo

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestDeleteDemo(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to delete an demo that don't exist
	// Then it should return an error and say that this resource does not exist
	result := tracetestcli.Exec(t, "delete demo --id some-demo", tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 1)
	require.Contains(result.StdErr, "Resource demos with ID some-demo not found") // TODO: update this message to singular

	// When I try to set up a new demo
	// Then it should be applied with success
	newDemoPath := env.GetTestResourcePath(t, "new-demo")
	helpers.InjectIdIntoDemoFile(t, newDemoPath, "")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply demo --file %s", newDemoPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	demo := helpers.UnmarshalYAML[types.DemoResource](t, result.StdOut)

	// When I try to delete the demo
	// Then it should delete with success
	command := fmt.Sprintf("delete demo --id %s", demo.Spec.Id)
	result = tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
	require.Contains(result.StdOut, "✔ Demo successfully deleted")

	// When I try to get an demo again
	// Then it should return a message saying that the environment was not found
	command = fmt.Sprintf("get demo --id %s --output yaml", demo.Spec.Id)
	result = tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
	require.Contains(result.StdOut, fmt.Sprintf("Resource demo with ID %s not found", demo.Spec.Id))
}
