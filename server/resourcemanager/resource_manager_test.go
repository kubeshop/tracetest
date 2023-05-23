package resourcemanager_test

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/pkg/id"
	rm "github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDifferentInputAndOuptutEncodings(t *testing.T) {
	// make sure to preserve indentation and spaces instead of tabs.
	// yaml is very sensitive to that, and fails silently
	inputYaml := `
type: Resource
spec:
  name: the name
  some_value: the value
		`

	expectedJSON := `{
		"type": "Resource",
		"spec": {
			"id": "1",
			"name": "the name",
			"some_value": "the value"
		}
	}`

	mockSample := sampleResource{
		ID:        "1",
		Name:      "the name",
		SomeValue: "the value",
	}

	//setup
	router := mux.NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	mockManager := new(sampleResourceManager)
	mockManager.On("Create", mockSample).Return(mockSample, nil)
	manager := rm.New[sampleResource](
		"Resource",
		"Resources",
		mockManager,
		rm.WithIDGen(func() id.ID {
			return id.ID("1")
		}),
		rm.WithOperations(rm.OperationCreate),
	)
	manager.RegisterRoutes(router)

	// prepare request
	req, err := http.NewRequest(
		http.MethodPost,
		testServer.URL+"/resources",
		strings.NewReader(inputYaml),
	)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "text/yaml")  // request content-type
	req.Header.Set("Accept", "application/json") // expected response content-type

	resp, err := testServer.Client().Do(req)
	require.NoError(t, err)

	actualBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	require.NoError(t, err)

	assert.JSONEq(t, expectedJSON, string(actualBody))
}

func TestSampleResource(t *testing.T) {
	sample := sampleResource{
		ID:        "1",
		Name:      "the name",
		SomeValue: "the value",
	}

	sampleUpdated := sampleResource{
		ID:        "1",
		Name:      "the name updated",
		SomeValue: "the value updated",
	}

	rmtests.TestResourceTypeWithErrorOperations(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: "SampleResource",
		ResourceTypePlural:   "SampleResources",
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) rm.Manager {
			mockManager := new(sampleResourceManager)
			manager := rm.New[sampleResource](
				"SampleResource",
				"SampleResources",
				mockManager,
				rm.WithIDGen(func() id.ID {
					return id.ID("3")
				}),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		Prepare: func(t *testing.T, op rmtests.Operation, manager rm.Manager) {
			mockManager := manager.Handler().(*sampleResourceManager)
			mockManager.Test(t)

			switch op {
			// Provisioning
			case rmtests.OperationProvisioningSuccess:
				mockManager.
					On("Provision", sample).
					Return(nil)

				// Create
			case rmtests.OperationCreateNoID:
				withGenID := sample
				withGenID.ID = id.ID("3")
				mockManager.
					On("Create", withGenID).
					Return(sample, nil)
			case rmtests.OperationCreateSuccess:
				mockManager.
					On("Create", sample).
					Return(sample, nil)
			case rmtests.OperationCreateInternalError:
				mockManager.
					On("Create", sample).
					Return(sampleResource{}, fmt.Errorf("some error"))

			// Update
			case rmtests.OperationUpdateNotFound:
				mockManager.
					On("Update", sampleUpdated).
					Return(sampleResource{}, sql.ErrNoRows)
			case rmtests.OperationUpdateSuccess:
				mockManager.
					On("Update", sampleUpdated).
					Return(sampleUpdated, nil)
			case rmtests.OperationUpdateInternalError:
				mockManager.
					On("Update", sampleUpdated).
					Return(sampleResource{}, fmt.Errorf("some error"))

			// Get
			case rmtests.OperationGetNotFound:
				mockManager.
					On("Get", sample.ID).
					Return(sampleResource{}, sql.ErrNoRows)
			case rmtests.OperationGetSuccess:
				mockManager.
					On("Get", sample.ID).
					Return(sample, nil)
			case rmtests.OperationGetInternalError:
				mockManager.
					On("Get", sample.ID).
					Return(sampleResource{}, fmt.Errorf("some error"))

			// Delete
			case rmtests.OperationDeleteNotFound:
				mockManager.
					On("Delete", sample.ID).
					Return(sql.ErrNoRows)
			case rmtests.OperationDeleteSuccess:
				mockManager.
					On("Delete", sample.ID).
					Return(nil)
				mockManager.
					On("Get", sample.ID).
					Return(sampleResource{}, sql.ErrNoRows)
			case rmtests.OperationDeleteInternalError:
				mockManager.
					On("Delete", sample.ID).
					Return(fmt.Errorf("some error"))

			// List
			case rmtests.OperationListSuccess:
				mockManager.
					On("Count", mock.Anything).
					Return(1, nil)
				mockManager.
					On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return([]sampleResource{sample}, nil)
			case rmtests.OperationListNoResults:
				mockManager.
					On("Count", mock.Anything).
					Return(0, nil)
				mockManager.
					On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return([]sampleResource{}, nil)
			case rmtests.OperationListSortSuccess:
				mockManager.
					On("Count", mock.Anything).
					Return(3, nil)

				prepareSortByID(mockManager)
				prepareSortByName(mockManager)
				prepareSortBySomeValue(mockManager)
			case rmtests.OperationListInternalError:
				mockManager.
					On("Count", mock.Anything).
					Return(0, fmt.Errorf("some error"))
			}
		},
		SampleJSON: `{
			"type": "SampleResource",
			"spec": {
				"id": "1",
				"name": "the name",
				"some_value": "the value"
			}
		}`,
		SampleJSONUpdated: `{
			"type": "SampleResource",
			"spec": {
				"id": "1",
				"name": "the name updated",
				"some_value": "the value updated"
			}
		}`,
	})
}

func TestRestrictedResource(t *testing.T) {
	sample := sampleResource{
		ID:        "1",
		Name:      "the name",
		SomeValue: "the value",
	}

	sampleUpdated := sampleResource{
		ID:        "1",
		Name:      "the name updated",
		SomeValue: "the value updated",
	}

	rmtests.TestResourceTypeWithErrorOperations(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: "RestrictedResource",
		ResourceTypePlural:   "RestrictedResources",
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) rm.Manager {
			mockManager := new(restrictedResourceManager)
			manager := rm.New[sampleResource](
				"RestrictedResource",
				"RestrictedResources",
				mockManager,
				rm.WithIDGen(func() id.ID {
					return id.ID("3")
				}),
				rm.WithOperations(mockManager.Operations()...),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		Prepare: func(t *testing.T, op rmtests.Operation, manager rm.Manager) {
			mockManager := manager.Handler().(*restrictedResourceManager)

			switch op {
			// Provisioning
			case rmtests.OperationProvisioningSuccess:
				mockManager.
					On("Provision", sample).
					Return(nil)
			// Update
			case rmtests.OperationUpdateNotFound:
				mockManager.
					On("Update", sampleUpdated).
					Return(sampleResource{}, sql.ErrNoRows)
			case rmtests.OperationUpdateSuccess:
				mockManager.
					On("Update", sampleUpdated).
					Return(sampleUpdated, nil)
			case rmtests.OperationUpdateInternalError:
				mockManager.
					On("Update", sampleUpdated).
					Return(sampleResource{}, fmt.Errorf("some error"))

				// Get
			case rmtests.OperationGetNotFound:
				mockManager.
					On("Get", sample.ID).
					Return(sampleResource{}, sql.ErrNoRows)
			case rmtests.OperationGetSuccess:
				mockManager.
					On("Get", sample.ID).
					Return(sample, nil)
			case rmtests.OperationGetInternalError:
				mockManager.
					On("Get", sample.ID).
					Return(sampleResource{}, fmt.Errorf("some error"))
			}
		},
		SampleJSON: `{
			"type": "RestrictedResource",
			"spec": {
				"id": "1",
				"name": "the name",
				"some_value": "the value"
			}
		}`,
		SampleJSONUpdated: `{
			"type": "RestrictedResource",
			"spec": {
				"id": "1",
				"name": "the name updated",
				"some_value": "the value updated"
			}
		}`,
	})
}

func TestAugmentedResource(t *testing.T) {
	sample := sampleResource{
		ID:        "1",
		Name:      "the name",
		SomeValue: "the value",
	}

	sampleAugmented := sampleResource{
		ID:                     "1",
		Name:                   "the name",
		SomeValue:              "the value",
		SomeAugmentedOnlyValue: "augmentation works",
	}

	rmtests.TestResourceTypeWithErrorOperations(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: "AugmentedResource",
		ResourceTypePlural:   "AugmentedResources",
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) rm.Manager {
			mockManager := new(augmentedResourceManager)
			manager := rm.New[sampleResource](
				"AugmentedResource",
				"AugmentedResources",
				mockManager,
				rm.WithOperations(mockManager.Operations()...),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		Prepare: func(t *testing.T, op rmtests.Operation, manager rm.Manager) {
			mockManager := manager.Handler().(*augmentedResourceManager)
			mockManager.Test(t)

			switch op {
			// Provisioning
			case rmtests.OperationProvisioningSuccess:
				mockManager.
					On("Provision", sample).
					Return(nil)

			// Get
			case rmtests.OperationGetNotFound:
				mockManager.
					On("Get", sample.ID).
					Return(sampleResource{}, sql.ErrNoRows)
			case rmtests.OperationGetSuccess:
				mockManager.
					On("Get", sample.ID).
					Return(sample, nil)
			case rmtests.OperationGetInternalError:
				mockManager.
					On("Get", sample.ID).
					Return(sampleResource{}, fmt.Errorf("some error"))

			// List
			case rmtests.OperationListSuccess:
				mockManager.
					On("Count", mock.Anything).
					Return(1, nil)
				mockManager.
					On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return([]sampleResource{sample}, nil)
			case rmtests.OperationListNoResults:
				mockManager.
					On("Count", mock.Anything).
					Return(0, nil)
				mockManager.
					On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return([]sampleResource{}, nil)
			case rmtests.OperationListSortSuccess:
				mockManager.
					On("Count", mock.Anything).
					Return(3, nil)

				prepareSortByID(mockManager)
				prepareSortByName(mockManager)
				prepareSortBySomeValue(mockManager)
			case rmtests.OperationListInternalError:
				mockManager.
					On("Count", mock.Anything).
					Return(0, fmt.Errorf("some error"))

			// Augmented
			case rmtests.OperationGetAugmentedSuccess:
				mockManager.
					On("GetAugmented", sampleAugmented.ID).
					Return(sampleAugmented, nil)
			case rmtests.OperationListAugmentedSuccess:
				mockManager.
					On("Count", mock.Anything).
					Return(1, nil)
				mockManager.
					On("ListAugmented", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return([]sampleResource{sampleAugmented}, nil)
			}
		},
		SampleJSON: `{
			"type": "AugmentedResource",
			"spec": {
				"id": "1",
				"name": "the name",
				"some_value": "the value"
			}
		}`,
		SampleJSONAugmented: `{
			"type": "AugmentedResource",
			"spec": {
				"id": "1",
				"name": "the name",
				"some_value": "the value",
				"some_augmented_value": "augmentation works"
			}
		}`,
	})
}

// test structures and mocks

type sampleResource struct {
	ID   id.ID  `json:"id"`
	Name string `json:"name"`

	SomeValue              string `json:"some_value"`
	SomeAugmentedOnlyValue string `json:"some_augmented_value,omitempty"`
}

func (sr sampleResource) HasID() bool {
	return sr.ID.String() != ""
}

func (sr sampleResource) GetID() id.ID {
	return sr.ID
}

func (sr sampleResource) Validate() error {
	return nil
}

type baseResourceManager struct {
	mock.Mock
}

func (m *baseResourceManager) Get(_ context.Context, id id.ID) (sampleResource, error) {
	args := m.Called(id)
	return args.Get(0).(sampleResource), args.Error(1)
}

func (m *baseResourceManager) Update(_ context.Context, s sampleResource) (sampleResource, error) {
	args := m.Called(s)
	return args.Get(0).(sampleResource), args.Error(1)
}

func (m *baseResourceManager) SetID(sr sampleResource, id id.ID) sampleResource {
	sr.ID = id
	return sr
}

func (m *baseResourceManager) Provision(_ context.Context, s sampleResource) error {
	args := m.Called(s)
	return args.Error(0)
}

type restrictedResourceManager struct {
	baseResourceManager
}

func (m *restrictedResourceManager) Operations() []rm.Operation {
	return []rm.Operation{
		rm.OperationGet,
		rm.OperationUpdate,
	}
}

type augmentedResourceManager struct {
	sampleResourceManager
}

func (m *augmentedResourceManager) Operations() []rm.Operation {
	return []rm.Operation{
		rm.OperationGet,
		rm.OperationGetAugmented,
		rm.OperationList,
		rm.OperationListAugmented,
	}
}

func (m *augmentedResourceManager) GetAugmented(_ context.Context, id id.ID) (sampleResource, error) {
	args := m.Called(id)
	return args.Get(0).(sampleResource), args.Error(1)
}

func (m *augmentedResourceManager) ListAugmented(_ context.Context, take, skip int, query, sortBy, sortDirection string) ([]sampleResource, error) {
	args := m.Called(take, skip, query, sortBy, sortDirection)
	return args.Get(0).([]sampleResource), args.Error(1)
}

type sampleResourceManager struct {
	baseResourceManager
}

func (m *sampleResourceManager) Create(_ context.Context, s sampleResource) (sampleResource, error) {
	args := m.Called(s)
	return args.Get(0).(sampleResource), args.Error(1)
}

func (m *sampleResourceManager) Delete(_ context.Context, id id.ID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *sampleResourceManager) SortingFields() []string {
	return []string{"id", "name", "some_value"}
}

func (m *sampleResourceManager) List(_ context.Context, take, skip int, query, sortBy, sortDirection string) ([]sampleResource, error) {
	args := m.Called(take, skip, query, sortBy, sortDirection)
	return args.Get(0).([]sampleResource), args.Error(1)
}

func (m *sampleResourceManager) Count(_ context.Context, query string) (int, error) {
	args := m.Called(query)
	return args.Int(0), args.Error(1)
}

type mockable interface {
	On(string, ...interface{}) *mock.Call
}

func prepareSortByID(m mockable) {
	m.On("List", mock.Anything, mock.Anything, mock.Anything, "id", "asc").
		Return([]sampleResource{
			{ID: "1", Name: "3", SomeValue: "3"},
			{ID: "2", Name: "1", SomeValue: "1"},
			{ID: "3", Name: "2", SomeValue: "2"},
		}, nil)
	m.On("List", mock.Anything, mock.Anything, mock.Anything, "id", "desc").
		Return([]sampleResource{
			{ID: "3", Name: "2", SomeValue: "2"},
			{ID: "2", Name: "1", SomeValue: "1"},
			{ID: "1", Name: "3", SomeValue: "3"},
		}, nil)
}

func prepareSortByName(m mockable) {
	m.On("List", mock.Anything, mock.Anything, mock.Anything, "name", "asc").
		Return([]sampleResource{
			{Name: "1", ID: "3", SomeValue: "3"},
			{Name: "2", ID: "1", SomeValue: "1"},
			{Name: "3", ID: "2", SomeValue: "2"},
		}, nil)
	m.On("List", mock.Anything, mock.Anything, mock.Anything, "name", "desc").
		Return([]sampleResource{
			{Name: "3", ID: "2", SomeValue: "2"},
			{Name: "2", ID: "1", SomeValue: "1"},
			{Name: "1", ID: "3", SomeValue: "3"},
		}, nil)
}

func prepareSortBySomeValue(m mockable) {
	m.On("List", mock.Anything, mock.Anything, mock.Anything, "some_value", "asc").
		Return([]sampleResource{
			{SomeValue: "1", ID: "3", Name: "3"},
			{SomeValue: "2", ID: "1", Name: "1"},
			{SomeValue: "3", ID: "2", Name: "2"},
		}, nil)
	m.On("List", mock.Anything, mock.Anything, mock.Anything, "some_value", "desc").
		Return([]sampleResource{
			{SomeValue: "3", ID: "2", Name: "2"},
			{SomeValue: "2", ID: "1", Name: "1"},
			{SomeValue: "1", ID: "3", Name: "3"},
		}, nil)
}
