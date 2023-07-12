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

func buildCreateRequest(body, rt string, ct contentTypeConverter, testServer *httptest.Server, t *testing.T) *http.Request {
	input := ct.fromJSON(body)
	url := fmt.Sprintf("%s/%s", testServer.URL, strings.ToLower(rt))

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(input))
	require.NoError(t, err)
	return req
}

const OperationCreateNoID Operation = "CreateNoID"

var createNoIDOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationCreateNoID,
	neededForOperation: rm.OperationCreate,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildCreateRequest(
			removeIDFromJSON(rt.SampleJSON),
			rt.ResourceTypePlural,
			ct,
			testServer,
			t,
		)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		dumpResponseIfNot(t, assert.Equal(t, 201, resp.StatusCode), resp)

		jsonBody := responseBodyJSON(t, resp, ct)

		clean := removeIDFromJSON(rt.SampleJSON)
		expected := ct.toJSON(clean)

		rt.customJSONComparer(t, OperationCreateNoID, expected, removeIDFromJSON(jsonBody))
		dumpResponseIfNot(t, assert.NotEmpty(t, extractID(jsonBody)), resp)
	},
})

const OperationCreateSuccess Operation = "CreateSuccess"

var createSuccessOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationCreateSuccess,
	neededForOperation: rm.OperationCreate,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildCreateRequest(
			rt.SampleJSON,
			rt.ResourceTypePlural,
			ct,
			testServer,
			t,
		)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		dumpResponseIfNot(t, assert.Equal(t, 201, resp.StatusCode), resp)

		jsonBody := responseBodyJSON(t, resp, ct)
		expected := ct.toJSON(rt.SampleJSON)

		rt.customJSONComparer(t, OperationCreateSuccess, expected, jsonBody)
	},
})

const OperationCreateInternalError Operation = "CreateInternalError"

var createInternalErrorOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationCreateInternalError,
	neededForOperation: rm.OperationCreate,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildCreateRequest(
			rt.SampleJSON,
			rt.ResourceTypePlural,
			ct,
			testServer,
			t,
		)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		assertInternalError(t, resp, ct, rt.ResourceTypeSingular, "creating")
	},
})
