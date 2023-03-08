package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	rm "github.com/kubeshop/tracetest/server/resourcemanager"
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

var updateSuccessOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationUpdateSuccess,
	neededForOperation: rm.OperationUpdate,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildUpdateRequest(rt, ct, testServer, t)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		t.Helper()
		require.Equal(t, 200, resp.StatusCode)

		jsonBody := responseBodyJSON(t, resp, ct)

		expected := ct.toJSON(rt.SampleJSONUpdated)

		require.JSONEq(t, expected, jsonBody)
	},
})

const OperationUpdateNotFound Operation = "UpdateNotFound"

var updateNotFoundOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationUpdateNotFound,
	neededForOperation: rm.OperationUpdate,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildUpdateRequest(rt, ct, testServer, t)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		t.Helper()
		require.Equal(t, 404, resp.StatusCode)
	},
})

const OperationUpdateInternalError Operation = "UpdateInternalError"

var updateInternalErrorOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationUpdateInternalError,
	neededForOperation: rm.OperationUpdate,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildUpdateRequest(rt, ct, testServer, t)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		assertInternalError(t, resp, ct, rt.ResourceType, "updating")
	},
})
