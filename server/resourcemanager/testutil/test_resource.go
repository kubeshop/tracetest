package testutil

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"

	rm "github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/testmock"
)

type ResourceTypeTest struct {
	ResourceTypeSingular string
	ResourceTypePlural   string
	RegisterManagerFn    func(*mux.Router, *sql.DB) rm.Manager
	Prepare              func(t *testing.T, operation Operation, manager rm.Manager)

	SampleJSON        string
	SampleJSONUpdated string

	// private fields
	sortFields         []string
	customJSONComparer func(t require.TestingT, operation Operation, firstValue, secondValue string)
}

type config struct {
	operations         operationTesters
	customJSONComparer func(t require.TestingT, operation Operation, firstValue, secondValue string)
}

type testOption func(*config)

func ExcludeOperations(ops ...Operation) testOption {
	return func(c *config) {
		c.operations = c.operations.exclude(ops...)
	}
}

func JSONComparer(comparer func(t require.TestingT, operation Operation, firstValue, secondValue string)) testOption {
	return func(c *config) {
		c.customJSONComparer = comparer
	}
}

func rawJSONComparer(t require.TestingT, operation Operation, firstValue, secondValue string) {
	require.JSONEq(t, firstValue, secondValue)
}

func TestResourceType(t *testing.T, rt ResourceTypeTest, opts ...testOption) {
	t.Helper()

	cfg := config{
		operations: defaultOperations,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	// consider customJSONComparer option
	if cfg.customJSONComparer == nil {
		cfg.customJSONComparer = rawJSONComparer
	}
	rt.customJSONComparer = cfg.customJSONComparer

	TestResourceTypeOperations(t, rt, cfg.operations)
}

func TestResourceTypeWithErrorOperations(t *testing.T, rt ResourceTypeTest) {
	t.Helper()

	// assumes default
	rt.customJSONComparer = rawJSONComparer

	TestResourceTypeOperations(t, rt, append(defaultOperations, errorOperations...))
}

func TestResourceTypeOperations(t *testing.T, rt ResourceTypeTest, operations []operationTester) {
	t.Parallel()
	t.Helper()

	if rt.ResourceTypeSingular == "" {
		t.Fatalf("ResourceTypeTest.ResourceTypeSingular not set")
	}

	if rt.ResourceTypePlural == "" {
		t.Fatalf("ResourceTypeTest.ResourceTypePlural not set")
	}

	t.Run(rt.ResourceTypeSingular, func(t *testing.T) {
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

			testOperationForContentType(t, op, ct, rt)
		})
	}
}

func testOperationForContentType(t *testing.T, op operationTester, ct contentTypeConverter, rt ResourceTypeTest) {
	t.Helper()

	db := testmock.CreateMigratedDatabase()
	defer db.Close()

	router := mux.NewRouter()

	testServer := httptest.NewServer(router)
	defer testServer.Close()

	manager := rt.RegisterManagerFn(router, db)

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
