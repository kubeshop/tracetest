package testdb_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTransactionRun(t *testing.T) {
	db, clean := getDB()
	defer clean()

	transaction := model.Transaction{
		Name:        "first test",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, db, "first step"),
			createTestWithName(t, db, "second step"),
		},
	}

	transaction, err := db.CreateTransaction(context.TODO(), transaction)
	require.NoError(t, err)

	run := model.NewTransactionRun(transaction)
	newRun, err := db.CreateTransactionRun(context.TODO(), run)
	require.NoError(t, err)

	assert.Equal(t, newRun.TransactionID, transaction.ID)
	assert.Equal(t, newRun.State, model.TransactionRunStateCreated)
}

func TestUpdateTransactionRun(t *testing.T) {
	db, clean := getDB()
	defer clean()

	transaction := model.Transaction{
		Name:        "first test",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, db, "first step"),
			createTestWithName(t, db, "second step"),
		},
	}

	transaction, err := db.CreateTransaction(context.TODO(), transaction)
	require.NoError(t, err)

	run := model.NewTransactionRun(transaction)
	run, err = db.CreateTransactionRun(context.TODO(), run)
	require.NoError(t, err)

	run.State = model.TransactionRunStateExecuting
	err = db.UpdateTransactionRun(context.TODO(), run)
	require.NoError(t, err)

	updatedRun, err := db.GetTransactionRun(context.TODO(), transaction.ID.String(), run.ID)
	require.NoError(t, err)

	assert.Equal(t, run.TransactionID, transaction.ID)
	assert.Equal(t, model.TransactionRunStateExecuting, updatedRun.State)
}

func TestDeleteTransactionRun(t *testing.T) {
	db, clean := getDB()
	defer clean()

	transaction := model.Transaction{
		Name:        "first test",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, db, "first step"),
			createTestWithName(t, db, "second step"),
		},
	}

	transaction, err := db.CreateTransaction(context.TODO(), transaction)
	require.NoError(t, err)

	run := model.NewTransactionRun(transaction)
	newRun, err := db.CreateTransactionRun(context.TODO(), run)
	require.NoError(t, err)

	err = db.DeleteTransactionRun(context.TODO(), newRun)
	require.NoError(t, err)

	_, err = db.GetTransactionRun(context.TODO(), transaction.ID.String(), newRun.ID)
	require.ErrorContains(t, err, "record not found")
}

func TestListTransactionRun(t *testing.T) {
	db, clean := getDB()
	defer clean()

	transaction := model.Transaction{
		Name:        "first test",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, db, "first step"),
			createTestWithName(t, db, "second step"),
		},
	}

	transaction2 := model.Transaction{
		Name:        "second transaction",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, db, "first step"),
			createTestWithName(t, db, "second step"),
		},
	}

	transaction, err := db.CreateTransaction(context.TODO(), transaction)
	require.NoError(t, err)

	transaction2, err = db.CreateTransaction(context.TODO(), transaction2)
	require.NoError(t, err)

	run1 := model.NewTransactionRun(transaction)
	newRun1, err := db.CreateTransactionRun(context.TODO(), run1)
	require.NoError(t, err)

	run2 := model.NewTransactionRun(transaction)
	newRun2, err := db.CreateTransactionRun(context.TODO(), run2)
	require.NoError(t, err)

	run3 := model.NewTransactionRun(transaction2)
	newRun3, err := db.CreateTransactionRun(context.TODO(), run3)
	require.NoError(t, err)

	runs, err := db.GetTransactionsRuns(context.TODO(), transaction.ID.String(), 20, 0)
	require.NoError(t, err)

	assert.Len(t, runs, 2)
	assert.Contains(t, runs, newRun1)
	assert.Contains(t, runs, newRun2)
	assert.NotContains(t, runs, newRun3)
}

func TestBug(t *testing.T) {
	db, clean := getDB()
	defer clean()

	ctx := context.TODO()

	transaction := model.Transaction{
		Name:        "first test",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, db, "first step"),
			createTestWithName(t, db, "second step"),
		},
	}

	transaction, err := db.CreateTransaction(ctx, transaction)
	require.NoError(t, err)

	run1 := model.NewTransactionRun(transaction)
	run1, err = db.CreateTransactionRun(ctx, run1)
	require.NoError(t, err)

	run2 := model.NewTransactionRun(transaction)
	run2, err = db.CreateTransactionRun(ctx, run2)
	require.NoError(t, err)

	runs, err := db.GetTransactionsRuns(ctx, transaction.ID.String(), 20, 0)
	require.NoError(t, err)
	assert.Contains(t, runs, run1)
	assert.Contains(t, runs, run2)

	transaction.Name = "another thing"
	newTransaction, err := db.UpdateTransaction(ctx, transaction)
	require.NoError(t, err)

	run3 := model.NewTransactionRun(newTransaction)
	run3, err = db.CreateTransactionRun(ctx, run3)
	require.NoError(t, err)

	runs, err = db.GetTransactionsRuns(ctx, newTransaction.ID.String(), 20, 0)
	require.NoError(t, err)
	assert.Len(t, runs, 3)
	assert.Contains(t, runs, run1)
	assert.Contains(t, runs, run2)
	assert.Contains(t, runs, run3)

}
