package configresource_test

import (
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

	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceType: "Config",
		RegisterManagerFn: func(router *mux.Router) any {
			db := testmock.MustCreateRandomMigratedDatabase(db)
			configRepo := configresource.Repository(db, id.GenerateID)

			manager := resourcemanager.New[configresource.Config]("Config", configRepo)
			manager.RegisterRoutes(router)

			return nil
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
