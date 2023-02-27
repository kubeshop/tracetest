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
		Name:             "test",
		AnalyticsEnabled: true,
	}

	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceType: "Config",
		RegisterManagerFn: func(router *mux.Router) any {
			db := testmock.MustCreateRandomMigratedDatabase(db)
			configRepo := configresource.Repository(db, id.GenerateID)

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
			}
		},
		SampleJSON: `{
			"type": "Config",
			"spec": {
				"id": "1",
				"name": "test",
				"analyticsEnabled": true
			}
		}`,

		SampleJSONUpdated: `{
			"type": "Config",
			"spec": {
				"id": "1",
				"name": "test updated",
				"analyticsEnabled": true
			}
		}`,
	})
}
