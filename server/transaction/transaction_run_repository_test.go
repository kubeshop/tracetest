package transaction_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/kubeshop/tracetest/server/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestWithName(t *testing.T, db test.Repository, name string) test.Test {
	t.Helper()
	test := test.Test{
		Name:        name,
		Description: "description",
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPRequest{
				URL: "http://localhost:3030/hello-instrumented",
			},
		},
	}

	updated, err := db.Create(context.TODO(), test)
	if err != nil {
		panic(err)
	}
	return updated
}

func getRepos() (*transaction.Repository, *transaction.RunRepository, test.Repository) {
	db := testmock.CreateMigratedDatabase()

	testRepo := test.NewRepository(db)
	testRunRepo := test.NewRunRepository(db)

	transactionRepo := transaction.NewRepository(db, testRepo)
	runRepo := transaction.NewRunRepository(db, testRunRepo)

	return transactionRepo, runRepo, testRepo
}

func getTransaction(t *testing.T, transactionRepo *transaction.Repository, testsRepo test.Repository) (transaction.Transaction, transactionFixture) {
	f := setupTransactionFixture(t, transactionRepo.DB())

	transaction := transaction.Transaction{
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
	transactionRepo, transactionRunRepo, testsRepo := getRepos()
	transactionObject, _ := getTransaction(t, transactionRepo, testsRepo)

	tr, err := transactionRunRepo.CreateRun(context.TODO(), transactionObject.NewRun())
	require.NoError(t, err)

	assert.Equal(t, tr.TransactionID, transactionObject.ID)
	assert.Equal(t, tr.State, transaction.TransactionRunStateCreated)
	assert.Len(t, tr.Steps, 0)
}

func TestUpdateTransactionRun(t *testing.T) {
	transactionRepo, transactionRunRepo, testsRepo := getRepos()
	transactionObject, fixture := getTransaction(t, transactionRepo, testsRepo)

	tr, err := transactionRunRepo.CreateRun(context.TODO(), transactionObject.NewRun())
	require.NoError(t, err)

	tr.State = transaction.TransactionRunStateExecuting
	tr.Steps = []test.Run{fixture.testRun}
	err = transactionRunRepo.UpdateRun(context.TODO(), tr)
	require.NoError(t, err)

	updatedRun, err := transactionRunRepo.GetTransactionRun(context.TODO(), transactionObject.ID, tr.ID)
	require.NoError(t, err)

	assert.Equal(t, tr.TransactionID, transactionObject.ID)
	assert.Equal(t, transaction.TransactionRunStateExecuting, updatedRun.State)
	assert.Len(t, tr.Steps, 1)
}

func TestDeleteTransactionRun(t *testing.T) {
	transactionRepo, transactionRunRepo, testsRepo := getRepos()
	transaction, _ := getTransaction(t, transactionRepo, testsRepo)

	tr, err := transactionRunRepo.CreateRun(context.TODO(), transaction.NewRun())
	require.NoError(t, err)

	err = transactionRunRepo.DeleteTransactionRun(context.TODO(), tr)
	require.NoError(t, err)

	_, err = transactionRunRepo.GetTransactionRun(context.TODO(), transaction.ID, tr.ID)
	require.ErrorContains(t, err, "no rows in result set")
}

func createTransaction(t *testing.T, repo *transaction.Repository, tran transaction.Transaction) transaction.Transaction {
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
	transactionRepo, transactionRunRepo, testsRepo := getRepos()

	t1 := createTransaction(t, transactionRepo, transaction.Transaction{
		Name:        "first test",
		Description: "description",
		Steps: []test.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	})

	t2 := createTransaction(t, transactionRepo, transaction.Transaction{
		Name:        "second transaction",
		Description: "description",
		Steps: []test.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	})

	run1, err := transactionRunRepo.CreateRun(context.TODO(), t1.NewRun())
	require.NoError(t, err)

	run2, err := transactionRunRepo.CreateRun(context.TODO(), t1.NewRun())
	require.NoError(t, err)

	_, err = transactionRunRepo.CreateRun(context.TODO(), t2.NewRun())
	require.NoError(t, err)

	runs, err := transactionRunRepo.GetTransactionsRuns(context.TODO(), t1.ID, 20, 0)
	require.NoError(t, err)

	assert.Len(t, runs, 2)
	assert.Equal(t, runs[0].ID, run2.ID)
	assert.Equal(t, runs[1].ID, run1.ID)
}

func TestBug(t *testing.T) {
	transactionRepo, transactionRunRepo, testsRepo := getRepos()

	ctx := context.TODO()

	transaction := createTransaction(t, transactionRepo, transaction.Transaction{
		Name:        "first test",
		Description: "description",
		Steps: []test.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	})

	run1, err := transactionRunRepo.CreateRun(ctx, transaction.NewRun())
	require.NoError(t, err)

	run2, err := transactionRunRepo.CreateRun(ctx, transaction.NewRun())
	require.NoError(t, err)

	runs, err := transactionRunRepo.GetTransactionsRuns(ctx, transaction.ID, 20, 0)
	require.NoError(t, err)
	assert.Equal(t, runs[0].ID, run2.ID)
	assert.Equal(t, runs[1].ID, run1.ID)

	transaction.Name = "another thing"
	_, err = transactionRepo.Update(ctx, transaction)
	require.NoError(t, err)

	newTransaction, err := transactionRepo.GetAugmented(context.TODO(), transaction.ID)
	require.NoError(t, err)

	run3, err := transactionRunRepo.CreateRun(ctx, newTransaction.NewRun())
	require.NoError(t, err)

	runs, err = transactionRunRepo.GetTransactionsRuns(ctx, newTransaction.ID, 20, 0)
	require.NoError(t, err)

	assert.Len(t, runs, 3)
	assert.Equal(t, runs[0].ID, run3.ID)
	assert.Equal(t, runs[1].ID, run2.ID)
	assert.Equal(t, runs[2].ID, run1.ID)

}
