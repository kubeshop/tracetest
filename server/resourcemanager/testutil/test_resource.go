package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
	"sigs.k8s.io/yaml"
)

type Operation string

type ResourceTypeTest struct {
	ResourceType      string
	RegisterManagerFn func(*mux.Router) any
	Prepare           func(t *testing.T, operation Operation, bridge any)

	SampleJSON        string
	SampleJSONUpdated string
}

type contentType struct {
	name        string
	contentType string
	fromJSON    func(input string) string
	toJSON      func(input string) string
}

type operationTester interface {
	buildRequest(*testing.T, *httptest.Server, contentType, ResourceTypeTest) *http.Request
	assertResponse(*testing.T, *http.Response, contentType, ResourceTypeTest)
	name() Operation
	postAssert(t *testing.T, ct contentType, rt ResourceTypeTest, testServer *httptest.Server)
}

var contentTypes = []contentType{
	{
		name:        "json",
		contentType: "application/json",
		fromJSON:    func(jsonSring string) string { return jsonSring },
		toJSON:      func(jsonSring string) string { return jsonSring },
	},

	{
		name:        "yaml",
		contentType: "text/yaml",
		fromJSON: func(jsonString string) string {
			y, err := yaml.JSONToYAML([]byte(jsonString))
			if err != nil {
				panic(err)
			}
			return string(y)
		},
		toJSON: func(yamlString string) string {
			j, err := yaml.YAMLToJSON([]byte(yamlString))
			if err != nil {
				panic(err)
			}
			return string(j)
		},
	},
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
			t.Run(string(op.name()), func(t *testing.T) {
				op := op
				t.Parallel()

				testOperation(t, op, rt)
			})
		}

	})
}

func testOperation(t *testing.T, op operationTester, rt ResourceTypeTest) {
	t.Helper()

	for _, ct := range contentTypes {
		t.Run(ct.name, func(t *testing.T) {
			ct := ct
			t.Parallel()

			testContentType(t, op, ct, rt)
		})
	}
}

func testContentType(t *testing.T, op operationTester, ct contentType, rt ResourceTypeTest) {
	t.Helper()

	router := mux.NewRouter()
	testServer := httptest.NewServer(router)
	testBridge := rt.RegisterManagerFn(router)

	if rt.Prepare != nil {
		rt.Prepare(t, op.name(), testBridge)
	}

	req := op.buildRequest(t, testServer, ct, rt)
	resp := doRequest(t, req, ct, testServer)

	op.assertResponse(t, resp, ct, rt)
	assert.Equal(t, ct.contentType, resp.Header.Get("Content-Type"))
	op.postAssert(t, ct, rt, testServer)
}

func doRequest(t *testing.T, req *http.Request, ct contentType, testServer *httptest.Server) *http.Response {
	req.Header.Set("Content-Type", ct.contentType)
	resp, err := testServer.Client().Do(req)
	require.NoError(t, err)

	return resp
}
