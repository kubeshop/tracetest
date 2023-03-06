package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

type ResourceTypeTest struct {
	ResourceType      string
	RegisterManagerFn func(*mux.Router) any
	Prepare           func(t *testing.T, operation Operation, bridge any)

	SampleJSON           string
	SampleJSONUpdated    string
	PaginationSortFields []string
}

func TestResourceType(t *testing.T, rt ResourceTypeTest) {
	t.Helper()

	TestResourceTypeOperations(t, rt, defaultOperations)
}

func TestResourceTypeWithErrorOperations(t *testing.T, rt ResourceTypeTest) {
	t.Helper()

	TestResourceTypeOperations(t, rt, append(defaultOperations, errorOperations...))
}

func TestResourceTypeOperations(t *testing.T, rt ResourceTypeTest, operations []operationTester) {
	t.Parallel()
	t.Helper()

	t.Run(rt.ResourceType, func(t *testing.T) {
		for _, op := range operations {
			t.Run(string(op.name), func(t *testing.T) {
				op := op
				t.Parallel()

				testOperation(t, op, rt)
			})
		}
	})
}

func testOperation(t *testing.T, op operationTester, rt ResourceTypeTest) {
	t.Helper()

	for _, ct := range contentTypeConverters {
		t.Run(ct.name, func(t *testing.T) {
			ct := ct
			t.Parallel()

			testOperationForContentType(t, op, ct, rt)
		})
	}
}

func testOperationForContentType(t *testing.T, op operationTester, ct contentTypeConverter, rt ResourceTypeTest) {
	t.Helper()

	router := mux.NewRouter()
	testServer := httptest.NewServer(router)
	testBridge := rt.RegisterManagerFn(router)

	if rt.Prepare != nil {
		rt.Prepare(t, op.name, testBridge)
	}

	req := op.buildRequest(t, testServer, ct, rt)
	resp := doRequest(t, req, ct.contentType, testServer)

	op.assertResponse(t, resp, ct, rt)
	assert.Equal(t, ct.contentType, resp.Header.Get("Content-Type"))
	if op.postAssert != nil {
		op.postAssert(t, ct, rt, testServer)
	}
}

func doRequest(t *testing.T, req *http.Request, contentType string, testServer *httptest.Server) *http.Response {
	req.Header.Set("Content-Type", contentType)
	resp, err := testServer.Client().Do(req)
	require.NoError(t, err)

	return resp
}
