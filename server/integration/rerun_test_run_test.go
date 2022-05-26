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

func TestRerun(t *testing.T) {
	importPokemonTest, err := testfixtures.GetPokemonTest()
	require.NoError(t, err)

	importPokemonTestRun, err := testfixtures.GetPokemonTestRun()
	require.NoError(t, err)

	newTestRun := rerunTestRun(t, importPokemonTest, importPokemonTestRun)
	assert.NoError(t, err)

	assert.NotEqual(t, importPokemonTestRun.Id, newTestRun.Id)
	assert.Equal(t, importPokemonTestRun.Request, newTestRun.Request)
	assert.Equal(t, importPokemonTestRun.Response, newTestRun.Response)
	assert.Equal(t, importPokemonTestRun.Trace, newTestRun.Trace)
}

func rerunTestRun(t *testing.T, test *openapi.Test, testRun *openapi.TestRun) *openapi.TestRun {
	url := fmt.Sprintf("%s/api/tests/%s/run/%s/rerun", endpointUrl, test.Id, testRun.Id)
	response, err := http.Post(url, "application/json", nil)
	require.NoError(t, err)

	if response.StatusCode != 200 {
		require.Fail(t, fmt.Sprintf("expected status 200, got %d", response.StatusCode))
	}

	responseBody := openapi.TestRun{}
	bodyContent, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	err = json.Unmarshal(bodyContent, &responseBody)
	require.NoError(t, err)

	return &responseBody
}
