package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/Jeffail/gabs/v2"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func dumpResponseIfNot(t *testing.T, success bool, resp *http.Response) {
	t.Helper()

	if success {
		return
	}

	b, _ := httputil.DumpResponse(resp, true)
	t.Log("\n", string(b))
	t.FailNow()
}

func generateRandomString() string {
	generator := id.NewRandGenerator()
	return generator.TraceID().String()
}

func RemoveFieldFromJSONResource(field, jsonResource string) string {
	jsonParsed, err := gabs.ParseJSON([]byte(jsonResource))
	if err != nil {
		panic(err)
	}

	field = "spec." + field

	jf := jsonParsed.Path(field)
	if jf == nil {
		// field not exists, do nothing
		return jsonResource
	}

	err = jsonParsed.DeleteP(field)
	if err != nil {
		panic(err)
	}

	return jsonParsed.String()
}

func removeIDFromJSON(input string) string {
	return RemoveFieldFromJSONResource("id", input)
}

func parseJSON(input string) map[string]any {
	parsed := map[string]any{}
	err := json.Unmarshal([]byte(input), &parsed)
	if err != nil {
		panic(err)
	}

	return parsed
}

func extractID(input string) string {
	parsed := parseJSON(input)
	id := parsed["spec"].(map[string]any)["id"]
	if id == nil {
		return ""
	}

	return id.(string)
}

func responseBody(t *testing.T, resp *http.Response) string {
	require.NotNil(t, resp.Body)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return string(body)
}

func responseBodyJSON(t *testing.T, resp *http.Response, ct contentTypeConverter) string {
	body := responseBody(t, resp)
	jsonBody := ct.toJSON(string(body))
	return jsonBody
}

func assertInternalError(t *testing.T, resp *http.Response, ct contentTypeConverter, resourceType, verb string) {
	dumpResponseIfNot(t, assert.Equal(t, 500, resp.StatusCode), resp)

	jsonBody := responseBodyJSON(t, resp, ct)

	// hacky way to get the types we want
	bodyValues := struct {
		Code  int    `json:"code"`
		Error string `json:"error"`
	}{}
	json.Unmarshal([]byte(jsonBody), &bodyValues)

	dumpResponseIfNot(t, assert.Equal(t, 500, bodyValues.Code), resp)
	dumpResponseIfNot(t, assert.Contains(t, bodyValues.Error, "error "+verb+" resource "+resourceType), resp)
}
