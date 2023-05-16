package lintern_resource_test

import (
	"database/sql"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/config/configresource"
	lintern_resource "github.com/kubeshop/tracetest/server/lintern/resource"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
)

func TestConfigResource(t *testing.T) {
	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: lintern_resource.ResourceName,
		ResourceTypePlural:   lintern_resource.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			repo := lintern_resource.NewRepository(db)

			manager := resourcemanager.New[configresource.Config](
				lintern_resource.ResourceName,
				lintern_resource.ResourceNamePlural,
				repo,
				resourcemanager.WithOperations(lintern_resource.Operations...),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		SampleJSON: `{
			"type": "Lintern",
			"spec": {
				"id": "current",
				"name": "Lintern",
				"enabled": true,
				"minimiumScore": 80,
				"plugins": [{
					name: "standards",
					enabled: true,
					required: true
				}, {
					name: "security",
					enabled: true, 
					required: true
				}]
			}
		}`,
		SampleJSONUpdated: `{
			"type": "Lintern",
			"spec": {
				"id": "current",
				"name": "Lintern",
				"enabled": true,
				"minimiumScore": 50,
				"plugins": [{
					name: "standards",
					enabled: false,
					required: false
				}, {
					name: "security",
					enabled: true, 
					required: true
				}]
			}
		}`,
	},
		rmtests.ExcludeOperations(
			rmtests.OperationGetNotFound,
			rmtests.OperationUpdateNotFound,
			rmtests.OperationListPaginatedSuccess,
			rmtests.OperationListNoResults,
		),
	)
}
