package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildQueryString(params map[string]string) string {
	queryString := ""

	if len(params) == 0 {
		return queryString
	}

	formattedParams := []string{}

	for key, value := range params {
		formattedParams = append(formattedParams, fmt.Sprintf("%s=%s", key, value))
	}

	return "?" + strings.Join(formattedParams, "&")
}

func compareFieldsGreaterThan(first, second any) bool {
	switch first.(type) {
	case int:
		return first.(int) >= second.(int)
	case float32:
		return first.(float32) >= second.(float32)
	case float64:
		return first.(float64) >= second.(float64)
	case string:
		return first.(string) >= second.(string)
	default:
		panic("type unknown for comparison")
	}
}

func compareFieldsLesserThan(first, second any) bool {
	switch first.(type) {
	case int:
		return first.(int) <= second.(int)
	case float32:
		return first.(float32) <= second.(float32)
	case float64:
		return first.(float64) <= second.(float64)
	case string:
		return first.(string) <= second.(string)
	default:
		panic("type unknown for comparison")
	}
}

func buildListRequest(resourceType string, paginationParams map[string]string, ct contentTypeConverter, testServer *httptest.Server, t *testing.T) *http.Request {
	queryString := buildQueryString(paginationParams)

	url := fmt.Sprintf("%s/%s/%s", testServer.URL, strings.ToLower(resourceType), queryString)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	return req
}

const OperationListNoResults Operation = "ListNoResults"

var listNoResultsOperation = operationTester{
	name: OperationListNoResults,
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
}

const OperationListSuccess Operation = "ListSuccess"

var listSuccessOperation = operationTester{
	name: OperationListSuccess,
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
}

const OperationListWithInvalidSortField Operation = "ListWithInvalidSortField"

var listWithInvalidSortFieldOperation = operationTester{
	name: OperationListWithInvalidSortField,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildListRequest(
			rt.ResourceType,
			map[string]string{
				"take":          "2",
				"skip":          "1",
				"sortBy":        rt.InvalidSortField,
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
}

const OperationListPaginatedAscendingSuccess Operation = "ListPaginatedAscendingSuccess"

var listPaginatedAscendingSuccessOperation = operationTester{
	name: OperationListPaginatedAscendingSuccess,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildListRequest(
			rt.ResourceType,
			map[string]string{
				"take":          "2",
				"skip":          "1",
				"sortBy":        rt.SortField,
				"sortDirection": "asc",
			},
			ct,
			testServer,
			t,
		)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		require.Equal(t, 200, resp.StatusCode)

		jsonBody := responseBodyJSON(t, resp, ct)

		var parsedJsonBody struct {
			Count int              `json:count`
			Items []map[string]any `json:items`
		}
		json.Unmarshal([]byte(jsonBody), &parsedJsonBody)

		require.Equal(t, 3, parsedJsonBody.Count)

		var prevVal any
		field := rt.SortField

		for _, item := range parsedJsonBody.Items {
			if prevVal == nil {
				prevVal = item[field]
				continue
			}
			assert.True(t, compareFieldsLesserThan(prevVal, item[field]))

			prevVal = item[field]
		}
	},
}

const OperationListPaginatedDescendingSuccess Operation = "ListPaginatedDescendingSuccess"

var listPaginatedDescendingSuccessOperation = operationTester{
	name: OperationListPaginatedDescendingSuccess,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildListRequest(
			rt.ResourceType,
			map[string]string{
				"take":          "2",
				"skip":          "1",
				"sortBy":        rt.SortField,
				"sortDirection": "desc",
			},
			ct,
			testServer,
			t,
		)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		require.Equal(t, 200, resp.StatusCode)

		jsonBody := responseBodyJSON(t, resp, ct)

		var parsedJsonBody struct {
			Count int              `json:count`
			Items []map[string]any `json:items`
		}
		json.Unmarshal([]byte(jsonBody), &parsedJsonBody)

		require.Equal(t, 3, parsedJsonBody.Count)

		var prevVal any
		field := rt.SortField

		for _, item := range parsedJsonBody.Items {
			if prevVal == nil {
				prevVal = item[field]
				continue
			}
			assert.True(t, compareFieldsGreaterThan(prevVal, item[field]))

			prevVal = item[field]
		}
	},
}

const OperationListInternalError Operation = "ListInternalError"

var listInternalErrorOperation = operationTester{
	name: OperationListInternalError,
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
}
