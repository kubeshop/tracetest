package executor_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kubeshop/tracetest/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutorIntegration(t *testing.T) {
	demoAppConfig, err := test.GetDemoApplicationInstance()
	require.NoError(t, err)

	resp, err := http.Get(fmt.Sprintf("http://%s/pokemon/healthcheck", demoAppConfig.Endpoint))
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
}
