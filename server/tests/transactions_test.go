package tests_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/tests"
	"github.com/stretchr/testify/require"
)

func setupTests(t *testing.T, db *sql.DB) model.Run {
	testsDB, err := testdb.Postgres(testdb.WithDB(db))
	require.NoError(t, err)

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
	}
	test, err = testsDB.CreateTest(context.TODO(), test)
	require.NoError(t, err)

	run := model.Run{
		State:   model.RunStateFinished,
		TraceID: id.NewRandGenerator().TraceID(),
		SpanID:  id.NewRandGenerator().SpanID(),
	}
	run, err = testsDB.CreateRun(context.TODO(), test, run)
	require.NoError(t, err)

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
	require.NoError(t, err)

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
	}

	_, err = testsDB.CreateTest(context.TODO(), test)
	require.NoError(t, err)

	return run

}

func TestTransactions(t *testing.T) {
	d := time.Date(2022, 01, 01, 14, 54, 0, 0, time.UTC)

	sample := tests.Transaction{
		ID:          "NiWVnxP4R",
		CreatedAt:   &d,
		Name:        "Verify Import",
		Description: "check the working of the import flow",
		StepIDs: []id.ID{
			"ezMn7bE4g",
			"2qOn7xPVg",
		},
	}

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
			transactionsRepo := tests.NewTransactionsRepository(db, testsDB.GetTransactionSteps)

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
			switch op {
			case rmtests.OperationGetSuccess,
				rmtests.OperationUpdateSuccess,
				rmtests.OperationDeleteSuccess,
				rmtests.OperationListSuccess:
				transactionRepo.Create(context.TODO(), sample)

			case rmtests.OperationListAugmentedSuccess,
				rmtests.OperationGetAugmentedSuccess:

				transactionRepo.Create(context.TODO(), sample)
				// create a local copy of sample, with all the data
				sample, err := transactionRepo.GetAugmented(context.TODO(), sample.ID)
				if err != nil {
					panic(err)
				}

				run := setupTests(t, transactionRepo.DB())
				transactionRepo.Create(context.TODO(), sample)

				tr, err := transactionRepo.CreateRun(context.TODO(), sample.NewRun())
				if err != nil {
					panic(err)
				}

				tr.CompletedAt = d
				tr.CreatedAt = d
				tr.CurrentTest = 1
				tr.Steps = []model.Run{run}

				err = transactionRepo.UpdateRun(context.TODO(), tr)
				if err != nil {
					panic(err)
				}

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
				"createdAt": "` + d.Format(time.RFC3339) + `",
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
						"specs":[],
						"outputs":[],
						"summary": {
							"runs": 1,
							"lastRun": {
								"time": "` + d.Format(time.RFC3339) + `",
								"passes": 2,
								"fails": 1
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
						"specs":[],
						"outputs":[],
						"summary": {
							"runs": 0
						}
					}
				],
				"summary": {
					"runs": 1,
					"lastRun": {
						"fails": 1,
						"passes": 2,
						"time": "` + d.Format(time.RFC3339) + `"
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
	expected = rmtests.RemoveFieldFromJSONResource("fullSteps.0.createdAt", expected)
	expected = rmtests.RemoveFieldFromJSONResource("fullSteps.1.createdAt", expected)

	actual := secondValue
	actual = rmtests.RemoveFieldFromJSONResource("fullSteps.0.createdAt", actual)
	actual = rmtests.RemoveFieldFromJSONResource("fullSteps.1.createdAt", actual)

	require.JSONEq(t, expected, actual)
}
