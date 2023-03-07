package configresource_test

import (
	"context"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/config/configresource"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/kubeshop/tracetest/server/testmock"
)

func TestConfigResource(t *testing.T) {

	db := testmock.MustGetRawTestingDatabase()
	sampleConfig := configresource.Config{
		ID:               "1",
		Name:             "test 1",
		AnalyticsEnabled: false,
	}
	secondSampleConfig := configresource.Config{
		ID:               "2",
		Name:             "test 2",
		AnalyticsEnabled: true,
	}
	thirdSampleConfig := configresource.Config{
		ID:               "3",
		Name:             "test 3",
		AnalyticsEnabled: false,
	}

	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceType: "Config",
		RegisterManagerFn: func(router *mux.Router) any {
			db := testmock.MustCreateRandomMigratedDatabase(db)
			configRepo := configresource.Repository(db)

			manager := resourcemanager.New[configresource.Config]("Config", configRepo, id.GenerateID)
			manager.RegisterRoutes(router)

			return configRepo
		},
		Prepare: func(t *testing.T, op rmtests.Operation, bridge any) {
			configRepo := bridge.(resourcemanager.ResourceHandler[configresource.Config])
			switch op {
			case rmtests.OperationGetSuccess,
				rmtests.OperationUpdateSuccess,
				rmtests.OperationDeleteSuccess,
				rmtests.OperationListSuccess:
				configRepo.Create(context.TODO(), sampleConfig)
			case rmtests.OperationListPaginatedSuccess:
				configRepo.Create(context.TODO(), sampleConfig)
				configRepo.Create(context.TODO(), secondSampleConfig)
				configRepo.Create(context.TODO(), thirdSampleConfig)
			}
		},
		SampleJSON: `{
			"type": "Config",
			"spec": {
				"id": "1",
				"name": "test 1",
				"analyticsEnabled": false
			}
		}`,
		SampleJSONUpdated: `{
			"type": "Config",
			"spec": {
				"id": "1",
				"name": "test updated",
				"analyticsEnabled": false
			}
		}`,
		SortField:        "id",
		InvalidSortField: "invalid_field",
	})
}
