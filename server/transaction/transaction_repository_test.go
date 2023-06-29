package transaction_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/kubeshop/tracetest/server/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type transactionFixture struct {
	t1      test.Test
	t2      test.Test
	testRun test.Run
}

func copyRun(testsDB test.RunRepository, run test.Run) test.Run {
	return createRun(testsDB, test.Test{
		ID:      run.TestID,
		Version: &run.TestVersion,
	})
}

func createRun(runRepository test.RunRepository, t test.Test) test.Run {
	run := test.Run{
		State:   test.RunStateFinished,
		TraceID: id.NewRandGenerator().TraceID(),
		SpanID:  id.NewRandGenerator().SpanID(),
	}
	run, err := runRepository.CreateRun(context.TODO(), t, run)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	return run
}

func setupTransactionFixture(t *testing.T, db *sql.DB) transactionFixture {
	testsDB := test.NewRepository(db)
	runDB := test.NewRunRepository(db)

	fixture := transactionFixture{}

	createdTest := test.Test{
		ID:          "ezMn7bE4g",
		Name:        "first test",
		Description: "description",
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPRequest{
				URL: "http://localhost:3030/hello-instrumented",
			},
		},
		Specs: test.Specs{
			{
				Name:     "some assertion",
				Selector: test.Selector{Query: "query"},
				Assertions: []test.Assertion{
					"attr:some_attr = 1",
				},
			},
		},
		Outputs: []test.Output{
			{Name: "output", Selector: "selector", Value: "value"},
		},
	}
	createdTest, err := testsDB.Create(context.TODO(), createdTest)
	require.NoError(t, err)
	fixture.t1 = createdTest

	fixture.testRun = createRun(runDB, createdTest)

	createdTest = test.Test{
		ID:          "2qOn7xPVg",
		Name:        "second test",
		Description: "description",
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPRequest{
				URL: "http://localhost:3030/hello-instrumented",
			},
		},
		Specs: test.Specs{
			{
				Name:     "some assertion",
				Selector: test.Selector{Query: "query"},
				Assertions: []test.Assertion{
					"attr:some_attr = 1",
				},
			},
		},
		Outputs: []test.Output{
			{
				Name:     "output",
				Selector: "selector",
				Value:    "value",
			},
		},
	}

	_, err = testsDB.Create(context.TODO(), createdTest)
	require.NoError(t, err)
	fixture.t2 = createdTest

	return fixture
}

func setupTests(t *testing.T, db *sql.DB) test.Run {
	f := setupTransactionFixture(t, db)

	return f.testRun
}

func TestDeleteTestsRelatedToTransactions(t *testing.T) {
	db := testmock.CreateMigratedDatabase()
	defer db.Close()

	testRepository := test.NewRepository(db)
	runRepository := test.NewRunRepository(db)
	transactionRepo := transaction.NewRepository(db, testRepository)
	transactionRunRepo := transaction.NewRunRepository(db, runRepository)

	transactionRepo.Create(context.TODO(), transactionSample)

	f := setupTransactionFixture(t, db)
	createTransactionRun(transactionRepo, transactionRunRepo, transactionSample, f.testRun)

	testRepository.Delete(context.TODO(), f.t1)
	testRepository.Delete(context.TODO(), f.t2)

	actual, err := transactionRepo.Get(context.TODO(), transactionSample.ID)
	assert.NoError(t, err)
	assert.Len(t, actual.StepIDs, 0)

}

var transactionSample = transaction.Transaction{
	ID:          "NiWVnxP4R",
	Name:        "Verify Import",
	Description: "check the working of the import flow",
	StepIDs: []id.ID{
		"ezMn7bE4g",
		"2qOn7xPVg",
	},
}

func TestTransactions(t *testing.T) {
	// sample2 := transaction.Transaction{
	// 	ID:          "sample2",
	// 	Name:        "Some Transaction",
	// 	Description: "Do important stuff",
	// 	StepIDs: []id.ID{
	// 		"ezMn7bE4g",
	// 	},
	// }

	// sample3 := transaction.Transaction{
	// 	ID:          "sample3",
	// 	Name:        "Some Transaction",
	// 	Description: "Do important stuff",
	// 	StepIDs: []id.ID{
	// 		"ezMn7bE4g",
	// 	},
	// }

	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: transaction.TransactionResourceName,
		ResourceTypePlural:   transaction.TransactionResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			testsDB := test.NewRepository(db)
			transactionsRepo := transaction.NewRepository(db, testsDB)

			manager := resourcemanager.New[transaction.Transaction](
				transaction.TransactionResourceName,
				transaction.TransactionResourceNamePlural,
				transactionsRepo,
				resourcemanager.CanBeAugmented(),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		Prepare: func(t *testing.T, op rmtests.Operation, manager resourcemanager.Manager) {
			transactionRepo := manager.Handler().(*transaction.Repository)
			runRepository := test.NewRunRepository(transactionRepo.DB())
			runRepo := transaction.NewRunRepository(transactionRepo.DB(), runRepository)

			switch op {
			case rmtests.OperationGetSuccess,
				rmtests.OperationUpdateSuccess,
				rmtests.OperationListSuccess:
				transactionRepo.Create(context.TODO(), transactionSample)

			case rmtests.OperationDeleteSuccess:
				transactionRepo.Create(context.TODO(), transactionSample)

				// test delete with more than 1 run
				run := setupTests(t, transactionRepo.DB())
				createTransactionRun(transactionRepo, runRepo, transactionSample, run)

				run = copyRun(runRepository, run)
				createTransactionRun(transactionRepo, runRepo, transactionSample, run)

			case rmtests.OperationListAugmentedSuccess,
				rmtests.OperationGetAugmentedSuccess:

				transactionRepo.Create(context.TODO(), transactionSample)
				run := setupTests(t, transactionRepo.DB())
				createTransactionRun(transactionRepo, runRepo, transactionSample, run)

				// TODO: reenable this tests when we figure out how to test it
				// problems:
				//   1. sort fields do not map 1:1 with the actual items. Example "last_run" maps to item.summary.last_run.time
				//   2. even if we pass some cusotm function to map fields, non augmented versions do not include the fields so we cannot assert them
				// case rmtests.OperationListSortSuccess:
				// 	transactionRepo.Create(context.TODO(), sample)
				// 	transactionRepo.Create(context.TODO(), sample2)
				// 	transactionRepo.Create(context.TODO(), sample3)
			}
		},
		SampleJSON: `{
			"type": "Transaction",
			"spec": {
				"id": "NiWVnxP4R",
				"name": "Verify Import",
				"description": "check the working of the import flow",
				"steps": [
					"ezMn7bE4g",
					"2qOn7xPVg"
				]
			}
		}`,
		SampleJSONAugmented: `{
			"type": "Transaction",
			"spec": {
				"id": "NiWVnxP4R",
				"createdAt": "REMOVEME",
				"version": 1,
				"name": "Verify Import",
				"description": "check the working of the import flow",
				"steps": [
					"ezMn7bE4g",
					"2qOn7xPVg"
				],
				"fullSteps": [
					{
						"id": "ezMn7bE4g",
						"name": "first test",
						"description": "description",
						"version": 1,
						"createdAt": "REMOVEME",
						"serviceUnderTest": {
							"triggerType": "http",
							"http": {
								"url": "http://localhost:3030/hello-instrumented"
							}
						},
						"specs":[
							{
								"name": "some assertion",
								"selector": "query",
								"assertions": [
									"attr:some_attr = 1"
								]
							}
						],
						"outputs":[
							{
								"name": "output",
								"selector": "selector",
								"value": "value"
							}
						],
						"summary": {
							"runs": 1,
							"lastRun": {
								"time": "REMOVEME",
								"passes": 2,
								"fails": 1,
								"analyzerScore": 0
							}
						}
					},
					{
						"id": "2qOn7xPVg",
						"name": "second test",
						"description": "description",
						"version": 1,
						"createdAt": "REMOVEME",
						"serviceUnderTest": {
							"triggerType": "http",
							"http": {
								"url": "http://localhost:3030/hello-instrumented"
							}
						},
						"specs":[
							{
								"name": "some assertion",
								"selector": "query",
								"assertions": [
									"attr:some_attr = 1"
								]
							}
						],
						"outputs":[
							{
								"name": "output",
								"selector": "selector",
								"value": "value"
							}
						],
						"summary": {
							"runs": 0,
							"lastRun": {
								"fails": 0,
								"passes": 0,
								"time": "REMOVEME",
								"analyzerScore": 0
							}
						}
					}
				],
				"summary": {
					"runs": 1,
					"lastRun": {
						"fails": 1,
						"passes": 2,
						"time": "REMOVEME",
						"analyzerScore": 0
					}
				}
			}
		}`,
		SampleJSONUpdated: `{
			"type": "Transaction",
			"spec": {
				"id": "NiWVnxP4R",
				"name": "Verify Import Updated",
				"description": "check import flow",
				"steps": [
					"ezMn7bE4g"
				]
			}
		}`,
	},
		rmtests.ExcludeOperations(rmtests.OperationListSortSuccess),
		rmtests.JSONComparer(compareJSON),
	)
}

func compareJSON(t require.TestingT, operation rmtests.Operation, firstValue, secondValue string) {
	expected := firstValue
	expected = rmtests.RemoveFieldFromJSONResource("createdAt", expected)
	expected = rmtests.RemoveFieldFromJSONResource("fullSteps.0.createdAt", expected)
	expected = rmtests.RemoveFieldFromJSONResource("fullSteps.1.createdAt", expected)
	expected = rmtests.RemoveFieldFromJSONResource("fullSteps.0.summary.lastRun.time", expected)
	expected = rmtests.RemoveFieldFromJSONResource("fullSteps.1.summary.lastRun.time", expected)
	expected = rmtests.RemoveFieldFromJSONResource("summary.lastRun.time", expected)

	actual := secondValue
	actual = rmtests.RemoveFieldFromJSONResource("createdAt", actual)
	actual = rmtests.RemoveFieldFromJSONResource("fullSteps.0.createdAt", actual)
	actual = rmtests.RemoveFieldFromJSONResource("fullSteps.1.createdAt", actual)
	actual = rmtests.RemoveFieldFromJSONResource("fullSteps.0.summary.lastRun.time", actual)
	actual = rmtests.RemoveFieldFromJSONResource("fullSteps.1.summary.lastRun.time", actual)
	actual = rmtests.RemoveFieldFromJSONResource("summary.lastRun.time", actual)

	require.JSONEq(t, expected, actual)
}

func createTransactionRun(transactionRepo *transaction.Repository, runRepo *transaction.RunRepository, tran transaction.Transaction, run test.Run) {
	updated, err := transactionRepo.GetAugmented(context.TODO(), tran.ID)
	if err != nil {
		panic(err)
	}

	tr, err := runRepo.CreateRun(context.TODO(), updated.NewRun())
	if err != nil {
		panic(err)
	}
	tr.Steps = []test.Run{run}

	err = runRepo.UpdateRun(context.TODO(), tr)
	if err != nil {
		panic(err)
	}
}
