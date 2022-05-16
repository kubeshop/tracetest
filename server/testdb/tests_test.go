package testdb_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/model"
	"github.com/kubeshop/tracetest/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTest(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := model.Test{
		Name:        "first test",
		Description: "description",
		ServiceUnderTest: model.ServiceUnderTest{
			Request: model.HTTPRequest{
				URL: "http://localhost:3030/hello-instrumented",
			},
		},
	}

	updated, err := db.CreateTest(context.TODO(), test)
	require.NoError(t, err)

	actual, err := db.GetTest(context.TODO(), updated.ID)
	require.NoError(t, err)
	assert.Equal(t, test, actual)
}

func TestUpdateTest(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := createTest(db)
	test.Name = "updated test"

	err := db.UpdateTest(context.TODO(), test)
	require.NoError(t, err)

	actual, err := db.GetTest(context.TODO(), test.ID)
	require.NoError(t, err)
	assert.Equal(t, test, actual)
}

func TestDeleteTest(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := createTest(db)

	err := db.DeleteTest(context.TODO(), test)
	require.NoError(t, err)

	actual, err := db.GetTest(context.TODO(), test.ID)
	assert.ErrorIs(t, err, testdb.ErrNotFound)
	assert.Empty(t, actual)

}

func TestGetTests(t *testing.T) {
	db, clean := getDB()
	defer clean()

	createTestWithName(db, "1")
	createTestWithName(db, "2")

	actual, err := db.GetTests(context.TODO(), 20, 0)
	require.NoError(t, err)
	assert.Len(t, actual, 2)

	actual, err = db.GetTests(context.TODO(), 20, 10)
	require.NoError(t, err)
	assert.Len(t, actual, 0)
}
