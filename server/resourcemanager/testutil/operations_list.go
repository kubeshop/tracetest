package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

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
}

// Tech debt: Even being one operation, we will split ASC and DESC in two just to keep
// "one request per operation". In the future we can change this structure to support both tests
// in one operationTester instance
const OperationListPaginatedSuccess Operation = "ListPaginatedSuccess"

var listPaginatedAscendingSuccessOperation = operationTester{
	name: OperationListPaginatedSuccess,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildListRequest(
			rt.ResourceType,
			map[string]string{
				"take":          "2",
				"skip":          "1",
				"sortBy":        rt.sortFields[0],
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
		field := rt.sortFields[0]

		for _, item := range parsedJsonBody.Items {
			if prevVal == nil {
				prevVal = item[field]
				continue
			}
			assert.LessOrEqual(t, prevVal, item[field])

			prevVal = item[field]
		}
	},
}

var listPaginatedDescendingSuccessOperation = operationTester{
	name: OperationListPaginatedSuccess,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildListRequest(
			rt.ResourceType,
			map[string]string{
				"take":          "2",
				"skip":          "1",
				"sortBy":        rt.sortFields[0],
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
		field := rt.sortFields[0]

		for _, item := range parsedJsonBody.Items {
			if prevVal == nil {
				prevVal = item[field]
				continue
			}
			assert.GreaterOrEqual(t, prevVal, item[field])

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
