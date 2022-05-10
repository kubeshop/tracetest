package executor_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kubeshop/tracetest/openapi"
	"github.com/kubeshop/tracetest/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutorIntegration(t *testing.T) {
	demoApp, err := test.GetDemoApplicationInstance()
	require.NoError(t, err)
	defer demoApp.Stop()

	resp, err := http.Get(fmt.Sprintf("http://%s/pokemon/healthcheck", demoApp.Endpoint()))
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
}

func createTest(demoApp test.DemoApp) (openapi.Test, error) {
	importPokemonEndpointURL := fmt.Sprintf("http://%s/pokemon/import")
	resp, err := 
}
