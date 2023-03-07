package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildDeleteRequest(rt ResourceTypeTest, ct contentTypeConverter, testServer *httptest.Server, t *testing.T) *http.Request {
	id := extractID(rt.SampleJSON)
	url := fmt.Sprintf(
		"%s/%s/%s",
		testServer.URL,
		strings.ToLower(rt.ResourceType),
		id,
	)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)
	return req
}

const OperationDeleteSuccess Operation = "DeleteSuccess"

var deleteSuccessOperation = buildSingleStepOperation(singleStepOperationTester{
	name: OperationDeleteSuccess,
	postAssert: func(t *testing.T, ct contentTypeConverter, rt ResourceTypeTest, testServer *httptest.Server) {
		req := buildGetRequest(rt, ct, testServer, t)
		resp := doRequest(t, req, ct.contentType, testServer)
		require.Equal(t, 404, resp.StatusCode)
	},
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildDeleteRequest(rt, ct, testServer, t)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		t.Helper()
		require.Equal(t, 204, resp.StatusCode)
		require.Empty(t, responseBody(t, resp))
	},
})

const OperationDeleteNotFound Operation = "DeleteNotFound"

var deleteNotFoundOperation = buildSingleStepOperation(singleStepOperationTester{
	name: OperationDeleteNotFound,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildDeleteRequest(rt, ct, testServer, t)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		t.Helper()
		require.Equal(t, 404, resp.StatusCode)
	},
})

const OperationDeleteInternalError Operation = "DeleteInternalError"

var deleteInternalErrorOperation = buildSingleStepOperation(singleStepOperationTester{
	name: OperationDeleteInternalError,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildDeleteRequest(rt, ct, testServer, t)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		assertInternalError(t, resp, ct, rt.ResourceType, "deleting")
	},
})
