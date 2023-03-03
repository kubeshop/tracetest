package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildCreateRequest(body, rt string, ct contentType, testServer *httptest.Server, t *testing.T) *http.Request {
	input := ct.fromJSON(body)
	url := fmt.Sprintf("%s/%s/", testServer.URL, strings.ToLower(rt))

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(input))
	require.NoError(t, err)
	return req
}

const OperationCreateNoID Operation = "CreateNoID"

type createNoIDOperation struct{}

func (op createNoIDOperation) postAssert(t *testing.T, ct contentType, rt ResourceTypeTest, testServer *httptest.Server) {
}

func (op createNoIDOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildCreateRequest(
		removeIDFromJSON(rt.SampleJSON),
		rt.ResourceType,
		ct,
		testServer,
		t,
	)
}

func (createNoIDOperation) name() Operation {
	return OperationCreateNoID
}

func (createNoIDOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	require.Equal(t, 201, resp.StatusCode)

	jsonBody := responseBodyJSON(t, resp, ct)

	clean := removeIDFromJSON(rt.SampleJSON)
	expected := ct.toJSON(clean)

	require.JSONEq(t, expected, removeIDFromJSON(jsonBody))
	require.NotEmpty(t, extractID(jsonBody))
}

const OperationCreateSuccess Operation = "CreateSuccess"

type createSuccessOperation struct{}

func (op createSuccessOperation) postAssert(t *testing.T, ct contentType, rt ResourceTypeTest, testServer *httptest.Server) {
}

func (op createSuccessOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildCreateRequest(
		rt.SampleJSON,
		rt.ResourceType,
		ct,
		testServer,
		t,
	)
}

func (createSuccessOperation) name() Operation {
	return OperationCreateSuccess
}

func (createSuccessOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	require.Equal(t, 201, resp.StatusCode)

	jsonBody := responseBodyJSON(t, resp, ct)
	expected := ct.toJSON(rt.SampleJSON)

	require.JSONEq(t, expected, jsonBody)
}

const OperationCreateInternalError Operation = "CreateInternalError"

type createInternalErrorOperation struct{}

func (op createInternalErrorOperation) postAssert(t *testing.T, ct contentType, rt ResourceTypeTest, testServer *httptest.Server) {
}

func (op createInternalErrorOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt ResourceTypeTest) *http.Request {
	return buildCreateRequest(
		rt.SampleJSON,
		rt.ResourceType,
		ct,
		testServer,
		t,
	)
}

func (createInternalErrorOperation) name() Operation {
	return OperationCreateInternalError
}

func (createInternalErrorOperation) assertResponse(t *testing.T, resp *http.Response, ct contentType, rt ResourceTypeTest) {
	assertInternalError(t, resp, ct, rt, "creating")
}
