package resourcemanager_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/stretchr/testify/mock"
)

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

type sampleResourceManager struct {
	mock.Mock
}

func (m *sampleResourceManager) SetID(sr sampleResource, id id.ID) sampleResource {
	sr.ID = id
	return sr
}

func (m *sampleResourceManager) Create(_ context.Context, s sampleResource) (sampleResource, error) {
	args := m.Called(s)
	return args.Get(0).(sampleResource), args.Error(1)
}

func (m *sampleResourceManager) Update(_ context.Context, s sampleResource) (sampleResource, error) {
	args := m.Called(s)
	return args.Get(0).(sampleResource), args.Error(1)
}

func (m *sampleResourceManager) Get(_ context.Context, id id.ID) (sampleResource, error) {
	args := m.Called(id)
	return args.Get(0).(sampleResource), args.Error(1)
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
		ResourceType: "SampleResource",
		RegisterManagerFn: func(router *mux.Router) any {
			mockManager := new(sampleResourceManager)
			manager := resourcemanager.New[sampleResource]("SampleResource", mockManager, func() id.ID {
				return id.ID("3")
			})
			manager.RegisterRoutes(router)

			return mockManager
		},
		Prepare: func(t *testing.T, op rmtests.Operation, bridge any) {
			mockManager := bridge.(*sampleResourceManager)
			mockManager.Test(t)
			switch op {
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
			case rmtests.OperationCreateInteralError:
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
			case rmtests.OperationUpdateInteralError:
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
			case rmtests.OperationGetInteralError:
				mockManager.
					On("Get", sample.ID).
					Return(sampleResource{}, fmt.Errorf("some error"))
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
