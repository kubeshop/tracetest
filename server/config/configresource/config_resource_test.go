package configresource_test

import (
	"fmt"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/config/configresource"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/stretchr/testify/mock"
)

type mockableRepo[T resourcemanager.Validator] struct {
	mock.Mock
	repo resourcemanager.ResourceHandler[T]
}

func (m *mockableRepo[T]) Create(in T) (T, error) {
	fmt.Println("ACA", m.ExpectedCalls)
	for _, call := range m.ExpectedCalls {
		if call.Method == "Create" {
			args := m.Called(in)
			return args.Get(0).(T), args.Error(1)
		}
	}

	return m.repo.Create(in)
}

func TestConfigResource(t *testing.T) {

	sample := configresource.Config{
		Name:             "test",
		AnalyticsEnabled: true,
	}

	sampleJSON := `{
		"type": "Config",
		"spec": {
			"id": "1",
			"name": "test",
			"analyticsEnabled": true
		}
	}`
	db := testmock.MustGetRawTestingDatabase()

	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceType: "Config",
		RegisterManagerFn: func(router *mux.Router) any {
			db := testmock.MustCreateRandomMigratedDatabase(db)
			configRepo := configresource.Repository(db, id.GenerateID)

			mockRepo := new(mockableRepo[configresource.Config])
			mockRepo.repo = configRepo

			manager := resourcemanager.New[configresource.Config]("Config", mockRepo)
			manager.RegisterRoutes(router)

			return mockRepo
		},
		Prepare: func(op rmtests.Operation, bridge any) {
			mockRepo := bridge.(*mockableRepo[configresource.Config])
			switch op {
			case rmtests.OperationCreateInteralError:
				mockRepo.On("Create", sample).Return(configresource.Config{}, fmt.Errorf("some error"))
			}
		},
		SampleJSON: sampleJSON,
	})
}
