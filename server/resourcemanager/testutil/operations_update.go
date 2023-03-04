package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildUpdateRequest(rt ResourceTypeTest, ct contentTypeConverter, testServer *httptest.Server, t *testing.T) *http.Request {
	input := ct.fromJSON(rt.SampleJSONUpdated)
	url := fmt.Sprintf("%s/%s/", testServer.URL, strings.ToLower(rt.ResourceType))

	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(input))
	require.NoError(t, err)
	return req
}

const OperationUpdateSuccess Operation = "UpdateSuccess"

type updateSuccessOperation struct{}

func (op updateSuccessOperation) postAssert(t *testing.T, ct contentTypeConverter, rt ResourceTypeTest, testServer *httptest.Server) {
}

func (op updateSuccessOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
	return buildUpdateRequest(rt, ct, testServer, t)
}

func (updateSuccessOperation) name() Operation {
	return OperationUpdateSuccess
}

func (updateSuccessOperation) assertResponse(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
	t.Helper()
	require.Equal(t, 200, resp.StatusCode)

	jsonBody := responseBodyJSON(t, resp, ct)

	expected := ct.toJSON(rt.SampleJSONUpdated)

	require.JSONEq(t, expected, jsonBody)
}

const OperationUpdateNotFound Operation = "UpdateNotFound"

type updateNotFoundOperation struct{}

func (op updateNotFoundOperation) postAssert(t *testing.T, ct contentTypeConverter, rt ResourceTypeTest, testServer *httptest.Server) {
}

func (op updateNotFoundOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
	return buildUpdateRequest(rt, ct, testServer, t)
}

func (updateNotFoundOperation) name() Operation {
	return OperationUpdateNotFound
}

func (updateNotFoundOperation) assertResponse(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
	t.Helper()
	require.Equal(t, 404, resp.StatusCode)
}

const OperationUpdateInternalError Operation = "UpdateInternalError"

type updateInternalErrorOperation struct{}

func (op updateInternalErrorOperation) postAssert(t *testing.T, ct contentTypeConverter, rt ResourceTypeTest, testServer *httptest.Server) {
}

func (op updateInternalErrorOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
	return buildUpdateRequest(rt, ct, testServer, t)
}

func (updateInternalErrorOperation) name() Operation {
	return OperationUpdateInternalError
}

func (updateInternalErrorOperation) assertResponse(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
	assertInternalError(t, resp, ct, rt, "updating")
}
