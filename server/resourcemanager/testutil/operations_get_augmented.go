package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	rm "github.com/kubeshop/tracetest/server/resourcemanager"
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
		require.Equal(t, 200, resp.StatusCode)

		jsonBody := responseBodyJSON(t, resp, ct)

		expected := ct.toJSON(rt.SampleJSONAugmented)

		rt.customJSONComparer(t, OperationGetAugmentedSuccess, expected, jsonBody)
	},
})
