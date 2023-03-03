package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildDeleteRequest(rt ResourceTypeTest, ct contentType, testServer *httptest.Server, t *testing.T) *http.Request {
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

type deleteSuccessOperation struct{}

func (op deleteSuccessOperation) postAssert(t *testing.T, ct contentType, rt ResourceTypeTest, testServer *httptest.Server) {
	req := buildGetRequest(rt, ct, testServer, t)

	resp := doRequest(t, req, ct, testServer)

	(getNotFoundOperation{}).assertResponse(t, resp, ct, rt)
}

func (op deleteSuccessOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildDeleteRequest(rt, ct, testServer, t)
}

func (deleteSuccessOperation) name() Operation {
	return OperationDeleteSuccess
}

func (deleteSuccessOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	t.Helper()
	require.Equal(t, 204, resp.StatusCode)
	require.Empty(t, responseBody(t, resp))
}

const OperationDeleteNotFound Operation = "DeleteNotFound"

type deleteNotFoundOperation struct{}

func (op deleteNotFoundOperation) postAssert(t *testing.T, ct contentType, rt ResourceTypeTest, testServer *httptest.Server) {
}

func (op deleteNotFoundOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildDeleteRequest(rt, ct, testServer, t)
}

func (deleteNotFoundOperation) name() Operation {
	return OperationDeleteNotFound
}

func (deleteNotFoundOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	t.Helper()
	require.Equal(t, 404, resp.StatusCode)
}

const OperationDeleteInternalError Operation = "DeleteInternalError"

type deleteInternalErrorOperation struct{}

func (op deleteInternalErrorOperation) postAssert(t *testing.T, ct contentType, rt ResourceTypeTest, testServer *httptest.Server) {
}

func (op deleteInternalErrorOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildDeleteRequest(rt, ct, testServer, t)
}

func (deleteInternalErrorOperation) name() Operation {
	return OperationDeleteInternalError
}

func (deleteInternalErrorOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	assertInternalError(t, resp, ct, rt, "deleting")
}
