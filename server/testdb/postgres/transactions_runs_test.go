package postgres_test

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

	run := transaction.NewRun()
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

	run := transaction.NewRun()
	run, err = db.CreateTransactionRun(context.TODO(), run)
	require.NoError(t, err)

	run.State = model.TransactionRunStateExecuting
	err = db.UpdateTransactionRun(context.TODO(), run)
	require.NoError(t, err)

	updatedRun, err := db.GetTransactionRun(context.TODO(), transaction.ID, run.ID)
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

	run := transaction.NewRun()
	newRun, err := db.CreateTransactionRun(context.TODO(), run)
	require.NoError(t, err)

	err = db.DeleteTransactionRun(context.TODO(), newRun)
	require.NoError(t, err)

	_, err = db.GetTransactionRun(context.TODO(), transaction.ID, newRun.ID)
	require.ErrorContains(t, err, "record not found")
}

func TestListTransactionRun(t *testing.T) {
	db, clean := getDB()
	defer clean()

	transaction, err := db.CreateTransaction(context.TODO(), model.Transaction{
		Name:        "first test",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, db, "first step"),
			createTestWithName(t, db, "second step"),
		},
	})
	require.NoError(t, err)

	transaction2, err := db.CreateTransaction(context.TODO(), model.Transaction{
		Name:        "second transaction",
		Description: "description",
		Steps: []model.Test{
			createTestWithName(t, db, "first step"),
			createTestWithName(t, db, "second step"),
		},
	})
	require.NoError(t, err)

	run1, err := db.CreateTransactionRun(context.TODO(), transaction.NewRun())
	require.NoError(t, err)

	run2, err := db.CreateTransactionRun(context.TODO(), transaction.NewRun())
	require.NoError(t, err)

	_, err = db.CreateTransactionRun(context.TODO(), transaction2.NewRun())
	require.NoError(t, err)

	runs, err := db.GetTransactionsRuns(context.TODO(), transaction.ID, 20, 0)
	require.NoError(t, err)

	assert.Len(t, runs, 2)
	assert.Equal(t, runs[0].ID, run2.ID)
	assert.Equal(t, runs[1].ID, run1.ID)
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

	run1, err := db.CreateTransactionRun(ctx, transaction.NewRun())
	require.NoError(t, err)

	run2, err := db.CreateTransactionRun(ctx, transaction.NewRun())
	require.NoError(t, err)

	runs, err := db.GetTransactionsRuns(ctx, transaction.ID, 20, 0)
	require.NoError(t, err)
	assert.Equal(t, runs[0].ID, run2.ID)
	assert.Equal(t, runs[1].ID, run1.ID)

	transaction.Name = "another thing"
	newTransaction, err := db.UpdateTransaction(ctx, transaction)
	require.NoError(t, err)

	run3, err := db.CreateTransactionRun(ctx, newTransaction.NewRun())
	require.NoError(t, err)

	runs, err = db.GetTransactionsRuns(ctx, newTransaction.ID, 20, 0)
	require.NoError(t, err)

	assert.Len(t, runs, 3)
	assert.Equal(t, runs[0].ID, run3.ID)
	assert.Equal(t, runs[1].ID, run2.ID)
	assert.Equal(t, runs[2].ID, run1.ID)

}
