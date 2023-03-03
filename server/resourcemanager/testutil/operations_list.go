package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildListRequest(rt string, ct contentType, testServer *httptest.Server, t *testing.T) *http.Request {
	url := fmt.Sprintf("%s/%s/", testServer.URL, strings.ToLower(rt))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	return req
}

const OperationListNoResults Operation = "ListNoResults"

type listNoResultsOperation struct{}

func (op listNoResultsOperation) postAssert(t *testing.T, ct contentType, rt ResourceTypeTest, testServer *httptest.Server) {
}

func (op listNoResultsOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildListRequest(
		rt.ResourceType,
		ct,
		testServer,
		t,
	)
}

func (listNoResultsOperation) name() Operation {
	return OperationListNoResults
}

func (listNoResultsOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	require.Equal(t, 200, resp.StatusCode)

	jsonBody := responseBodyJSON(t, resp, ct)

	expected := `{
		"count": 0,
		"items": []
	}`

	require.JSONEq(t, expected, jsonBody)
}

const OperationListSuccess Operation = "ListSuccess"

type listSuccessOperation struct{}

func (op listSuccessOperation) postAssert(t *testing.T, ct contentType, rt ResourceTypeTest, testServer *httptest.Server) {
}

func (op listSuccessOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildListRequest(
		rt.ResourceType,
		ct,
		testServer,
		t,
	)
}

func (listSuccessOperation) name() Operation {
	return OperationListSuccess
}

func (listSuccessOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	require.Equal(t, 200, resp.StatusCode)

	jsonBody := responseBodyJSON(t, resp, ct)

	expected := `{
		"count": 1,
		"items": [` + ct.toJSON(rt.SampleJSON) + `]
	}`

	require.JSONEq(t, expected, jsonBody)
}

const OperationListInteralError Operation = "ListInteralError"

type listInteralErrorOperation struct{}

func (op listInteralErrorOperation) postAssert(t *testing.T, ct contentType, rt ResourceTypeTest, testServer *httptest.Server) {
}

func (op listInteralErrorOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildListRequest(
		rt.ResourceType,
		ct,
		testServer,
		t,
	)
}

func (listInteralErrorOperation) name() Operation {
	return OperationListInteralError
}

func (listInteralErrorOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	assertInternalError(t, resp, ct, rt, "listing")
}
