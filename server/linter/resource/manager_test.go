package linter_resource_test

import (
	"database/sql"
	"testing"

	"github.com/gorilla/mux"
	linter_resource "github.com/kubeshop/tracetest/server/linter/resource"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
)

func TestlinterResource(t *testing.T) {
	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: linter_resource.ResourceName,
		ResourceTypePlural:   linter_resource.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			repo := linter_resource.NewRepository(db)

			manager := resourcemanager.New[linter_resource.Linter](
				linter_resource.ResourceName,
				linter_resource.ResourceNamePlural,
				repo,
				resourcemanager.WithOperations(linter_resource.Operations...),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		SampleJSON: `{
			"type": "linter",
			"spec": {
				"id": "current",
				"name": "linter",
				"enabled": true,
				"minimumScore": 80,
				"plugins": [
					{
						"name": "standards",
						"enabled": true,
						"required": true
					},
					{
						"name": "security",
						"enabled": true,
						"required": true
					},
					{
						"name": "common",
						"enabled": true,
						"required": true
					}
				]
			}
		}`,
		SampleJSONUpdated: `{
			"type": "linter",
			"spec": {
				"id": "current",
				"name": "linter",
				"enabled": true,
				"minimumScore": 50,
				"plugins": [
					{
						"name": "standards",
						"enabled": false,
						"required": false
					},
					{
						"name": "security",
						"enabled": true,
						"required": true
					},
					{
						"name": "common",
						"enabled": true,
						"required": true
					}
				]
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
