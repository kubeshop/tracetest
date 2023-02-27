package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	defaultOperations = []operationTester{
		createSuccessOperation{},
		updateSuccessOperation{},
		getSuccessOperation{},
	}

	errorOperations = []operationTester{
		createInteralErrorOperation{},
		updateInteralErrorOperation{},
		getInteralErrorOperation{},
	}
)

func removeID(input map[string]any) map[string]any {
	out := map[string]any{}
	out["type"] = input["type"]
	newSpec := map[string]any{}
	for k, v := range input["spec"].(map[string]any) {
		if k == "id" {
			continue
		}
		newSpec[k] = v
	}

	out["spec"] = newSpec

	return out
}

func parseJSON(input string) map[string]any {
	parsed := map[string]any{}
	err := json.Unmarshal([]byte(input), &parsed)
	if err != nil {
		panic(err)
	}

	return parsed
}

func removeIDFromJSON(input string) string {

	clean := removeID(parseJSON(input))

	out, err := json.Marshal(clean)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func extractID(input string) string {
	parsed := parseJSON(input)
	id := parsed["spec"].(map[string]any)["id"]
	if id == nil {
		return ""
	}

	return id.(string)
}

func responseBodyJSON(t *testing.T, resp *http.Response, ct contentType) string {
	require.NotNil(t, resp.Body)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	jsonBody := ct.toJSON(string(body))
	return jsonBody
}

func assertInternalError(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest, verb string) {
	require.Equal(t, 500, resp.StatusCode)

	jsonBody := responseBodyJSON(t, resp, ct)

	// hacky way to get the types we want
	bodyValues := struct {
		Code  int    `json:"code"`
		Error string `json:"error"`
	}{}
	json.Unmarshal([]byte(jsonBody), &bodyValues)

	require.Equal(t, 500, bodyValues.Code)
	require.Contains(t, bodyValues.Error, "error "+verb+" resource "+rt.ResourceType)
}
