package testrunner_test

import (
	"database/sql"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
)

func TestTestRunnerResource(t *testing.T) {
	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: testrunner.ResourceName,
		ResourceTypePlural:   testrunner.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			repo := testrunner.NewRepository(db)

			manager := resourcemanager.New[testrunner.TestRunner](
				testrunner.ResourceName,
				testrunner.ResourceNamePlural,
				repo,
				resourcemanager.DisableDelete(),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		SampleJSON: `{
			"type": "TestRunner",
			"spec": {
				"id": "current",
				"name": "default",
				"requiredGates": [
					"test-specs"
				]
			}
		}`,
		SampleJSONUpdated: `{
			"type": "TestRunner",
			"spec": {
				"id": "current",
				"name": "default",
				"requiredGates": [
					"test-specs"
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
