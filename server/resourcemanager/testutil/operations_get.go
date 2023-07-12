package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	rm "github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getRequestForID(id string, rt ResourceTypeTest, testServer *httptest.Server) (*http.Request, error) {
	url := fmt.Sprintf(
		"%s/%s/%s",
		testServer.URL,
		strings.ToLower(rt.ResourceTypePlural),
		id,
	)

	return http.NewRequest(http.MethodGet, url, nil)
}

func buildGetRequest(rt ResourceTypeTest, ct contentTypeConverter, testServer *httptest.Server, t *testing.T) *http.Request {
	id := extractID(rt.SampleJSON)
	req, err := getRequestForID(id, rt, testServer)
	require.NoError(t, err)
	return req
}

const OperationGetSuccess Operation = "GetSuccess"

var getSuccessOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationGetSuccess,
	neededForOperation: rm.OperationGet,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildGetRequest(rt, ct, testServer, t)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		t.Helper()
		dumpResponseIfNot(t, assert.Equal(t, 200, resp.StatusCode), resp)

		jsonBody := responseBodyJSON(t, resp, ct)

		expected := ct.toJSON(rt.SampleJSON)

		rt.customJSONComparer(t, OperationGetSuccess, expected, jsonBody)
	},
})

const OperationGetNotFound Operation = "GetNotFound"

var getNotFoundOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationGetNotFound,
	neededForOperation: rm.OperationGet,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildGetRequest(rt, ct, testServer, t)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		t.Helper()
		dumpResponseIfNot(t, assert.Equal(t, 404, resp.StatusCode), resp)
	},
})

const OperationGetInternalError Operation = "GetInternalError"

var getInternalErrorOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationGetInternalError,
	neededForOperation: rm.OperationGet,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildGetRequest(rt, ct, testServer, t)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		assertInternalError(t, resp, ct, rt.ResourceTypeSingular, "getting")
	},
})
