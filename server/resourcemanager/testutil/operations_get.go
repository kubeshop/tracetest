package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildGetRequest(rt ResourceTypeTest, ct contentType, testServer *httptest.Server, t *testing.T) *http.Request {
	id := extractID(rt.SampleJSON)
	url := fmt.Sprintf(
		"%s/%s/%s",
		testServer.URL,
		strings.ToLower(rt.ResourceType),
		id,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	return req
}

const OperationGetSuccess Operation = "GetSuccess"

type getSuccessOperation struct{}

func (op getSuccessOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildGetRequest(rt, ct, testServer, t)
}

func (getSuccessOperation) name() Operation {
	return OperationGetSuccess
}

func (getSuccessOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	t.Helper()
	require.Equal(t, 200, resp.StatusCode)

	jsonBody := responseBodyJSON(t, resp, ct)

	expected := ct.toJSON(rt.SampleJSON)

	require.JSONEq(t, expected, jsonBody)
}

const OperationGetNotFound Operation = "GetNotFound"

type getNotFoundOperation struct{}

func (op getNotFoundOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildGetRequest(rt, ct, testServer, t)
}

func (getNotFoundOperation) name() Operation {
	return OperationGetNotFound
}

func (getNotFoundOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	t.Helper()
	require.Equal(t, 404, resp.StatusCode)
}

const OperationGetInteralError Operation = "GetInteralError"

type getInteralErrorOperation struct{}

func (op getInteralErrorOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildGetRequest(rt, ct, testServer, t)
}

func (getInteralErrorOperation) name() Operation {
	return OperationGetInteralError
}

func (getInteralErrorOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	assertInternalError(t, resp, ct, rt, "getting")
}
