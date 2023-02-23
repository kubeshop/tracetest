package testutil

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var operations = []operationTester{
	createSuccessOperation{},
	createInteralErrorOperation{},
}

const OperationCreateSuccess Operation = "CreateSuccess"

type createSuccessOperation struct{}

func (op createSuccessOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt *ResourceTypeTest) *http.Request {
	input := ct.fromJSON(rt.SampleNewJSON)
	url := fmt.Sprintf("%s/%s/", testServer.URL, strings.ToLower(rt.ResourceType))

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(input))
	require.NoError(t, err)

	return req
}

func (createSuccessOperation) name() Operation {
	return OperationCreateSuccess
}

func (createSuccessOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt *ResourceTypeTest) {
	assert.Equal(t, resp.StatusCode, 201)

	require.NotNil(t, resp.Body)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.JSONEq(t, rt.SampleCreatedJSON, ct.toJSON(string(body)))
}

const OperationCreateInteralError Operation = "CreateInteralError"

type createInteralErrorOperation struct{}

func (op createInteralErrorOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt *ResourceTypeTest) *http.Request {
	input := ct.fromJSON(rt.SampleNewJSON)
	url := fmt.Sprintf("%s/%s/", testServer.URL, strings.ToLower(rt.ResourceType))

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(input))
	require.NoError(t, err)

	return req
}

func (createInteralErrorOperation) name() Operation {
	return OperationCreateInteralError
}

func (createInteralErrorOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt *ResourceTypeTest) {
	assert.Equal(t, 500, resp.StatusCode)

	require.NotNil(t, resp.Body)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Contains(t, string(body), "error creating resource "+rt.ResourceType)
}
