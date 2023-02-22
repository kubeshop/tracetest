package resourcemanager_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type sampleResource struct {
	ID   string `mapstructure:"id"`
	Name string `mapstructure:"name"`

	SomeValue string `mapstructure:"some_value"`
}

func (sr sampleResource) Validate() error {
	return nil
}

type sampleResourceManager struct {
	mock.Mock
}

func (m *sampleResourceManager) Create(s sampleResource) (sampleResource, error) {
	args := m.Called(s)
	return args.Get(0).(sampleResource), args.Error(1)
}

func TestSampleResource(t *testing.T) {

	resourceTest := struct {
		resourceType      string
		registerManagerFn func(*mux.Router) *mock.Mock

		sampleNew     any
		sampleNewJSON string

		sampleCreated     any
		sampleCreatedJSON string
	}{
		resourceType: "SampleResource",
		registerManagerFn: func(router *mux.Router) *mock.Mock {
			mockManager := new(sampleResourceManager)
			manager := resourcemanager.New[sampleResource]("SampleResource", mockManager)
			manager.RegisterRoutes(router)

			return &mockManager.Mock
		},

		sampleNew: sampleResource{
			Name:      "test",
			SomeValue: "the value",
		},
		sampleNewJSON: `{
			"type": "SampleResource",
			"spec": {
				"name": "test",
				"some_value": "the value"
			}
		}`,

		sampleCreated: sampleResource{
			ID:        "1",
			Name:      "test",
			SomeValue: "the value",
		},
		sampleCreatedJSON: `{
			"type": "SampleResource",
			"spec": {
				"id": "1",

				"name": "test",
				"some_value": "the value"
			}
		}`,
	}

	t.Run(resourceTest.resourceType, func(t *testing.T) {

		rt := resourceTest
		t.Parallel()

		contentTypes := []struct {
			name     string
			fromJSON func(input string) string
			toJSON   func(input string) string
		}{
			{
				name:     "application/json",
				fromJSON: func(jsonSring string) string { return jsonSring },
				toJSON:   func(jsonSring string) string { return jsonSring },
			},

			{
				name: "text/yaml",
				fromJSON: func(jsonString string) string {
					var parsed map[string]any
					err := json.Unmarshal([]byte(jsonString), &parsed)
					if err != nil {
						panic(err)
					}
					out, err := yaml.Marshal(parsed)
					if err != nil {
						panic(err)
					}
					return string(out)
				},
				toJSON: func(yamlString string) string {
					var parsed map[string]any
					err := yaml.Unmarshal([]byte(yamlString), &parsed)
					if err != nil {
						panic(err)
					}
					out, err := json.Marshal(parsed)
					if err != nil {
						panic(err)
					}
					return string(out)
				},
			},
		}

		for _, ct := range contentTypes {
			t.Run(ct.name, func(t *testing.T) {
				ct := ct
				t.Parallel()

				router := mux.NewRouter()
				testServer := httptest.NewServer(router)
				mockManager := rt.registerManagerFn(router)

				mockManager.
					On("Create", rt.sampleNew).
					Return(rt.sampleCreated, nil)

				input := ct.fromJSON(rt.sampleNewJSON)
				url := fmt.Sprintf("%s/%s/", testServer.URL, strings.ToLower(rt.resourceType))
				req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(input))
				require.NoError(t, err)

				req.Header.Set("Content-Type", ct.name)

				resp, err := testServer.Client().Do(req)
				require.NoError(t, err)

				assert.Equal(t, resp.StatusCode, 201)

				require.NotNil(t, resp.Body)
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				assert.Equal(t, 201, resp.StatusCode)
				assert.JSONEq(t, rt.sampleCreatedJSON, ct.toJSON(string(body)))
			})
		}

	})

}
