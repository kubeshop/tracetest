package tests_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/kubeshop/tracetest/server/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestWithName(t *testing.T, db model.Repository, name string) model.Test {
	t.Helper()
	test := model.Test{
		Name:        name,
		Description: "description",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				URL: "http://localhost:3030/hello-instrumented",
			},
		},
	}

	updated, err := db.CreateTest(context.TODO(), test)
	if err != nil {
		panic(err)
	}
	return updated
}

func getRepos() (*tests.TransactionsRepository, model.Repository) {
	db := testmock.GetRawTestingDatabase()

	testsRepo, err := testdb.Postgres(testdb.WithDB(db))
	if err != nil {
		panic(err)
	}

	transactionsRepo := tests.NewTransactionsRepository(db, testsRepo.GetTransactionSteps)

	return transactionsRepo, testsRepo
}

func TestCreateTransactionRun(t *testing.T) {
	transactionsRepo, testsRepo := getRepos()

	transaction := tests.Transaction{
		Name:        "first test",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	}

	transaction, err := transactionsRepo.Create(context.TODO(), transaction)
	require.NoError(t, err)

	run := transaction.NewRun()
	newRun, err := transactionsRepo.CreateRun(context.TODO(), run)
	require.NoError(t, err)

	assert.Equal(t, newRun.TransactionID, transaction.ID)
	assert.Equal(t, newRun.State, tests.TransactionRunStateCreated)
}

func TestUpdateTransactionRun(t *testing.T) {
	transactionsRepo, testsRepo := getRepos()

	transaction := tests.Transaction{
		Name:        "first test",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	}

	transaction, err := transactionsRepo.Create(context.TODO(), transaction)
	require.NoError(t, err)

	run := transaction.NewRun()
	run, err = transactionsRepo.CreateRun(context.TODO(), run)
	require.NoError(t, err)

	run.State = tests.TransactionRunStateExecuting
	err = transactionsRepo.UpdateRun(context.TODO(), run)
	require.NoError(t, err)

	updatedRun, err := transactionsRepo.GetTransactionRun(context.TODO(), transaction.ID, run.ID)
	require.NoError(t, err)

	assert.Equal(t, run.TransactionID, transaction.ID)
	assert.Equal(t, tests.TransactionRunStateExecuting, updatedRun.State)
}

func TestDeleteTransactionRun(t *testing.T) {
	transactionsRepo, testsRepo := getRepos()

	transaction := tests.Transaction{
		Name:        "first test",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	}

	transaction, err := transactionsRepo.Create(context.TODO(), transaction)
	require.NoError(t, err)

	run := transaction.NewRun()
	newRun, err := transactionsRepo.CreateRun(context.TODO(), run)
	require.NoError(t, err)

	err = transactionsRepo.DeleteTransactionRun(context.TODO(), newRun)
	require.NoError(t, err)

	_, err = transactionsRepo.GetTransactionRun(context.TODO(), transaction.ID, newRun.ID)
	require.ErrorContains(t, err, "record not found")
}

func TestListTransactionRun(t *testing.T) {
	transactionsRepo, testsRepo := getRepos()

	transaction, err := transactionsRepo.Create(context.TODO(), tests.Transaction{
		Name:        "first test",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	})
	require.NoError(t, err)

	transaction2, err := transactionsRepo.Create(context.TODO(), tests.Transaction{
		Name:        "second transaction",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	})
	require.NoError(t, err)

	run1, err := transactionsRepo.CreateRun(context.TODO(), transaction.NewRun())
	require.NoError(t, err)

	run2, err := transactionsRepo.CreateRun(context.TODO(), transaction.NewRun())
	require.NoError(t, err)

	_, err = transactionsRepo.CreateRun(context.TODO(), transaction2.NewRun())
	require.NoError(t, err)

	runs, err := transactionsRepo.GetTransactionsRuns(context.TODO(), transaction.ID, 20, 0)
	require.NoError(t, err)

	assert.Len(t, runs, 2)
	assert.Equal(t, runs[0].ID, run2.ID)
	assert.Equal(t, runs[1].ID, run1.ID)
}

func TestBug(t *testing.T) {
	transactionsRepo, testsRepo := getRepos()

	ctx := context.TODO()

	transaction := tests.Transaction{
		Name:        "first test",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	}

	transaction, err := transactionsRepo.Create(ctx, transaction)
	require.NoError(t, err)

	run1, err := transactionsRepo.CreateRun(ctx, transaction.NewRun())
	require.NoError(t, err)

	run2, err := transactionsRepo.CreateRun(ctx, transaction.NewRun())
	require.NoError(t, err)

	runs, err := transactionsRepo.GetTransactionsRuns(ctx, transaction.ID, 20, 0)
	require.NoError(t, err)
	assert.Equal(t, runs[0].ID, run2.ID)
	assert.Equal(t, runs[1].ID, run1.ID)

	transaction.Name = "another thing"
	newTransaction, err := transactionsRepo.Update(ctx, transaction)
	require.NoError(t, err)

	run3, err := transactionsRepo.CreateRun(ctx, newTransaction.NewRun())
	require.NoError(t, err)

	runs, err = transactionsRepo.GetTransactionsRuns(ctx, newTransaction.ID, 20, 0)
	require.NoError(t, err)

	assert.Len(t, runs, 3)
	assert.Equal(t, runs[0].ID, run3.ID)
	assert.Equal(t, runs[1].ID, run2.ID)
	assert.Equal(t, runs[2].ID, run1.ID)

}
