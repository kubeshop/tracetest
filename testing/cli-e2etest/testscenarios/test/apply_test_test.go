package test

import (
	"fmt"
	"testing"

	"atomicgo.dev/assert"
	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestApplyTest(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new test
	// Then it should be applied with success
	testPath := env.GetTestResourcePath(t, "list")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply test --file %s", testPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	// When I try to get a test
	// Then it should return the test applied on the last step
	result = tracetestcli.Exec(t, "get test --id fH_8AulVR", tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	listTest := helpers.UnmarshalYAML[types.TestResource](t, result.StdOut)
	assert.Equal("Test", listTest.Type)
	assert.Equal("fH_8AulVR", listTest.Spec.ID)
	assert.Equal("Pokeshop - List", listTest.Spec.Name)
	assert.Equal("List Pokemon", listTest.Spec.Description)
	assert.Equal("http", listTest.Spec.Trigger.Type)
	assert.Equal("http://demo-api:8081/pokemon?take=20&skip=0", listTest.Spec.Trigger.HTTPRequest.URL)
	assert.Equal("GET", listTest.Spec.Trigger.HTTPRequest.Method)
	assert.Equal("", listTest.Spec.Trigger.HTTPRequest.Body)
	require.Len(listTest.Spec.Trigger.HTTPRequest.Headers, 1)
	assert.Equal("Content-Type", listTest.Spec.Trigger.HTTPRequest.Headers[0].Key)
	assert.Equal("application/json", listTest.Spec.Trigger.HTTPRequest.Headers[0].Value)
}
