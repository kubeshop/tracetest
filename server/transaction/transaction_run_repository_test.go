<<<<<<<< HEAD:server/transaction/transaction_run_repository_test.go
package transaction_test
========
package transactions_test
>>>>>>>> 7fb86839 (fix: move transactions to it's own module (#2664)):server/transactions/transactions_runs_test.go

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/testmock"
<<<<<<<< HEAD:server/transaction/transaction_run_repository_test.go
	"github.com/kubeshop/tracetest/server/transaction"
========
	"github.com/kubeshop/tracetest/server/transactions"
>>>>>>>> 7fb86839 (fix: move transactions to it's own module (#2664)):server/transactions/transactions_runs_test.go
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

<<<<<<<< HEAD:server/transaction/transaction_run_repository_test.go
func getRepos() (*transaction.Repository, *transaction.RunRepository, model.Repository) {
========
func getRepos() (*transactions.TransactionsRepository, model.Repository) {
>>>>>>>> 7fb86839 (fix: move transactions to it's own module (#2664)):server/transactions/transactions_runs_test.go
	db := testmock.GetRawTestingDatabase()

	testsRepo, err := testdb.Postgres(testdb.WithDB(db))
	if err != nil {
		panic(err)
	}

<<<<<<<< HEAD:server/transaction/transaction_run_repository_test.go
	transactionRepo := transaction.NewRepository(db, testsRepo)
========
	transactionRepo := transactions.NewTransactionsRepository(db, testsRepo)
>>>>>>>> 7fb86839 (fix: move transactions to it's own module (#2664)):server/transactions/transactions_runs_test.go

	runRepo := transaction.NewRunRepository(db, testsRepo)

	return transactionRepo, runRepo, testsRepo
}

<<<<<<<< HEAD:server/transaction/transaction_run_repository_test.go
func getTransaction(t *testing.T, transactionRepo *transaction.Repository, testsRepo model.TestRepository) (transaction.Transaction, transactionFixture) {
	f := setupTransactionFixture(t, transactionRepo.DB())

	transaction := transaction.Transaction{
========
func getTransaction(t *testing.T, transactionRepo *transactions.TransactionsRepository, testsRepo model.TestRepository) (transactions.Transaction, transactionFixture) {
	f := setupTransactionFixture(t, transactionRepo.DB())

	transaction := transactions.Transaction{
>>>>>>>> 7fb86839 (fix: move transactions to it's own module (#2664)):server/transactions/transactions_runs_test.go
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

<<<<<<<< HEAD:server/transaction/transaction_run_repository_test.go
	assert.Equal(t, tr.TransactionID, transactionObject.ID)
	assert.Equal(t, tr.State, transaction.TransactionRunStateCreated)
========
	assert.Equal(t, tr.TransactionID, transaction.ID)
	assert.Equal(t, tr.State, transactions.TransactionRunStateCreated)
>>>>>>>> 7fb86839 (fix: move transactions to it's own module (#2664)):server/transactions/transactions_runs_test.go
	assert.Len(t, tr.Steps, 0)
}

func TestUpdateTransactionRun(t *testing.T) {
	transactionRepo, transactionRunRepo, testsRepo := getRepos()
	transactionObject, fixture := getTransaction(t, transactionRepo, testsRepo)

	tr, err := transactionRunRepo.CreateRun(context.TODO(), transactionObject.NewRun())
	require.NoError(t, err)

<<<<<<<< HEAD:server/transaction/transaction_run_repository_test.go
	tr.State = transaction.TransactionRunStateExecuting
========
	tr.State = transactions.TransactionRunStateExecuting
>>>>>>>> 7fb86839 (fix: move transactions to it's own module (#2664)):server/transactions/transactions_runs_test.go
	tr.Steps = []model.Run{fixture.testRun}
	err = transactionRunRepo.UpdateRun(context.TODO(), tr)
	require.NoError(t, err)

	updatedRun, err := transactionRunRepo.GetTransactionRun(context.TODO(), transactionObject.ID, tr.ID)
	require.NoError(t, err)

<<<<<<<< HEAD:server/transaction/transaction_run_repository_test.go
	assert.Equal(t, tr.TransactionID, transactionObject.ID)
	assert.Equal(t, transaction.TransactionRunStateExecuting, updatedRun.State)
========
	assert.Equal(t, tr.TransactionID, transaction.ID)
	assert.Equal(t, transactions.TransactionRunStateExecuting, updatedRun.State)
>>>>>>>> 7fb86839 (fix: move transactions to it's own module (#2664)):server/transactions/transactions_runs_test.go
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

<<<<<<<< HEAD:server/transaction/transaction_run_repository_test.go
func createTransaction(t *testing.T, repo *transaction.Repository, tran transaction.Transaction) transaction.Transaction {
========
func createTransaction(t *testing.T, repo *transactions.TransactionsRepository, tran transactions.Transaction) transactions.Transaction {
>>>>>>>> 7fb86839 (fix: move transactions to it's own module (#2664)):server/transactions/transactions_runs_test.go
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

<<<<<<<< HEAD:server/transaction/transaction_run_repository_test.go
	t1 := createTransaction(t, transactionRepo, transaction.Transaction{
========
	t1 := createTransaction(t, transactionRepo, transactions.Transaction{
>>>>>>>> 7fb86839 (fix: move transactions to it's own module (#2664)):server/transactions/transactions_runs_test.go
		Name:        "first test",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	})

<<<<<<<< HEAD:server/transaction/transaction_run_repository_test.go
	t2 := createTransaction(t, transactionRepo, transaction.Transaction{
========
	t2 := createTransaction(t, transactionRepo, transactions.Transaction{
>>>>>>>> 7fb86839 (fix: move transactions to it's own module (#2664)):server/transactions/transactions_runs_test.go
		Name:        "second transaction",
		Description: "description",
		Steps: []model.Test{
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

<<<<<<<< HEAD:server/transaction/transaction_run_repository_test.go
	transaction := createTransaction(t, transactionRepo, transaction.Transaction{
========
	transaction := createTransaction(t, transactionRepo, transactions.Transaction{
>>>>>>>> 7fb86839 (fix: move transactions to it's own module (#2664)):server/transactions/transactions_runs_test.go
		Name:        "first test",
		Description: "description",
		Steps: []model.Test{
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
