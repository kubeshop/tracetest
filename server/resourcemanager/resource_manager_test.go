package resourcemanager_test

import (
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

var (
	sampleResourceNew_json = `{
		"type": "SampleResource",
		"spec": {
			"name": "test",
			"some_value": "the value"
		}
	}`
	sampleResourceNew = sampleResource{
		Name:      "test",
		SomeValue: "the value",
	}

	sampleResourceCreated_json = `{
		"type": "SampleResource",
		"spec": {
			"id": "1",

			"name": "test",
			"some_value": "the value"
		}
	}`
	sampleResourceCreated = sampleResource{
		ID:        "1",
		Name:      "test",
		SomeValue: "the value",
	}
)

func TestSampleResource(t *testing.T) {
	router := mux.NewRouter()
	testServer := httptest.NewServer(router)

	mockManager := new(sampleResourceManager)
	manager := resourcemanager.New[sampleResource]("SampleResource", mockManager)
	manager.RegisterRoutes(router)

	mockManager.
		On("Create", sampleResourceNew).
		Return(sampleResourceCreated, nil)

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/sampleresource/", strings.NewReader(sampleResourceNew_json))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)

	resp, err := testServer.Client().Do(req)
	require.NoError(t, err)

	assert.Equal(t, resp.StatusCode, 201)

	body := ""
	if resp.Body != nil {
		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		body = string(b)
	}

	assert.Equal(t, 201, resp.StatusCode)
	assert.JSONEq(t, sampleResourceCreated_json, body)

}
