package environment

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestListEnvironments(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new environment
	// Then it should be applied with success
	newEnvironmentPath := env.GetTestResourcePath(t, "new-environment")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply environment --file %s", newEnvironmentPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	environmentVars := helpers.UnmarshalYAML[types.EnvironmentResource](t, result.StdOut)
	require.Equal("Environment", environmentVars.Type)
	require.Equal(".env", environmentVars.Spec.ID)
	require.Equal(".env", environmentVars.Spec.Name)

	// When I try to set up a another environment
	// Then it should be applied with success
	anotherEnvironmentPath := env.GetTestResourcePath(t, "another-environment")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply environment --file %s", anotherEnvironmentPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	environmentVars = helpers.UnmarshalYAML[types.EnvironmentResource](t, result.StdOut)
	require.Equal("Environment", environmentVars.Type)
	require.Equal("another-env", environmentVars.Spec.ID)
	require.Equal("another-env", environmentVars.Spec.Name)

	// When I try to set up a third environment
	// Then it should be applied with success
	oneMoreEnvironmentPath := env.GetTestResourcePath(t, "one-more-environment")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply environment --file %s", oneMoreEnvironmentPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	environmentVars = helpers.UnmarshalYAML[types.EnvironmentResource](t, result.StdOut)
	require.Equal("Environment", environmentVars.Type)
	require.Equal("one-more-env", environmentVars.Spec.ID)
	require.Equal("one-more-env", environmentVars.Spec.Name)

	// When I try to list these environments by an invalid field
	// Then I should receive an error

	result = tracetestcli.Exec(t, "list environment --sortBy id --output yaml", tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 1)
	require.Contains(result.StdErr, "invalid sort field: id") // TODO: think on how to improve this error handling

	// When I try to list these environments by a valid field and in YAML format
	// Then I should receive three environments

	result = tracetestcli.Exec(t, "list environment --sortBy name --sortDirection asc --output yaml", tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	environmentVarsList := helpers.UnmarshalYAMLSequence[types.EnvironmentResource](t, result.StdOut)
	require.Len(environmentVarsList, 3)

	anotherEnvironmentVars := environmentVarsList[0]
	require.Equal("Environment", anotherEnvironmentVars.Type)
	require.Equal("another-env", anotherEnvironmentVars.Spec.ID)
	require.Equal("another-env", anotherEnvironmentVars.Spec.Name)
	require.Len(anotherEnvironmentVars.Spec.Values, 2)
	require.Equal("Here", anotherEnvironmentVars.Spec.Values[0].Key)
	require.Equal("We", anotherEnvironmentVars.Spec.Values[0].Value)
	require.Equal("Come", anotherEnvironmentVars.Spec.Values[1].Key)
	require.Equal("Again", anotherEnvironmentVars.Spec.Values[1].Value)

	environmentVars = environmentVarsList[1]
	require.Equal("Environment", environmentVars.Type)
	require.Equal(".env", environmentVars.Spec.ID)
	require.Equal(".env", environmentVars.Spec.Name)
	require.Len(environmentVars.Spec.Values, 2)
	require.Equal("FIRST_VAR", environmentVars.Spec.Values[0].Key)
	require.Equal("some-value", environmentVars.Spec.Values[0].Value)
	require.Equal("SECOND_VAR", environmentVars.Spec.Values[1].Key)
	require.Equal("another_value", environmentVars.Spec.Values[1].Value)

	oneMoreEnvironmentVars := environmentVarsList[2]
	require.Equal("Environment", oneMoreEnvironmentVars.Type)
	require.Equal("one-more-env", oneMoreEnvironmentVars.Spec.ID)
	require.Equal("one-more-env", oneMoreEnvironmentVars.Spec.Name)
	require.Len(oneMoreEnvironmentVars.Spec.Values, 2)
	require.Equal("This", oneMoreEnvironmentVars.Spec.Values[0].Key)
	require.Equal("Is", oneMoreEnvironmentVars.Spec.Values[0].Value)
	require.Equal("The", oneMoreEnvironmentVars.Spec.Values[1].Key)
	require.Equal("Third", oneMoreEnvironmentVars.Spec.Values[1].Value)
}
