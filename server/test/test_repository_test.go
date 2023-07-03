package test_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtest "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/stretchr/testify/require"
)

var (
	excludedOperations = rmtest.ExcludeOperations(
		// Update
		rmtest.OperationUpdateInternalError,
		rmtest.OperationUpdateNotFound,
		rmtest.OperationUpdateSuccess,

		// Provisioning (Update)
		rmtest.OperationProvisioningTypeNotSupported,
		rmtest.OperationProvisioningError,
		rmtest.OperationProvisioningSuccess,

		// Delete
		rmtest.OperationDeleteInternalError,
		rmtest.OperationDeleteNotFound,
		rmtest.OperationDeleteSuccess,

		// Get
		rmtest.OperationGetAugmentedSuccess,
		rmtest.OperationGetInternalError,
		rmtest.OperationGetNotFound,
		rmtest.OperationGetSuccess,

		// List
		rmtest.OperationListAugmentedSuccess,
		rmtest.OperationListInternalError,
		rmtest.OperationListNoResults,
		rmtest.OperationListSortSuccess,
		rmtest.OperationListSuccess,
		rmtest.OperationListWithInvalidSortField,
	)
	jsonComparer = rmtest.JSONComparer(testJsonComparer)
)

func testJsonComparer(t require.TestingT, operation rmtest.Operation, firstValue, secondValue string) {
	expected := firstValue
	expected = rmtest.RemoveFieldFromJSONResource("createdAt", expected)
	expected = rmtest.RemoveFieldFromJSONResource("specs.0.selector.parsedSelector", expected)
	expected = rmtest.RemoveFieldFromJSONResource("summary", expected)

	actual := secondValue
	actual = rmtest.RemoveFieldFromJSONResource("createdAt", actual)
	actual = rmtest.RemoveFieldFromJSONResource("specs.0.selector.parsedSelector", actual)
	actual = rmtest.RemoveFieldFromJSONResource("summary", actual)

	require.JSONEq(t, expected, actual)
}

func registerManagerFn(router *mux.Router, db *sql.DB) resourcemanager.Manager {
	testRepo := test.NewRepository(db)

	manager := resourcemanager.New[test.Test](
		test.ResourceName,
		test.ResourceNamePlural,
		testRepo,
		resourcemanager.CanBeAugmented(),
	)
	manager.RegisterRoutes(router)

	return manager
}

func getScenarioPreparation(sample test.Test) func(t *testing.T, op rmtest.Operation, manager resourcemanager.Manager) {
	return func(t *testing.T, op rmtest.Operation, manager resourcemanager.Manager) {
		testRepo := manager.Handler().(test.Repository)
		switch op {
		case rmtest.OperationGetSuccess,
			rmtest.OperationUpdateSuccess,
			rmtest.OperationListSuccess:
			testRepo.Create(context.TODO(), sample)

		case rmtest.OperationDeleteSuccess:
			testRepo.Create(context.TODO(), sample)

		case rmtest.OperationListAugmentedSuccess,
			rmtest.OperationGetAugmentedSuccess:
			testRepo.Create(context.TODO(), sample)
		}
	}
}

func Test_TestResourceWithHTTPTrigger(t *testing.T) {
	var testSample = test.Test{
		ID:          "NiWVnxP4R",
		Name:        "Verify Import",
		Description: "check the working of the import flow",
		Trigger: trigger.Trigger{
			Type: "http",
			HTTP: &trigger.HTTPRequest{
				Method: "GET",
				URL:    "http://localhost:11633/api/tests",
			},
		},
		Specs: test.Specs{
			{
				Name:       "check user id exists",
				Selector:   test.Selector{Query: `span[name = "span name"]`},
				Assertions: []test.Assertion{`attr:user_id != ""`},
			},
		},
		Outputs: test.Outputs{
			{
				Name:     "USER_ID",
				Selector: test.SpanQuery(`span[name = "span name"]`),
				Value:    `attr:user_id`,
			},
		},
	}

	testSpec := rmtest.ResourceTypeTest{
		ResourceTypeSingular: test.ResourceName,
		ResourceTypePlural:   test.ResourceNamePlural,
		RegisterManagerFn:    registerManagerFn,
		Prepare:              getScenarioPreparation(testSample),
		SampleJSON: `{
			"type": "test",
			"spec": {
				"id": "NiWVnxP4R",
				"name": "Verify Import",
				"description": "check the working of the import flow",
				"trigger": {
					"type": "http",
					"httpRequest": {
						"method": "GET",
						"url": "http://localhost:11633/api/tests"
					}
				},
				"specs": [
					{
						"name": "check user id exists",
						"selector": { "query": "span[name = \"span name\"]" },
						"assertions": [ "attr:user_id != \"\"" ]
					}
				],
				"outputs": [
					{
						"name": "USER_ID",
						"selector": "span[name = \"span name\"]",
						"value": "attr:user_id"
					}
				],
				"summary": {
					"runs": 1,
					"lastRun": {
						"fails": 1,
						"passes": 2
					}
				}
			}
		}`,
	}

	rmtest.TestResourceType(t, testSpec, excludedOperations, jsonComparer)
}
