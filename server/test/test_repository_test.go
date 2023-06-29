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
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/stretchr/testify/require"
)

var (
	excludedOperations = rmtest.ExcludeOperations(
		// List
		rmtest.OperationListSortSuccess, // we need to think how to deal with augmented fields on sorting
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

func getScenarioPreparation(sample, secondSample, thirdSample test.Test) func(t *testing.T, op rmtest.Operation, manager resourcemanager.Manager) {
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

		case rmtest.OperationListSortSuccess:
			testRepo.Create(context.TODO(), sample)
			testRepo.Create(context.TODO(), secondSample)
			testRepo.Create(context.TODO(), thirdSample)
		}
	}
}

func TestTestResourceWithHTTPTrigger(t *testing.T) {
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

	var secondTestSample = test.Test{
		ID:          "NiWVnjahsdvR",
		Name:        "Another Test",
		Description: "another test description",
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
		Outputs: test.Outputs{},
	}

	var thirdTestSample = test.Test{
		ID:          "oau3si2y6d",
		Name:        "One More Test",
		Description: "one more test description",
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
		Outputs: test.Outputs{},
	}

	testSpec := rmtest.ResourceTypeTest{
		ResourceTypeSingular: test.ResourceName,
		ResourceTypePlural:   test.ResourceNamePlural,
		RegisterManagerFn:    registerManagerFn,
		Prepare:              getScenarioPreparation(testSample, secondTestSample, thirdTestSample),
		SampleJSON: `{
			"type": "Test",
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
				]
			}
		}`,
		SampleJSONAugmented: `{
			"type": "Test",
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
				"version": 1,
				"summary": {
					"runs": 0,
					"lastRun": {
						"fails": 0,
						"passes": 0
					}
				}
			}
		}`,
		SampleJSONUpdated: `{
			"type": "Test",
			"spec": {
				"id": "NiWVnxP4R",
				"name": "Verify Import Updated",
				"description": "check the working of the import flow updated",
				"trigger": {
					"type": "http",
					"httpRequest": {
						"method": "GET",
						"url": "http://localhost:11633/api/tests"
					}
				},
				"specs": [
					{
						"name": "check user id exists updated",
						"selector": { "query": "span[name = \"span name updated\"]" },
						"assertions": [ "attr:user_id != \"\"" ]
					}
				]
			}
		}`,
	}

	testmock.StartTestEnvironment() //TODO: remove it later
	rmtest.TestResourceType(t, testSpec, excludedOperations, jsonComparer)
}
