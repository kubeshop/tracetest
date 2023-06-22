package config_test

import (
	"database/sql"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/config/configresource"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
)

func TestConfigResource(t *testing.T) {
	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: configresource.ResourceName,
		ResourceTypePlural:   configresource.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			configRepo := configresource.NewRepository(db)

			manager := resourcemanager.New[configresource.Config](
				configresource.ResourceName,
				configresource.ResourceNamePlural,
				configRepo,
				resourcemanager.WithOperations(configresource.Operations...),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		SampleJSON: `{
			"type": "Config",
			"spec": {
				"id": "current",
				"name": "Config",
				"analyticsEnabled": true
			}
		}`,
		SampleJSONUpdated: `{
			"type": "Config",
			"spec": {
				"id": "current",
				"name": "Config",
				"analyticsEnabled": false
			}
		}`,
	},
		rmtests.ExcludeOperations(
			rmtests.OperationGetNotFound,
			rmtests.OperationUpdateNotFound,
			rmtests.OperationListSortSuccess,
			rmtests.OperationListNoResults,
		),
	)
}
