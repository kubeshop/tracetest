package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildCreateRequest(rt ResourceTypeTest, ct contentType, testServer *httptest.Server, t *testing.T) *http.Request {
	clean := removeIDFromJSON(rt.SampleJSON)
	input := ct.fromJSON(clean)
	url := fmt.Sprintf("%s/%s/", testServer.URL, strings.ToLower(rt.ResourceType))

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(input))
	require.NoError(t, err)
	return req
}

const OperationCreateSuccess Operation = "CreateSuccess"

type createSuccessOperation struct{}

func (op createSuccessOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildCreateRequest(rt, ct, testServer, t)
}

func (createSuccessOperation) name() Operation {
	return OperationCreateSuccess
}

func (createSuccessOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	require.Equal(t, 201, resp.StatusCode)

	jsonBody := responseBodyJSON(t, resp, ct)

	clean := removeIDFromJSON(rt.SampleJSON)
	expected := ct.toJSON(clean)

	require.JSONEq(t, expected, removeIDFromJSON(jsonBody))
	require.NotEmpty(t, extractID(jsonBody))
}

const OperationCreateInteralError Operation = "CreateInteralError"

type createInteralErrorOperation struct{}

func (op createInteralErrorOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildCreateRequest(rt, ct, testServer, t)
}

func (createInteralErrorOperation) name() Operation {
	return OperationCreateInteralError
}

func (createInteralErrorOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	assertInternalError(t, resp, ct, rt, "creating")
}
