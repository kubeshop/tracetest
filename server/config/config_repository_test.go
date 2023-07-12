package config_test

import (
	"database/sql"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
)

func TestConfigResource(t *testing.T) {
	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: config.ResourceName,
		ResourceTypePlural:   config.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			configRepo := config.NewRepository(db)

			manager := resourcemanager.New[config.Config](
				config.ResourceName,
				config.ResourceNamePlural,
				configRepo,
				resourcemanager.DisableDelete(),
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
