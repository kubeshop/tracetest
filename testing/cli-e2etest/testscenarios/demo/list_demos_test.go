package demo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func addListDemoPreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new environment
	// Then it should be applied with success
	newDemoPath := env.GetTestResourcePath(t, "new-demo")
	helpers.InjectIdIntoDemoFile(t, newDemoPath, "")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply demo --file %s", newDemoPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	// When I try to set up another demo
	// Then it should be applied with success
	anotherDemoPath := env.GetTestResourcePath(t, "another-demo")
	helpers.InjectIdIntoDemoFile(t, anotherDemoPath, "")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply demo --file %s", anotherDemoPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestListDemos(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("list no environments", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And there is no demos
		result := tracetestcli.Exec(t, "list demo --sortBy name --sortDirection asc --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		demos := helpers.UnmarshalYAMLSequence[types.DemoResource](t, result.StdOut)
		require.Len(demos, 0)
	})

	addListDemoPreReqs(t, env)

	t.Run("list with invalid sortBy field", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these demos by an invalid field
		// Then I should receive an error
		result := tracetestcli.Exec(t, "list demo --sortBy invalid-field --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 1)
		require.Contains(result.StdErr, "invalid sort field: invalid-field") // TODO: think on how to improve this error handling
	})

	t.Run("list with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these demos by a valid field and in YAML format
		// Then I should receive two demos
		result := tracetestcli.Exec(t, "list demo --sortBy name --sortDirection asc --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		demos := helpers.UnmarshalYAMLSequence[types.DemoResource](t, result.StdOut)
		require.Len(demos, 2)

		anotherDemo := demos[0]
		require.Equal("Demo", anotherDemo.Type)
		require.Equal("another-dev", anotherDemo.Spec.Name)
		require.Equal("pokeshop", anotherDemo.Spec.Type)
		require.True(anotherDemo.Spec.Enabled)
		require.Equal("new-dev-grpc:9091", anotherDemo.Spec.Pokeshop.GrpcEndpoint)
		require.Equal("http://new-dev-endpoint:1234", anotherDemo.Spec.Pokeshop.HttpEndpoint)

		demo := demos[1]
		require.Equal("Demo", demo.Type)
		require.Equal("dev", demo.Spec.Name)
		require.Equal("otelstore", demo.Spec.Type)
		require.True(demo.Spec.Enabled)
		require.Equal("http://dev-cart:8082", demo.Spec.OTelStore.CartEndpoint)
		require.Equal("http://dev-checkout:8083", demo.Spec.OTelStore.CheckoutEndpoint)
		require.Equal("http://dev-frontend:9000", demo.Spec.OTelStore.FrontendEndpoint)
		require.Equal("http://dev-product:8081", demo.Spec.OTelStore.ProductCatalogEndpoint)
	})

	t.Run("list with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these demos by a valid field and in JSON format
		// Then I should receive two demos
		result := tracetestcli.Exec(t, "list demo --sortBy name --sortDirection asc --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		demos := helpers.UnmarshalJSON[[]types.DemoResource](t, result.StdOut)
		require.Len(demos, 2)

		anotherDemo := demos[0]
		require.Equal("Demo", anotherDemo.Type)
		require.Equal("another-dev", anotherDemo.Spec.Name)
		require.Equal("pokeshop", anotherDemo.Spec.Type)
		require.True(anotherDemo.Spec.Enabled)
		require.Equal("new-dev-grpc:9091", anotherDemo.Spec.Pokeshop.GrpcEndpoint)
		require.Equal("http://new-dev-endpoint:1234", anotherDemo.Spec.Pokeshop.HttpEndpoint)

		demo := demos[1]
		require.Equal("Demo", demo.Type)
		require.Equal("dev", demo.Spec.Name)
		require.Equal("otelstore", demo.Spec.Type)
		require.True(demo.Spec.Enabled)
		require.Equal("http://dev-cart:8082", demo.Spec.OTelStore.CartEndpoint)
		require.Equal("http://dev-checkout:8083", demo.Spec.OTelStore.CheckoutEndpoint)
		require.Equal("http://dev-frontend:9000", demo.Spec.OTelStore.FrontendEndpoint)
		require.Equal("http://dev-product:8081", demo.Spec.OTelStore.ProductCatalogEndpoint)
	})

	t.Run("list with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these demos by a valid field and in pretty format
		// Then it should print a table with 5 lines printed: header, separator, two demos and empty line
		result := tracetestcli.Exec(t, "list demo --sortBy name --sortDirection asc --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		lines := strings.Split(result.StdOut, "\n")
		require.Len(lines, 5)

		require.Contains(lines[2], "another-dev")
		require.Contains(lines[3], "dev")
	})

	t.Run("list with YAML format skipping the first and taking one item", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these demos by a valid field, paging options and in YAML format
		// Then I should receive one demo
		result := tracetestcli.Exec(t, "list demo --sortBy name --sortDirection asc --skip 1 --take 1 --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		demos := helpers.UnmarshalYAMLSequence[types.DemoResource](t, result.StdOut)
		require.Len(demos, 1)

		demo := demos[0]
		require.Equal("Demo", demo.Type)
		require.Equal("dev", demo.Spec.Name)
		require.Equal("otelstore", demo.Spec.Type)
		require.True(demo.Spec.Enabled)
		require.Equal("http://dev-cart:8082", demo.Spec.OTelStore.CartEndpoint)
		require.Equal("http://dev-checkout:8083", demo.Spec.OTelStore.CheckoutEndpoint)
		require.Equal("http://dev-frontend:9000", demo.Spec.OTelStore.FrontendEndpoint)
		require.Equal("http://dev-product:8081", demo.Spec.OTelStore.ProductCatalogEndpoint)
	})
}
