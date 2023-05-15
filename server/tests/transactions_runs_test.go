package tests_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/kubeshop/tracetest/server/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestWithName(t *testing.T, db model.TestRepository, name string) model.Test {
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

	transactionRepo := tests.NewTransactionsRepository(db, testsRepo.GetTransactionSteps)

	return transactionRepo, testsRepo
}

func getTransaction(t *testing.T, transactionRepo *tests.TransactionsRepository, testsRepo model.TestRepository) (tests.Transaction, transactionFixture) {
	f := setupTransactionFixture(t, transactionRepo.DB())

	transaction := tests.Transaction{
		ID:          id.NewRandGenerator().ID(),
		Name:        "first test",
		Description: "description",
		StepIDs: []id.ID{
			f.t1.ID,
			f.t2.ID,
		},
	}

	_, err := transactionRepo.Create(context.TODO(), transaction)
	require.NoError(t, err)

	transaction, err = transactionRepo.GetAugmented(context.TODO(), transaction.ID)
	require.NoError(t, err)

	return transaction, f
}

func TestCreateTransactionRun(t *testing.T) {
	transactionRepo, testsRepo := getRepos()
	transaction, _ := getTransaction(t, transactionRepo, testsRepo)

	tr, err := transactionRepo.CreateRun(context.TODO(), transaction.NewRun())
	require.NoError(t, err)

	assert.Equal(t, tr.TransactionID, transaction.ID)
	assert.Equal(t, tr.State, tests.TransactionRunStateCreated)
	assert.Len(t, tr.Steps, 0)
}

func TestUpdateTransactionRun(t *testing.T) {
	transactionRepo, testsRepo := getRepos()
	transaction, fixture := getTransaction(t, transactionRepo, testsRepo)

	tr, err := transactionRepo.CreateRun(context.TODO(), transaction.NewRun())
	require.NoError(t, err)

	tr.State = tests.TransactionRunStateExecuting
	tr.Steps = []model.Run{fixture.testRun}
	err = transactionRepo.UpdateRun(context.TODO(), tr)
	require.NoError(t, err)

	updatedRun, err := transactionRepo.GetTransactionRun(context.TODO(), transaction.ID, tr.ID)
	require.NoError(t, err)

	assert.Equal(t, tr.TransactionID, transaction.ID)
	assert.Equal(t, tests.TransactionRunStateExecuting, updatedRun.State)
	assert.Len(t, tr.Steps, 1)
}

func TestDeleteTransactionRun(t *testing.T) {
	transactionRepo, testsRepo := getRepos()
	transaction, _ := getTransaction(t, transactionRepo, testsRepo)

	tr, err := transactionRepo.CreateRun(context.TODO(), transaction.NewRun())
	require.NoError(t, err)

	err = transactionRepo.DeleteTransactionRun(context.TODO(), tr)
	require.NoError(t, err)

	_, err = transactionRepo.GetTransactionRun(context.TODO(), transaction.ID, tr.ID)
	require.ErrorContains(t, err, "no rows in result set")
}

func createTransaction(t *testing.T, repo *tests.TransactionsRepository, tran tests.Transaction) tests.Transaction {
	one := 1
	tran.ID = id.GenerateID()
	tran.Version = &one
	for _, step := range tran.Steps {
		tran.StepIDs = append(tran.StepIDs, step.ID)
	}
	_, err := repo.Create(context.TODO(), tran)
	require.NoError(t, err)

	tran, err = repo.GetAugmented(context.TODO(), tran.ID)
	require.NoError(t, err)

	return tran
}

func TestListTransactionRun(t *testing.T) {
	transactionRepo, testsRepo := getRepos()

	t1 := createTransaction(t, transactionRepo, tests.Transaction{
		Name:        "first test",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	})

	t2 := createTransaction(t, transactionRepo, tests.Transaction{
		Name:        "second transaction",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	})

	run1, err := transactionRepo.CreateRun(context.TODO(), t1.NewRun())
	require.NoError(t, err)

	run2, err := transactionRepo.CreateRun(context.TODO(), t1.NewRun())
	require.NoError(t, err)

	_, err = transactionRepo.CreateRun(context.TODO(), t2.NewRun())
	require.NoError(t, err)

	runs, err := transactionRepo.GetTransactionsRuns(context.TODO(), t1.ID, 20, 0)
	require.NoError(t, err)

	assert.Len(t, runs, 2)
	assert.Equal(t, runs[0].ID, run2.ID)
	assert.Equal(t, runs[1].ID, run1.ID)
}

func TestBug(t *testing.T) {
	transactionRepo, testsRepo := getRepos()

	ctx := context.TODO()

	transaction := createTransaction(t, transactionRepo, tests.Transaction{
		Name:        "first test",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	})

	run1, err := transactionRepo.CreateRun(ctx, transaction.NewRun())
	require.NoError(t, err)

	run2, err := transactionRepo.CreateRun(ctx, transaction.NewRun())
	require.NoError(t, err)

	runs, err := transactionRepo.GetTransactionsRuns(ctx, transaction.ID, 20, 0)
	require.NoError(t, err)
	assert.Equal(t, runs[0].ID, run2.ID)
	assert.Equal(t, runs[1].ID, run1.ID)

	transaction.Name = "another thing"
	_, err = transactionRepo.Update(ctx, transaction)
	require.NoError(t, err)

	newTransaction, err := transactionRepo.GetAugmented(context.TODO(), transaction.ID)
	require.NoError(t, err)

	run3, err := transactionRepo.CreateRun(ctx, newTransaction.NewRun())
	require.NoError(t, err)

	runs, err = transactionRepo.GetTransactionsRuns(ctx, newTransaction.ID, 20, 0)
	require.NoError(t, err)

	assert.Len(t, runs, 3)
	assert.Equal(t, runs[0].ID, run3.ID)
	assert.Equal(t, runs[1].ID, run2.ID)
	assert.Equal(t, runs[2].ID, run1.ID)

}
