package testutil

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

type Operation string

const (
	OperationCreateSuccess Operation = "CreateSuccess"
)

type ResourceTypeTest struct {
	ResourceType      string
	RegisterManagerFn func(*mux.Router) any
	Prepare           func(operation Operation, bridge any)

	SampleNew     any
	SampleNewJSON string

	SampleCreated     any
	SampleCreatedJSON string
}

var contentTypes = []struct {
	name        string
	contentType string
	fromJSON    func(input string) string
	toJSON      func(input string) string
}{
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

func TestResourceType(t *testing.T, rt *ResourceTypeTest) {
	t.Helper()

	t.Run(rt.ResourceType, func(t *testing.T) {

		rt := rt
		t.Parallel()

		for _, ct := range contentTypes {
			t.Run(ct.name, func(t *testing.T) {
				ct := ct
				t.Parallel()

				router := mux.NewRouter()
				testServer := httptest.NewServer(router)
				testBridge := rt.RegisterManagerFn(router)

				if rt.Prepare != nil {
					rt.Prepare(OperationCreateSuccess, testBridge)
				}

				input := ct.fromJSON(rt.SampleNewJSON)
				url := fmt.Sprintf("%s/%s/", testServer.URL, strings.ToLower(rt.ResourceType))
				req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(input))
				require.NoError(t, err)

				req.Header.Set("Content-Type", ct.contentType)

				resp, err := testServer.Client().Do(req)
				require.NoError(t, err)

				assert.Equal(t, resp.StatusCode, 201)

				require.NotNil(t, resp.Body)
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				assert.Equal(t, 201, resp.StatusCode)
				assert.JSONEq(t, rt.SampleCreatedJSON, ct.toJSON(string(body)))
			})
		}

	})
}
