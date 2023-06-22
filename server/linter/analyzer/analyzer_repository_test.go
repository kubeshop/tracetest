package analyzer_test

import (
	"database/sql"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
)

func TestlinterResource(t *testing.T) {
	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: analyzer.ResourceName,
		ResourceTypePlural:   analyzer.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			repo := analyzer.NewRepository(db)

			manager := resourcemanager.New[analyzer.Linter](
				analyzer.ResourceName,
				analyzer.ResourceNamePlural,
				repo,
				resourcemanager.WithOperations(analyzer.Operations...),
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
