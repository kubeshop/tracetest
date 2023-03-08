package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	rm "github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildQueryString(params map[string]string) url.Values {
	vals := url.Values{}

	for k, v := range params {
		vals.Add(k, v)
	}

	return vals
}

func buildListRequest(resourceType string, paginationParams map[string]string, ct contentTypeConverter, testServer *httptest.Server, t *testing.T) *http.Request {
	qs := buildQueryString(paginationParams)

	url := fmt.Sprintf("%s/%s/", testServer.URL, strings.ToLower(resourceType))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	req.URL.RawQuery = qs.Encode()

	return req
}

const OperationListNoResults Operation = "ListNoResults"

var listNoResultsOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationListNoResults,
	neededForOperation: rm.OperationList,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildListRequest(
			rt.ResourceType,
			map[string]string{},
			ct,
			testServer,
			t,
		)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		require.Equal(t, 200, resp.StatusCode)

		jsonBody := responseBodyJSON(t, resp, ct)

		expected := `{
			"count": 0,
			"items": []
		}`

		require.JSONEq(t, expected, jsonBody)
	},
})

const OperationListSuccess Operation = "ListSuccess"

var listSuccessOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationListSuccess,
	neededForOperation: rm.OperationList,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildListRequest(
			rt.ResourceType,
			map[string]string{},
			ct,
			testServer,
			t,
		)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		require.Equal(t, 200, resp.StatusCode)

		jsonBody := responseBodyJSON(t, resp, ct)

		expected := `{
			"count": 1,
			"items": [` + ct.toJSON(rt.SampleJSON) + `]
		}`

		require.JSONEq(t, expected, jsonBody)
	},
})

const OperationListWithInvalidSortField Operation = "ListWithInvalidSortField"

var listWithInvalidSortFieldOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationListWithInvalidSortField,
	neededForOperation: rm.OperationList,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		invalidSortField := generateRandomString()

		return buildListRequest(
			rt.ResourceType,
			map[string]string{
				"take":          "2",
				"skip":          "1",
				"sortBy":        invalidSortField,
				"sortDirection": "asc",
			},
			ct,
			testServer,
			t,
		)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		require.Equal(t, 400, resp.StatusCode)
	},
})

const OperationListPaginatedSuccess Operation = "ListPaginatedSuccess"

var listPaginatedSuccessOperation = operationTester{
	name:               OperationListPaginatedSuccess,
	neededForOperation: rm.OperationList,
	getSteps: func(t *testing.T, rt ResourceTypeTest) []operationTesterStep {
		steps := []operationTesterStep{}

		if len(rt.sortFields) < 1 {
			panic(fmt.Errorf("trying to test list pagination but no sort field was provided"))
		}
		for _, sortField := range rt.sortFields {
			steps = append(steps,
				buildPaginationOperationStep("asc", sortField),
				buildPaginationOperationStep("desc", sortField),
			)
		}

		return steps
	},
}

func buildPaginationOperationStep(sortDirection, sortField string) operationTesterStep {
	return operationTesterStep{
		buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
			return buildListRequest(
				rt.ResourceType,
				map[string]string{
					"take":          "2",
					"skip":          "1",
					"sortBy":        sortField,
					"sortDirection": sortDirection,
				},
				ct,
				testServer,
				t,
			)
		},
		assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
			sortField := sortField
			require.Equal(t, 200, resp.StatusCode)

			jsonBody := responseBodyJSON(t, resp, ct)

			var parsedJsonBody struct {
				Count int              `json:"count"`
				Items []map[string]any `json:"items"`
			}
			json.Unmarshal([]byte(jsonBody), &parsedJsonBody)

			require.Equal(t, 3, parsedJsonBody.Count)
			require.Greater(t, len(parsedJsonBody.Items), 1)

			// we skip the 1st item, so starting in 1 instead of 0
			// makes things match later when comparing to len(parsedJsonBody.Items)
			asserted := 1
			var prevVal any
			for _, item := range parsedJsonBody.Items {
				itemSpec := item["spec"].(map[string]any)
				if prevVal == nil {
					prevVal = itemSpec[sortField]
					continue
				}

				if sortDirection == "asc" {
					assert.LessOrEqual(t, prevVal, itemSpec[sortField])
					asserted++
				} else {
					assert.GreaterOrEqual(t, prevVal, itemSpec[sortField])
					asserted++
				}
			}
			assert.Equal(t, len(parsedJsonBody.Items), asserted)
		},
	}
}

const OperationListInternalError Operation = "ListInternalError"

var listInternalErrorOperation = buildSingleStepOperation(singleStepOperationTester{
	name:               OperationListInternalError,
	neededForOperation: rm.OperationList,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildListRequest(
			rt.ResourceType,
			map[string]string{},
			ct,
			testServer,
			t,
		)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		assertInternalError(t, resp, ct, rt.ResourceType, "listing")
	},
})
