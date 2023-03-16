package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"

	rm "github.com/kubeshop/tracetest/server/resourcemanager"
)

type ResourceTypeTest struct {
	ResourceType      string
	RegisterManagerFn func(*mux.Router) rm.Manager
	Prepare           func(t *testing.T, operation Operation, manager rm.Manager)

	SampleJSON        string
	SampleJSONUpdated string

	// private files
	sortFields []string
}

type config struct {
	operations operationTesters
}

type testOption func(*config)

func ExcludeOperations(ops ...Operation) testOption {
	return func(c *config) {
		c.operations = c.operations.exclude(ops...)
	}
}

func TestResourceType(t *testing.T, rt ResourceTypeTest, opts ...testOption) {
	t.Helper()

	cfg := config{
		operations: defaultOperations,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	TestResourceTypeOperations(t, rt, cfg.operations)
}

func TestResourceTypeWithErrorOperations(t *testing.T, rt ResourceTypeTest) {
	t.Helper()

	TestResourceTypeOperations(t, rt, append(defaultOperations, errorOperations...))
}

func TestResourceTypeOperations(t *testing.T, rt ResourceTypeTest, operations []operationTester) {
	t.Parallel()
	t.Helper()

	t.Run(rt.ResourceType, func(t *testing.T) {
		testProvisioning(t, rt)

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
	manager := rt.RegisterManagerFn(router)

	sortable, ok := manager.Handler().(rm.SortableHandler)
	if ok {
		rt.sortFields = sortable.SortingFields()
	}

	if rt.Prepare != nil {
		rt.Prepare(t, op.name, manager)
	}

	if !op.needsToRun(manager.EnabledOperations()) {
		t.Skipf("operation '%s' not enabled", op.name)
	}

	operationSteps := op.getSteps(t, rt)

	for _, step := range operationSteps {
		req := step.buildRequest(t, testServer, ct, rt)
		resp := doRequest(t, req, ct.contentType, testServer)

		step.assertResponse(t, resp, ct, rt)
		assert.Equal(t, ct.contentType, resp.Header.Get("Content-Type"))
		if step.postAssert != nil {
			step.postAssert(t, ct, rt, testServer)
		}
	}
}

func doRequest(t *testing.T, req *http.Request, contentType string, testServer *httptest.Server) *http.Response {
	req.Header.Set("Content-Type", contentType)
	resp, err := testServer.Client().Do(req)
	require.NoError(t, err)

	return resp
}
