package resourcemanager_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/id"
	rm "github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/stretchr/testify/mock"
)

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

	prepareSortByID := func(m *sampleResourceManager) {
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

	prepareSortByName := func(m *sampleResourceManager) {
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

	prepareSortBySomeValue := func(m *sampleResourceManager) {
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

	rmtests.TestResourceTypeWithErrorOperations(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: "SampleResource",
		ResourceTypePlural:   "SampleResources",
		RegisterManagerFn: func(router *mux.Router) rm.Manager {
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
			case rmtests.OperationListPaginatedSuccess:
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
		RegisterManagerFn: func(router *mux.Router) rm.Manager {
			mockManager := new(restrictedResourceManager)
			manager := rm.New[sampleResource](
				"RestrictedResource",
				"RestrictedResources",
				mockManager,
				rm.WithIDGen(func() id.ID {
					return id.ID("3")
				}),
				rm.WithOperations(rm.OperationGet, rm.OperationUpdate),
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

// test structures and mocks

type sampleResource struct {
	ID   id.ID  `mapstructure:"id"`
	Name string `mapstructure:"name"`

	SomeValue string `mapstructure:"some_value"`
}

func (sr sampleResource) HasID() bool {
	return sr.ID.String() != ""
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
