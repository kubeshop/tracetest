package integration_test

import (
	"bytes"
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

func TestAddAssertionToTest(t *testing.T) {
	test, err := testfixtures.GetPokemonTest()
	require.NoError(t, err)

	setDefinition(t, test)

	updatedTest := getTest(t, test.Id)
	assert.Equal(t, test.Version+1, updatedTest.Version)
}

func setDefinition(t *testing.T, test *openapi.Test) {
	url := fmt.Sprintf("%s/api/tests/%s/definition", endpointUrl, test.Id)
	body := openapi.TestDefinition{
		Definitions: []openapi.TestDefinitionDefinitions{
			{
				Selector: `span[tracetest.span.type="http"]`,
				Assertions: []openapi.Assertion{
					{
						Attribute:  "http.status_code",
						Comparator: "=",
						Expected:   "200",
					},
				},
			},
		},
	}
	bodyJson, err := json.Marshal(body)
	require.NoError(t, err)

	bodyBuffer := bytes.NewBuffer(bodyJson)
	req, err := http.NewRequest(http.MethodPut, url, bodyBuffer)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	if resp.StatusCode != 204 {
		require.Fail(t, fmt.Sprintf("expected status 204, got %d", resp.StatusCode))
	}
}

func getTest(t *testing.T, id string) openapi.Test {
	url := fmt.Sprintf("%s/api/tests/%s", endpointUrl, id)
	req, err := http.Get(url)
	require.NoError(t, err)

	bodyJsonBytes, err := ioutil.ReadAll(req.Body)
	require.NoError(t, err)

	var test openapi.Test
	err = json.Unmarshal(bodyJsonBytes, &test)
	require.NoError(t, err)

	return test
}
