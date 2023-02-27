package resourcemanager_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/stretchr/testify/mock"
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

func (m *sampleResourceManager) Create(_ context.Context, s sampleResource) (sampleResource, error) {
	args := m.Called(s)
	return args.Get(0).(sampleResource), args.Error(1)
}

func (m *sampleResourceManager) Update(_ context.Context, s sampleResource) (sampleResource, error) {
	args := m.Called(s)
	return args.Get(0).(sampleResource), args.Error(1)
}

func (m *sampleResourceManager) Get(_ context.Context, id string) (sampleResource, error) {
	args := m.Called(id)
	return args.Get(0).(sampleResource), args.Error(1)
}

func TestSampleResource(t *testing.T) {

	sample := sampleResource{
		Name:      "the name",
		SomeValue: "the value",
	}

	sampleWithID := sampleResource{
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
			manager := resourcemanager.New[sampleResource]("SampleResource", mockManager)
			manager.RegisterRoutes(router)

			return mockManager
		},
		Prepare: func(op rmtests.Operation, bridge any) {
			mockManager := bridge.(*sampleResourceManager)
			switch op {
			// Create
			case rmtests.OperationCreateSuccess:
				mockManager.
					On("Create", sample).
					Return(sampleWithID, nil)
			case rmtests.OperationCreateInteralError:
				mockManager.
					On("Create", sample).
					Return(sampleResource{}, fmt.Errorf("some error"))

			// Update
			case rmtests.OperationUpdateSuccess:
				mockManager.
					On("Update", sampleUpdated).
					Return(sampleUpdated, nil)
			case rmtests.OperationUpdateInteralError:
				mockManager.
					On("Update", sampleUpdated).
					Return(sampleResource{}, fmt.Errorf("some error"))

			// Get
			case rmtests.OperationGetSuccess:
				mockManager.
					On("Get", sampleWithID.ID).
					Return(sampleWithID, nil)
			case rmtests.OperationGetInteralError:
				mockManager.
					On("Get", sampleWithID.ID).
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
