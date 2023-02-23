package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var operations = []operationTester{
	createSuccessOperation{},
}

const OperationCreateSuccess Operation = "CreateSuccess"

type createSuccessOperation struct{}

func (op createSuccessOperation) buildRequest(t *testing.T, testServer *httptest.Server, ct contentType, rt *ResourceTypeTest) *http.Request {
	input := ct.fromJSON(rt.SampleNewJSON)
	url := fmt.Sprintf("%s/%s/", testServer.URL, strings.ToLower(rt.ResourceType))

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(input))
	require.NoError(t, err)

	return req
}

func (_ createSuccessOperation) name() Operation {
	return OperationCreateSuccess
}
