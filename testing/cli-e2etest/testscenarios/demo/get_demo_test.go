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

func addGetDemoPreReqs(t *testing.T, env environment.Manager) string {
	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new demo
	// Then it should be applied with success
	newDemoPath := env.GetTestResourcePath(t, "new-demo")
	helpers.InjectIdIntoDemoFile(t, newDemoPath, "")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply demo --file %s", newDemoPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	demo := helpers.UnmarshalYAML[types.DemoResource](t, result.StdOut)
	return demo.Spec.Id
}

func TestGetDemo(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("get with no demo initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And no demo registered

		// When I try to get a demo on yaml mode
		// Then it should return a error message
		result := tracetestcli.Exec(t, "get demo --id some-demo --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "Resource demo with ID some-demo not found")
	})

	registeredDemoId := addGetDemoPreReqs(t, env)

	t.Run("get with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a demo already set

		// When I try to get a demo on yaml mode
		// Then it should print a YAML
		command := fmt.Sprintf("get demo --id %s --output yaml", registeredDemoId)
		result := tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		demo := helpers.UnmarshalYAML[types.DemoResource](t, result.StdOut)

		require.Equal("Demo", demo.Type)
		require.Equal("dev", demo.Spec.Name)
		require.Equal("otelstore", demo.Spec.Type)
		require.True(demo.Spec.Enabled)
		require.Equal("http://dev-cart:8082", demo.Spec.OTelStore.CartEndpoint)
		require.Equal("http://dev-checkout:8083", demo.Spec.OTelStore.CheckoutEndpoint)
		require.Equal("http://dev-frontend:9000", demo.Spec.OTelStore.FrontendEndpoint)
		require.Equal("http://dev-product:8081", demo.Spec.OTelStore.ProductCatalogEndpoint)
	})

	t.Run("get with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a demo already set

		// When I try to get a demo on json mode
		// Then it should print a json
		command := fmt.Sprintf("get demo --id %s --output json", registeredDemoId)
		result := tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		demo := helpers.UnmarshalJSON[types.DemoResource](t, result.StdOut)

		require.Equal("Demo", demo.Type)
		require.Equal("dev", demo.Spec.Name)
		require.Equal("otelstore", demo.Spec.Type)
		require.True(demo.Spec.Enabled)
		require.Equal("http://dev-cart:8082", demo.Spec.OTelStore.CartEndpoint)
		require.Equal("http://dev-checkout:8083", demo.Spec.OTelStore.CheckoutEndpoint)
		require.Equal("http://dev-frontend:9000", demo.Spec.OTelStore.FrontendEndpoint)
		require.Equal("http://dev-product:8081", demo.Spec.OTelStore.ProductCatalogEndpoint)
	})

	t.Run("get with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a demo already set

		// When I try to get a demo on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, demo item and empty line
		command := fmt.Sprintf("get demo --id %s --output pretty", registeredDemoId)
		result := tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		parsedTable := helpers.UnmarshalTable(t, result.StdOut)
		require.Len(parsedTable, 1)

		singleLine := parsedTable[0]

		require.NotEmpty(singleLine["ID"]) // demo resource generates a random ID each time
		require.Equal("dev", singleLine["NAME"])
		require.Equal("otelstore", singleLine["TYPE"])
		require.Equal("true", singleLine["ENABLED"])
	})
}
