package tests_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/kubeshop/tracetest/server/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type transactionFixture struct {
	t1      model.Test
	t2      model.Test
	testRun model.Run
}

func copyRun(testsDB model.RunRepository, run model.Run) model.Run {
	return createRun(testsDB, model.Test{
		ID:      run.TestID,
		Version: run.TestVersion,
	})
}

func createRun(testsDB model.RunRepository, test model.Test) model.Run {
	run := model.Run{
		State:   model.RunStateFinished,
		TraceID: id.NewRandGenerator().TraceID(),
		SpanID:  id.NewRandGenerator().SpanID(),
	}
	run, err := testsDB.CreateRun(context.TODO(), test, run)
	if err != nil {
		panic(err)
	}

	run.Results.Results = (maps.Ordered[model.SpanQuery, []model.AssertionResult]{}).
		MustAdd("query", []model.AssertionResult{
			{
				Results: []model.SpanAssertionResult{
					{CompareErr: nil},
					{CompareErr: nil},
					{CompareErr: fmt.Errorf("some error")},
				},
			},
		})

	err = testsDB.UpdateRun(context.TODO(), run)
	if err != nil {
		panic(err)
	}
	return run
}

func setupTransactionFixture(t *testing.T, db *sql.DB) transactionFixture {
	testsDB, err := testdb.Postgres(testdb.WithDB(db))
	require.NoError(t, err)

	fixture := transactionFixture{}

	test := model.Test{
		ID:          "ezMn7bE4g",
		Name:        "first test",
		Description: "description",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				URL: "http://localhost:3030/hello-instrumented",
			},
		},
		Specs: (maps.Ordered[model.SpanQuery, model.NamedAssertions]{}).
			MustAdd("query", model.NamedAssertions{
				Name: "some assertion",
				Assertions: []model.Assertion{
					"attr:some_attr = 1",
				},
			}),
		Outputs: (maps.Ordered[string, model.Output]{}).
			MustAdd("output", model.Output{
				Selector: "selector",
				Value:    "value",
			}),
	}
	test, err = testsDB.CreateTest(context.TODO(), test)
	require.NoError(t, err)
	fixture.t1 = test

	fixture.testRun = createRun(testsDB, test)

	test = model.Test{
		ID:          "2qOn7xPVg",
		Name:        "second test",
		Description: "description",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				URL: "http://localhost:3030/hello-instrumented",
			},
		},
		Specs: (maps.Ordered[model.SpanQuery, model.NamedAssertions]{}).
			MustAdd("query", model.NamedAssertions{
				Name: "some assertion",
				Assertions: []model.Assertion{
					"attr:some_attr = 1",
				},
			}),
		Outputs: (maps.Ordered[string, model.Output]{}).
			MustAdd("output", model.Output{
				Selector: "selector",
				Value:    "value",
			}),
	}

	_, err = testsDB.CreateTest(context.TODO(), test)
	require.NoError(t, err)
	fixture.t2 = test

	return fixture
}

func setupTests(t *testing.T, db *sql.DB) model.Run {
	f := setupTransactionFixture(t, db)

	return f.testRun
}

func TestDeleteTestsRelatedToTransactions(t *testing.T) {
	db := testmock.CreateMigratedDatabase()
	defer db.Close()

	testsDB, err := testdb.Postgres(testdb.WithDB(db))
	if err != nil {
		panic(err)
	}
	transactionRepo := tests.NewTransactionsRepository(db, testsDB)

	transactionRepo.Create(context.TODO(), transactionSample)

	f := setupTransactionFixture(t, db)
	createTransactionRun(transactionRepo, transactionSample, f.testRun)

	testsDB.DeleteTest(context.TODO(), f.t1)
	testsDB.DeleteTest(context.TODO(), f.t2)

	actual, err := transactionRepo.Get(context.TODO(), transactionSample.ID)
	assert.NoError(t, err)
	assert.Len(t, actual.StepIDs, 0)

}

var transactionSample = tests.Transaction{
	ID:          "NiWVnxP4R",
	Name:        "Verify Import",
	Description: "check the working of the import flow",
	StepIDs: []id.ID{
		"ezMn7bE4g",
		"2qOn7xPVg",
	},
}

func TestTransactions(t *testing.T) {
	// sample2 := tests.Transaction{
	// 	ID:          "sample2",
	// 	Name:        "Some Transaction",
	// 	Description: "Do important stuff",
	// 	StepIDs: []id.ID{
	// 		"ezMn7bE4g",
	// 	},
	// }

	// sample3 := tests.Transaction{
	// 	ID:          "sample3",
	// 	Name:        "Some Transaction",
	// 	Description: "Do important stuff",
	// 	StepIDs: []id.ID{
	// 		"ezMn7bE4g",
	// 	},
	// }

	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: tests.TransactionResourceName,
		ResourceTypePlural:   tests.TransactionResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			testsDB, err := testdb.Postgres(testdb.WithDB(db))
			if err != nil {
				panic(err)
			}
			transactionsRepo := tests.NewTransactionsRepository(db, testsDB)

			manager := resourcemanager.New[tests.Transaction](
				tests.TransactionResourceName,
				tests.TransactionResourceNamePlural,
				transactionsRepo,
				resourcemanager.CanBeAugmented(),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		Prepare: func(t *testing.T, op rmtests.Operation, manager resourcemanager.Manager) {
			transactionRepo := manager.Handler().(*tests.TransactionsRepository)
			testsDB, err := testdb.Postgres(testdb.WithDB(transactionRepo.DB()))
			if err != nil {
				panic(err)
			}
			switch op {
			case rmtests.OperationGetSuccess,
				rmtests.OperationUpdateSuccess,
				rmtests.OperationListSuccess:
				transactionRepo.Create(context.TODO(), transactionSample)

			case rmtests.OperationDeleteSuccess:
				transactionRepo.Create(context.TODO(), transactionSample)

				// test delete with more than 1 run
				run := setupTests(t, transactionRepo.DB())
				createTransactionRun(transactionRepo, transactionSample, run)

				run = copyRun(testsDB, run)
				createTransactionRun(transactionRepo, transactionSample, run)

			case rmtests.OperationListAugmentedSuccess,
				rmtests.OperationGetAugmentedSuccess:

				transactionRepo.Create(context.TODO(), transactionSample)
				run := setupTests(t, transactionRepo.DB())
				createTransactionRun(transactionRepo, transactionSample, run)

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

func createTransactionRun(repo *tests.TransactionsRepository, tran tests.Transaction, run model.Run) {
	updated, err := repo.GetAugmented(context.TODO(), tran.ID)
	if err != nil {
		panic(err)
	}

	tr, err := repo.CreateRun(context.TODO(), updated.NewRun())
	if err != nil {
		panic(err)
	}
	tr.Steps = []model.Run{run}

	err = repo.UpdateRun(context.TODO(), tr)
	if err != nil {
		panic(err)
	}
}
