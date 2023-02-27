package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildUpdateRequest(rt ResourceTypeTest, ct contentType, testServer *httptest.Server, t *testing.T) *http.Request {
	input := ct.fromJSON(rt.SampleJSONUpdated)
	url := fmt.Sprintf("%s/%s/", testServer.URL, strings.ToLower(rt.ResourceType))

	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(input))
	require.NoError(t, err)
	return req
}

const OperationUpdateSuccess Operation = "UpdateSuccess"

type updateSuccessOperation struct{}

func (op updateSuccessOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildUpdateRequest(rt, ct, testServer, t)
}

func (updateSuccessOperation) name() Operation {
	return OperationUpdateSuccess
}

func (updateSuccessOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	t.Helper()
	require.Equal(t, 200, resp.StatusCode)

	jsonBody := responseBodyJSON(t, resp, ct)

	expected := ct.toJSON(rt.SampleJSONUpdated)

	require.JSONEq(t, expected, jsonBody)
}

const OperationUpdateNotFound Operation = "UpdateNotFound"

type updateNotFoundOperation struct{}

func (op updateNotFoundOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildUpdateRequest(rt, ct, testServer, t)
}

func (updateNotFoundOperation) name() Operation {
	return OperationUpdateNotFound
}

func (updateNotFoundOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	t.Helper()
	require.Equal(t, 404, resp.StatusCode)
}

const OperationUpdateInteralError Operation = "UpdateInteralError"

type updateInteralErrorOperation struct{}

func (op updateInteralErrorOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildUpdateRequest(rt, ct, testServer, t)
}

func (updateInteralErrorOperation) name() Operation {
	return OperationUpdateInteralError
}

func (updateInteralErrorOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	assertInternalError(t, resp, ct, rt, "updating")
}
