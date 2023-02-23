package testutil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var operations = []operationTester{
	createSuccessOperation{},
	createInteralErrorOperation{},
}

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

const OperationCreateSuccess Operation = "CreateSuccess"

type createSuccessOperation struct{}

func (op createSuccessOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	clean := removeIDFromJSON(rt.SampleJSON)
	input := ct.fromJSON(clean)
	url := fmt.Sprintf("%s/%s/", testServer.URL, strings.ToLower(rt.ResourceType))

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(input))
	require.NoError(t, err)

	return req
}

func (createSuccessOperation) name() Operation {
	return OperationCreateSuccess
}

func (createSuccessOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	require.Equal(t, 201, resp.StatusCode)

	require.NotNil(t, resp.Body)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	jsonBody := ct.toJSON(string(body))

	clean := removeIDFromJSON(rt.SampleJSON)
	expected := ct.toJSON(clean)

	require.JSONEq(t, expected, removeIDFromJSON(jsonBody))
	require.NotEmpty(t, extractID(jsonBody))
}

const OperationCreateInteralError Operation = "CreateInteralError"

type createInteralErrorOperation struct{}

func (op createInteralErrorOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	clean := removeIDFromJSON(rt.SampleJSON)
	input := ct.fromJSON(clean)
	url := fmt.Sprintf("%s/%s/", testServer.URL, strings.ToLower(rt.ResourceType))

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(input))
	require.NoError(t, err)

	return req
}

func (createInteralErrorOperation) name() Operation {
	return OperationCreateInteralError
}

func (createInteralErrorOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	require.Equal(t, 500, resp.StatusCode)

	require.NotNil(t, resp.Body)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Contains(t, string(body), "error creating resource "+rt.ResourceType)
}
