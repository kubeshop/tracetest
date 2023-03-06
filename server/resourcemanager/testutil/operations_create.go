package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildCreateRequest(body, rt string, ct contentTypeConverter, testServer *httptest.Server, t *testing.T) *http.Request {
	input := ct.fromJSON(body)
	url := fmt.Sprintf("%s/%s/", testServer.URL, strings.ToLower(rt))

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(input))
	require.NoError(t, err)
	return req
}

const OperationCreateNoID Operation = "CreateNoID"

var createNoIDOperation = operationTester{
	name: OperationCreateNoID,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildCreateRequest(
			removeIDFromJSON(rt.SampleJSON),
			rt.ResourceType,
			ct,
			testServer,
			t,
		)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		require.Equal(t, 201, resp.StatusCode)

		jsonBody := responseBodyJSON(t, resp, ct)

		clean := removeIDFromJSON(rt.SampleJSON)
		expected := ct.toJSON(clean)

		require.JSONEq(t, expected, removeIDFromJSON(jsonBody))
		require.NotEmpty(t, extractID(jsonBody))
	},
}

const OperationCreateSuccess Operation = "CreateSuccess"

var createSuccessOperation = operationTester{
	name: OperationCreateSuccess,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildCreateRequest(
			rt.SampleJSON,
			rt.ResourceType,
			ct,
			testServer,
			t,
		)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		require.Equal(t, 201, resp.StatusCode)

		jsonBody := responseBodyJSON(t, resp, ct)
		expected := ct.toJSON(rt.SampleJSON)

		require.JSONEq(t, expected, jsonBody)
	},
}

const OperationCreateInternalError Operation = "CreateInternalError"

var createInternalErrorOperation = operationTester{
	name: OperationCreateInternalError,
	buildRequest: func(t *testing.T, testServer *httptest.Server, ct contentTypeConverter, rt ResourceTypeTest) *http.Request {
		return buildCreateRequest(
			rt.SampleJSON,
			rt.ResourceType,
			ct,
			testServer,
			t,
		)
	},
	assertResponse: func(t *testing.T, resp *http.Response, ct contentTypeConverter, rt ResourceTypeTest) {
		assertInternalError(t, resp, ct, rt.ResourceType, "creating")
	},
}
