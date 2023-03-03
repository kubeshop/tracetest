package config_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestDemoConfig(t *testing.T) {
	t.Run("DefaultValues", func(t *testing.T) {
		cfg, err := config.New(nil)
		require.NoError(t, err)

		defaultEndponts := map[string]string{
			"PokeshopHttp":       "",
			"PokeshopGrpc":       "",
			"OtelFrontend":       "",
			"OtelProductCatalog": "",
			"OtelCart":           "",
			"OtelCheckout":       "",
		}

		assert.DeepEqual(t, []string(nil), cfg.DemoEnabled())
		assert.DeepEqual(t, defaultEndponts, cfg.DemoEndpoints())
	})

	t.Run("File", func(t *testing.T) {
		t.Parallel()

		cfg := configFromFile(t, "./testdata/demo.yaml")

		expectedEndpoints := map[string]string{
			"PokeshopHttp":       "http://demo-pokemon-api.demo",
			"PokeshopGrpc":       "demo-pokemon-api.demo:8082",
			"OtelFrontend":       "http://otel-frontend.otel-demo:8084",
			"OtelProductCatalog": "http://otel-productcatalogservice.otel-demo:3550",
			"OtelCart":           "http://otel-cartservice.otel-demo:7070",
			"OtelCheckout":       "http://otel-checkoutservice.otel-demo:5050",
		}

		assert.DeepEqual(t, []string{"pokeshop", "otel"}, cfg.DemoEnabled())
		assert.DeepEqual(t, expectedEndpoints, cfg.DemoEndpoints())

	})
}
