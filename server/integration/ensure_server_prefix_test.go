package integration_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerPrefix(t *testing.T) {
	_, err := testfixtures.GetTracetestApp(
		testfixtures.WithServerPrefix("/tracetest"),
		testfixtures.WithHttpPort(8000),
		testfixtures.WithWebSocketPort(8001),
	)
	require.NoError(t, err)

	expectedEndpoint := "http://localhost:8000/tracetest"
	tests := getTests(t, expectedEndpoint)
	assert.NotNil(t, tests)
}

func getTests(t *testing.T, endpoint string) []openapi.Test {
	url := fmt.Sprintf("%s/api/tests", endpoint)
	resp, err := http.Get(url)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	bodyJsonBytes, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var tests []openapi.Test
	err = json.Unmarshal(bodyJsonBytes, &tests)
	require.NoError(t, err)

	return tests
}
