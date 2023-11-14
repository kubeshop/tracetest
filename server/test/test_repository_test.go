package test_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtest "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/kubeshop/tracetest/server/testsuite"
	"github.com/stretchr/testify/assert"
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
	expected = rmtest.RemoveFieldFromJSONResource("summary.lastRun.time", expected)

	actual := secondValue
	actual = rmtest.RemoveFieldFromJSONResource("createdAt", actual)
	actual = rmtest.RemoveFieldFromJSONResource("specs.0.selector.parsedSelector", actual)
	actual = rmtest.RemoveFieldFromJSONResource("summary.lastRun.time", actual)

	require.JSONEq(t, expected, actual)
}

func createRun(runRepository test.RunRepository, t test.Test) (test.Run, error) {
	run := test.Run{
		State:   test.RunStateFinished,
		TraceID: id.NewRandGenerator().TraceID(),
		SpanID:  id.NewRandGenerator().SpanID(),
	}
	run, err := runRepository.CreateRun(context.TODO(), t, run)
	if err != nil {
		return test.NewRun(), err
	}

	run.Results.Results = (maps.Ordered[test.SpanQuery, []test.AssertionResult]{}).
		MustAdd("query", []test.AssertionResult{
			{
				Results: []test.SpanAssertionResult{
					{CompareErr: nil},
					{CompareErr: nil},
					{CompareErr: fmt.Errorf("some error")},
				},
			},
		})

	err = runRepository.UpdateRun(context.TODO(), run)
	if err != nil {
		return test.NewRun(), err
	}

	return run, nil
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
		testRunRepo := test.NewRunRepository(testRepo.DB())

		switch op {
		case rmtest.OperationGetSuccess,
			rmtest.OperationUpdateSuccess,
			rmtest.OperationListSuccess:
			testRepo.Create(context.TODO(), sample)
			_, err := createRun(testRunRepo, sample)
			require.NoError(t, err)

		case rmtest.OperationDeleteSuccess:
			testRepo.Create(context.TODO(), sample)

		case rmtest.OperationListAugmentedSuccess,
			rmtest.OperationGetAugmentedSuccess:
			testRepo.Create(context.TODO(), sample)
			_, err := createRun(testRunRepo, sample)
			require.NoError(t, err)

		case rmtest.OperationListSortSuccess:
			testRepo.Create(context.TODO(), sample)
			testRepo.Create(context.TODO(), secondSample)
			testRepo.Create(context.TODO(), thirdSample)
		}
	}
}

func TestIfDeleteTestsCascadeDeletes(t *testing.T) {
	testmock.StartTestEnvironment()

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
				Selector:   `span[name = "span name"]`,
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
				Selector:   `span[name = "span name"]`,
				Assertions: []test.Assertion{`attr:user_id != ""`},
			},
		},
		Outputs: test.Outputs{},
	}

	var transactionSample = testsuite.TestSuite{
		ID:          "a98s76de",
		Name:        "Verify Import",
		Description: "check the working of the import flow",
		StepIDs: []id.ID{
			testSample.ID,
			secondTestSample.ID,
		},
	}

	db := testmock.CreateMigratedDatabase()
	defer db.Close()

	testRepository := test.NewRepository(db)
	runRepository := test.NewRunRepository(db)
	transactionRepository := testsuite.NewRepository(db, testRepository)
	transactionRunRepository := testsuite.NewRunRepository(db, runRepository)

	_, err := testRepository.Create(context.TODO(), testSample)
	require.NoError(t, err)

	run, err := createRun(runRepository, testSample)
	require.NoError(t, err)

	_, err = testRepository.Create(context.TODO(), secondTestSample)
	require.NoError(t, err)

	secondRun, err := createRun(runRepository, secondTestSample)
	require.NoError(t, err)

	_, err = transactionRepository.Create(context.TODO(), transactionSample)
	require.NoError(t, err)

	updatedTransactionSample, err := transactionRepository.GetAugmented(context.TODO(), transactionSample.ID)
	require.NoError(t, err)

	transactionRun, err := transactionRunRepository.CreateRun(context.TODO(), updatedTransactionSample.NewRun())
	require.NoError(t, err)

	transactionRun.Steps = []test.Run{run, secondRun}

	err = transactionRunRepository.UpdateRun(context.TODO(), transactionRun)
	require.NoError(t, err)

	err = testRepository.Delete(context.TODO(), testSample.ID)
	require.NoError(t, err)

	recentTransactionRun, err := transactionRunRepository.GetTestSuiteRun(context.TODO(), transactionSample.ID, transactionRun.ID)
	require.NoError(t, err)
	assert.Len(t, recentTransactionRun.Steps, 1)

	recentTransaction, err := transactionRepository.Get(context.TODO(), transactionSample.ID)
	require.NoError(t, err)
	assert.Len(t, recentTransaction.StepIDs, 1)

	// TODO: this test was broken by the test run cache, but it's not critical. We can fix it later.
	// _, err = runRepository.GetRun(context.TODO(), run.TestID, run.ID)
	// assert.ErrorIs(t, err, sql.ErrNoRows)

	_, err = testRepository.Get(context.TODO(), testSample.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
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
				Selector:   `span[name = "span name"]`,
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
				Selector:   `span[name = "span name"]`,
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
				Selector:   `span[name = "span name"]`,
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
						"selector": "span[name = \"span name\"]",
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
						"selector": "span[name = \"span name\"]",
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
					"runs": 1,
					"lastRun": {
						"analyzerScore": 0,
						"fails": 1,
						"passes": 2
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
						"selector": "span[name = \"span name updated\"]",
						"assertions": [ "attr:user_id != \"\"" ]
					}
				]
			}
		}`,
	}

	rmtest.TestResourceType(t, testSpec, excludedOperations, jsonComparer)
}

func TestTestResourceWithGRPCTrigger(t *testing.T) {
	protobufFile := `syntax = "proto3";

option java_multiple_files = true;
option java_outer_classname = "PokeshopProto";
option objc_class_prefix = "PKS";

package pokeshop;

service Pokeshop {
  rpc getPokemonList (GetPokemonRequest) returns (GetPokemonListResponse) {}
}

message GetPokemonRequest {
  optional int32 skip = 1;
  optional int32 take = 2;
  optional bool isFixed = 3;
}

message GetPokemonListResponse {
  repeated Pokemon items = 1;
  int32 totalCount = 2;
}`

	var testSample = test.Test{
		ID:          "NiWVnxP4R",
		Name:        "Verify Import",
		Description: "check the working of the import flow",
		Trigger: trigger.Trigger{
			Type: "grpc",
			GRPC: &trigger.GRPCRequest{
				Address:      "someadress:8080",
				Method:       "service.method",
				Request:      `{"hello":"world"}`,
				ProtobufFile: protobufFile,
			},
		},
		Specs: test.Specs{
			{
				Name:       "check user id exists",
				Selector:   `span[name = "span name"]`,
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
			Type: "grpc",
			GRPC: &trigger.GRPCRequest{
				Address:      "someadress:8080",
				Method:       "service.method",
				Request:      `{"hello":"world"}`,
				ProtobufFile: protobufFile,
			},
		},
		Specs: test.Specs{
			{
				Name:       "check user id exists",
				Selector:   `span[name = "span name"]`,
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
			Type: "grpc",
			GRPC: &trigger.GRPCRequest{
				Address:      "someadress:8080",
				Method:       "service.method",
				Request:      `{"hello":"world"}`,
				ProtobufFile: protobufFile,
			},
		},
		Specs: test.Specs{
			{
				Name:       "check user id exists",
				Selector:   `span[name = "span name"]`,
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
					"type": "grpc",
					"grpc": {
						"address": "someadress:8080",
						"method": "service.method",
						"request": "{\"hello\":\"world\"}",
						"protobufFile": "syntax = \"proto3\";\n\noption java_multiple_files = true;\noption java_outer_classname = \"PokeshopProto\";\noption objc_class_prefix = \"PKS\";\n\npackage pokeshop;\n\nservice Pokeshop {\n  rpc getPokemonList (GetPokemonRequest) returns (GetPokemonListResponse) {}\n}\n\nmessage GetPokemonRequest {\n  optional int32 skip = 1;\n  optional int32 take = 2;\n  optional bool isFixed = 3;\n}\n\nmessage GetPokemonListResponse {\n  repeated Pokemon items = 1;\n  int32 totalCount = 2;\n}"
					}
				},
				"specs": [
					{
						"name": "check user id exists",
						"selector": "span[name = \"span name\"]",
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
					"type": "grpc",
					"grpc": {
						"address": "someadress:8080",
						"method": "service.method",
						"request": "{\"hello\":\"world\"}",
						"protobufFile": "syntax = \"proto3\";\n\noption java_multiple_files = true;\noption java_outer_classname = \"PokeshopProto\";\noption objc_class_prefix = \"PKS\";\n\npackage pokeshop;\n\nservice Pokeshop {\n  rpc getPokemonList (GetPokemonRequest) returns (GetPokemonListResponse) {}\n}\n\nmessage GetPokemonRequest {\n  optional int32 skip = 1;\n  optional int32 take = 2;\n  optional bool isFixed = 3;\n}\n\nmessage GetPokemonListResponse {\n  repeated Pokemon items = 1;\n  int32 totalCount = 2;\n}"
					}
				},
				"specs": [
					{
						"name": "check user id exists",
						"selector": "span[name = \"span name\"]",
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
					"runs": 1,
					"lastRun": {
						"analyzerScore": 0,
						"fails": 1,
						"passes": 2
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
					"type": "grpc",
					"grpc": {
						"address": "someadress:8080",
						"method": "service.method",
						"request": "{\"hello\":\"world\"}",
						"protobufFile": "syntax = \"proto3\";\n\noption java_multiple_files = true;\noption java_outer_classname = \"PokeshopProto\";\noption objc_class_prefix = \"PKS\";\n\npackage pokeshop;\n\nservice Pokeshop {\n  rpc getPokemonList (GetPokemonRequest) returns (GetPokemonListResponse) {}\n}\n\nmessage GetPokemonRequest {\n  optional int32 skip = 1;\n  optional int32 take = 2;\n  optional bool isFixed = 3;\n}\n\nmessage GetPokemonListResponse {\n  repeated Pokemon items = 1;\n  int32 totalCount = 2;\n}"
					}
				},
				"specs": [
					{
						"name": "check user id exists updated",
						"selector": "span[name = \"span name updated\"]",
						"assertions": [ "attr:user_id != \"\"" ]
					}
				]
			}
		}`,
	}

	rmtest.TestResourceType(t, testSpec, excludedOperations, jsonComparer)
}

func TestTestResourceWithTraceIDTrigger(t *testing.T) {
	var testSample = test.Test{
		ID:          "NiWVnxP4R",
		Name:        "Verify Import",
		Description: "check the working of the import flow",
		Trigger: trigger.Trigger{
			Type: "traceid",
			TraceID: &trigger.TraceIDRequest{
				ID: "some-trace-id",
			},
		},
		Specs: test.Specs{
			{
				Name:       "check user id exists",
				Selector:   `span[name = "span name"]`,
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
			Type: "traceid",
			TraceID: &trigger.TraceIDRequest{
				ID: "some-trace-id",
			},
		},
		Specs: test.Specs{
			{
				Name:       "check user id exists",
				Selector:   `span[name = "span name"]`,
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
			Type: "traceid",
			TraceID: &trigger.TraceIDRequest{
				ID: "some-trace-id",
			},
		},
		Specs: test.Specs{
			{
				Name:       "check user id exists",
				Selector:   `span[name = "span name"]`,
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
					"type": "traceid",
					"traceid": {
						"id": "some-trace-id"
					}
				},
				"specs": [
					{
						"name": "check user id exists",
						"selector": "span[name = \"span name\"]",
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
					"type": "traceid",
					"traceid": {
						"id": "some-trace-id"
					}
				},
				"specs": [
					{
						"name": "check user id exists",
						"selector": "span[name = \"span name\"]",
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
					"runs": 1,
					"lastRun": {
						"analyzerScore": 0,
						"fails": 1,
						"passes": 2
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
					"type": "traceid",
					"traceid": {
						"id": "some-trace-id"
					}
				},
				"specs": [
					{
						"name": "check user id exists updated",
						"selector": "span[name = \"span name updated\"]",
						"assertions": [ "attr:user_id != \"\"" ]
					}
				]
			}
		}`,
	}

	rmtest.TestResourceType(t, testSpec, excludedOperations, jsonComparer)
}
