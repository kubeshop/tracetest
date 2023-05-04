package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	rm "github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildAugmentedGetRequest(rt ResourceTypeTest, ct contentTypeConverter, testServer *httptest.Server, t *testing.T) *http.Request {
	id := extractID(rt.SampleJSONAugmented)
	req, err := getRequestForID(id, rt, testServer)
	require.NoError(t, err)
	req.Header.Set(rm.HeaderAugmented, "true")
	return req
}

const OperationGetAugmentedSuccess Operation = "GetAugmentedSuccess"

var getAugmentedSuccessOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationGetAugmentedSuccess,
	neededForOperation: rm.OperationGetAugmented,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildAugmentedGetRequest(rt, ct, testServer, t)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		t.Helper()
		dumpResponseIfNot(t, assert.Equal(t, 200, resp.StatusCode), resp)

		jsonBody := responseBodyJSON(t, resp, ct)

		expected := ct.toJSON(rt.SampleJSONAugmented)

		rt.customJSONComparer(t, OperationGetAugmentedSuccess, expected, jsonBody)
	},
})

func buildAugmentedListRequest(rt ResourceTypeTest, ct contentTypeConverter, testServer *httptest.Server, t *testing.T) *http.Request {
	req := buildListRequest(
		rt.ResourceTypePlural,
		map[string]string{},
		ct,
		testServer,
		t,
	)
	req.Header.Set(rm.HeaderAugmented, "true")
	return req
}

const OperationListAugmentedSuccess Operation = "ListAugmentedSuccess"

var ListAugmentedSuccessOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationListAugmentedSuccess,
	neededForOperation: rm.OperationListAugmented,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildAugmentedListRequest(rt, ct, testServer, t)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		t.Helper()
		dumpResponseIfNot(t, assert.Equal(t, 200, resp.StatusCode), resp)

		jsonBody := responseBodyJSON(t, resp, ct)

		expected := ct.toJSON(rt.SampleJSONAugmented)

		rt.customJSONComparer(t, OperationGetAugmentedSuccess, expected, jsonBody)
	},
})
