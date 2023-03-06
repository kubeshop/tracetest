package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

const OperationListPaginatedAscendingSuccess Operation = "ListPaginatedAscendingSuccess"

var listPaginatedAscendingSuccessOperation = operationTester{
	name: OperationListPaginatedAscendingSuccess,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildListRequest(
			rt.ResourceType,
			map[string]string{
				"take":          "2",
				"skip":          "1",
				"sortBy":        "id",
				"sortDirection": "asc",
				// TODO: think how to use "query" param
			},
			ct,
			testServer,
			t,
		)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		require.Equal(t, 200, resp.StatusCode)

		jsonBody := responseBodyJSON(t, resp, ct)

		expected := `{
			"count": 3,
			"items": ` + ct.toJSON(rt.SamplePaginatedAscJSON) + `
		}`

		require.JSONEq(t, expected, jsonBody)
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
				"sortBy":        "id",
				"sortDirection": "desc",
				// TODO: think how to use "query" param
			},
			ct,
			testServer,
			t,
		)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		require.Equal(t, 200, resp.StatusCode)

		jsonBody := responseBodyJSON(t, resp, ct)

		expected := `{
			"count": 3,
			"items": ` + ct.toJSON(rt.SamplePaginatedDescJSON) + `
		}`

		require.JSONEq(t, expected, jsonBody)
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
		assertInternalError(t, resp, ct, rt, "listing")
	},
}
