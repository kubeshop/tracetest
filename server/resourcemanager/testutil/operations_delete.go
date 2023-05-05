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
	"golang.org/x/exp/slices"
)

func buildDeleteRequest(rt ResourceTypeTest, ct contentTypeConverter, testServer *httptest.Server, t *testing.T) *http.Request {
	id := extractID(rt.SampleJSON)
	url := fmt.Sprintf(
		"%s/%s/%s",
		testServer.URL,
		strings.ToLower(rt.ResourceTypePlural),
		id,
	)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)
	return req
}

const OperationDeleteSuccess Operation = "DeleteSuccess"

var deleteSuccessOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationDeleteSuccess,
	neededForOperation: rm.OperationDelete,
	postAssert: func(t *testing.T, ct contentTypeConverter, rt ResourceTypeTest, testServer *httptest.Server) {
		if slices.Contains(rt.operationsWithoutPostAssert, OperationDeleteSuccess) {
			return
		}

		req := buildGetRequest(rt, ct, testServer, t)
		resp := doRequest(t, req, ct.contentType, testServer)
		dumpResponseIfNot(t, assert.Equal(t, 404, resp.StatusCode), resp)
	},
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildDeleteRequest(rt, ct, testServer, t)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		t.Helper()
		dumpResponseIfNot(t, assert.Equal(t, 204, resp.StatusCode), resp)
		dumpResponseIfNot(t, assert.Empty(t, responseBody(t, resp)), resp)
	},
})

const OperationDeleteNotFound Operation = "DeleteNotFound"

var deleteNotFoundOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationDeleteNotFound,
	neededForOperation: rm.OperationDelete,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildDeleteRequest(rt, ct, testServer, t)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		t.Helper()
		dumpResponseIfNot(t, assert.Equal(t, 404, resp.StatusCode), resp)
	},
})

const OperationDeleteInternalError Operation = "DeleteInternalError"

var deleteInternalErrorOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationDeleteInternalError,
	neededForOperation: rm.OperationDelete,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildDeleteRequest(rt, ct, testServer, t)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		assertInternalError(t, resp, ct, rt.ResourceTypeSingular, "deleting")
	},
})
