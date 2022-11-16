package testdb_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {
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

	updated, err := db.CreateTransaction(context.TODO(), transaction)
	require.NoError(t, err)

	actual, err := db.GetLatestTransactionVersion(context.TODO(), updated.ID)
	require.NoError(t, err)
	assert.Equal(t, transaction.Name, actual.Name)
	assert.Equal(t, transaction.Description, actual.Description)
	assert.False(t, actual.CreatedAt.IsZero())

	require.Equal(t, len(transaction.Steps), len(actual.Steps))
	for i, step := range transaction.Steps {
		actualStep := actual.Steps[i]
		assert.Equal(t, step.ID, actualStep.ID)
	}
}

func TestUpdateTransactionStepsOrder(t *testing.T) {
	db, clean := getDB()
	defer clean()

	transaction := createTransaction(t, db)

	transaction.Steps = []model.Test{
		createTestWithName(t, db, "first step"),
		createTestWithName(t, db, "second step"),
	}

	newTransaction, err := db.UpdateTransaction(context.TODO(), transaction)
	require.NoError(t, err)

	for i, step := range transaction.Steps {
		newStep := newTransaction.Steps[i]
		assert.Equal(t, step.ID, newStep.ID)
	}

	transaction.Steps = []model.Test{
		transaction.Steps[0],
		createTestWithName(t, db, "new second step"),
		transaction.Steps[1],
	}

	newTransaction, err = db.UpdateTransaction(context.TODO(), transaction)
	require.NoError(t, err)

	for i, step := range transaction.Steps {
		newStep := newTransaction.Steps[i]
		assert.Equal(t, step.ID, newStep.ID)
	}

	// old version must exist
	oldTransaction, err := db.GetTransactionVersion(context.TODO(), newTransaction.ID, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, oldTransaction.Version)
	assert.Len(t, oldTransaction.Steps, 2)
}

func TestDeleteTransaction(t *testing.T) {
	db, clean := getDB()
	defer clean()

	transaction := createTransaction(t, db)

	err := db.DeleteTransaction(context.TODO(), transaction)
	require.NoError(t, err)

	actual, err := db.GetLatestTransactionVersion(context.TODO(), transaction.ID)
	assert.ErrorIs(t, err, testdb.ErrNotFound)
	assert.Empty(t, actual)

}

func TestGetLatestTransactionVersion(t *testing.T) {
	db, clean := getDB()
	defer clean()

	transaction := createTransaction(t, db)
	transaction.Name = "1 v2"
	transaction.Version = 2

	_, err := db.UpdateTransaction(context.TODO(), transaction)
	require.NoError(t, err)

	latestTransaction, err := db.GetLatestTransactionVersion(context.TODO(), transaction.ID)
	assert.NoError(t, err)
	assert.Equal(t, "1 v2", latestTransaction.Name)
	assert.Equal(t, 2, latestTransaction.Version)
}

func TestGetTransactionVersion(t *testing.T) {
	db, clean := getDB()
	defer clean()

	transaction := createTransaction(t, db)
	transaction.Name = "1 v2"

	_, err := db.UpdateTransaction(context.TODO(), transaction)
	require.NoError(t, err)

	transaction.Name = "1 v3"

	_, err = db.UpdateTransaction(context.TODO(), transaction)
	require.NoError(t, err)

	latestTransaction, err := db.GetTransactionVersion(context.TODO(), transaction.ID, 2)
	assert.NoError(t, err)
	assert.Equal(t, "1 v2", latestTransaction.Name)
	assert.Equal(t, 2, latestTransaction.Version)
}

func TestGetTransactions(t *testing.T) {
	db, clean := getDB()
	defer clean()

	createTransactionWithName(t, db, "one")
	createTransactionWithName(t, db, "two")
	createTransactionWithName(t, db, "three")

	t.Run("Order", func(t *testing.T) {
		actual, err := db.GetTransactions(context.TODO(), 20, 0, "", "", "")
		require.NoError(t, err)

		assert.Len(t, actual.Items, 3)
		assert.Equal(t, 3, actual.TotalCount)

		// test order
		assert.Equal(t, actual.TotalCount, 3)
		assert.Equal(t, "three", actual.Items[0].Name)
		assert.Equal(t, "two", actual.Items[1].Name)
		assert.Equal(t, "one", actual.Items[2].Name)
	})

	t.Run("Pagination", func(t *testing.T) {
		actual, err := db.GetTransactions(context.TODO(), 20, 10, "", "", "")
		require.NoError(t, err)

		assert.Equal(t, actual.TotalCount, 3)
		assert.Len(t, actual.Items, 0)
	})

	t.Run("SortByCreated", func(t *testing.T) {
		actual, err := db.GetTransactions(context.TODO(), 20, 0, "", "created", "")
		require.NoError(t, err)

		// test order
		assert.Equal(t, "three", actual.Items[0].Name)
		assert.Equal(t, "two", actual.Items[1].Name)
		assert.Equal(t, "one", actual.Items[2].Name)
	})

	t.Run("SortByNameAsc", func(t *testing.T) {
		actual, err := db.GetTransactions(context.TODO(), 20, 0, "", "name", "asc")
		require.NoError(t, err)

		// test order
		assert.Equal(t, "one", actual.Items[0].Name)
		assert.Equal(t, "three", actual.Items[1].Name)
		assert.Equal(t, "two", actual.Items[2].Name)
	})

	t.Run("SortByNameDesc", func(t *testing.T) {
		actual, err := db.GetTransactions(context.TODO(), 20, 0, "", "name", "desc")
		require.NoError(t, err)

		// test order
		assert.Equal(t, "two", actual.Items[0].Name)
		assert.Equal(t, "three", actual.Items[1].Name)
		assert.Equal(t, "one", actual.Items[2].Name)
	})

	t.Run("SearchByName", func(t *testing.T) {
		_, _ = db.CreateTransaction(context.TODO(), model.Transaction{Name: "VerySpecificName"})
		actual, err := db.GetTransactions(context.TODO(), 10, 0, "specif", "", "")
		require.NoError(t, err)
		assert.Len(t, actual.Items, 1)
		assert.Equal(t, actual.TotalCount, 1)

		assert.Equal(t, "VerySpecificName", actual.Items[0].Name)
	})

	t.Run("SearchByDescription", func(t *testing.T) {
		_, _ = db.CreateTransaction(context.TODO(), model.Transaction{Description: "VeryUniqueText"})

		actual, err := db.GetTransactions(context.TODO(), 10, 0, "nique", "", "")
		require.NoError(t, err)
		assert.Len(t, actual.Items, 1)
		assert.Equal(t, actual.TotalCount, 1)

		assert.Equal(t, "VeryUniqueText", actual.Items[0].Description)
	})
}

func TestGetTransactionsWithMultipleVersions(t *testing.T) {
	db, clean := getDB()
	defer clean()

	transaction1 := createTransactionWithName(t, db, "1")
	transaction1.Name = "1 v2"

	_, err := db.UpdateTransaction(context.TODO(), transaction1)
	require.NoError(t, err)

	transaction2 := createTransactionWithName(t, db, "2")
	transaction2.Name = "2 v2"

	_, err = db.UpdateTransaction(context.TODO(), transaction2)
	require.NoError(t, err)

	tests, err := db.GetTransactions(context.TODO(), 20, 0, "", "", "")
	assert.NoError(t, err)
	assert.Len(t, tests.Items, 2)
	assert.Equal(t, 2, tests.TotalCount)

	for _, test := range tests.Items {
		assert.Equal(t, 2, test.Version)
	}
}

func TestGetTransactionsStepsWithMultipleVersions(t *testing.T) {
	db, clean := getDB()
	defer clean()

	ctx := context.TODO()

	transaction := createTransaction(t, db)

	test1 := createTestWithName(t, db, "first step")
	test2 := createTestWithName(t, db, "second step")

	transaction.Steps = []model.Test{
		test1,
		test2,
	}

	_, err := db.UpdateTransaction(context.TODO(), transaction)
	require.NoError(t, err)

	updatedTransaction, err := db.GetLatestTransactionVersion(ctx, transaction.ID)
	require.NoError(t, err)

	assert.Len(t, updatedTransaction.Steps, 2)

	test1.Name = "new test name"
	updatedTest, err := db.UpdateTest(ctx, test1)
	require.NoError(t, err)

	updatedTransaction, err = db.GetLatestTransactionVersion(ctx, transaction.ID)
	require.NoError(t, err)

	assert.Len(t, updatedTransaction.Steps, 2)

	assert.Equal(t, updatedTransaction.Steps[0].ID, updatedTest.ID)
	assert.Equal(t, updatedTransaction.Steps[0].Version, updatedTest.Version)

	assert.Equal(t, updatedTransaction.Steps[1].ID, test2.ID)
	assert.Equal(t, updatedTransaction.Steps[1].Version, test2.Version)
}
