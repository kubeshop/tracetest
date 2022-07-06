package testdb_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTest(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := model.Test{
		Name:        "first test",
		Description: "description",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				URL: "http://localhost:3030/hello-instrumented",
			},
		},
	}

	updated, err := db.CreateTest(context.TODO(), test)
	require.NoError(t, err)

	actual, err := db.GetLatestTestVersion(context.TODO(), updated.ID)
	require.NoError(t, err)
	assert.Equal(t, test.Name, actual.Name)
	assert.Equal(t, test.Description, actual.Description)
	assert.Equal(t, test.ServiceUnderTest, actual.ServiceUnderTest)
	assert.Equal(t, test.ReferenceRun, actual.ReferenceRun)
	assert.Equal(t, test.Definition, actual.Definition)
	assert.False(t, actual.CreatedAt.IsZero())
}

func TestUpdateTest(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := createTest(t, db)
	test.Name = "updated test"

	err := db.UpdateTestVersion(context.TODO(), test)
	require.NoError(t, err)

	actual, err := db.GetLatestTestVersion(context.TODO(), test.ID)
	require.NoError(t, err)

	assert.False(t, actual.CreatedAt.IsZero())
	actual.CreatedAt = test.CreatedAt // hack to ignore created at in equal comparation
	assert.Equal(t, test, actual)
}

func TestDeleteTest(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := createTest(t, db)

	err := db.DeleteTest(context.TODO(), test)
	require.NoError(t, err)

	actual, err := db.GetLatestTestVersion(context.TODO(), test.ID)
	assert.ErrorIs(t, err, testdb.ErrNotFound)
	assert.Empty(t, actual)

}

func TestGetLatestTestVersion(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := createTestWithName(t, db, "1")
	test.Name = "1 v2"
	test.Version = 2

	_, err := db.UpdateTest(context.TODO(), test)
	require.NoError(t, err)

	latestTest, err := db.GetLatestTestVersion(context.TODO(), test.ID)
	assert.NoError(t, err)
	assert.Equal(t, "1 v2", latestTest.Name)
	assert.Equal(t, 2, latestTest.Version)
}

func TestGetTests(t *testing.T) {
	db, clean := getDB()
	defer clean()

	createTestWithName(t, db, "1")
	createTestWithName(t, db, "2")
	createTestWithName(t, db, "3")

	actual, err := db.GetTests(context.TODO(), 20, 0)
	require.NoError(t, err)
	assert.Len(t, actual, 3)

	// test order
	assert.Equal(t, "3", actual[0].Name)
	assert.Equal(t, "2", actual[1].Name)
	assert.Equal(t, "1", actual[2].Name)

	actual, err = db.GetTests(context.TODO(), 20, 10)
	require.NoError(t, err)
	assert.Len(t, actual, 0)
}

func TestGetTestsWithMultipleVersions(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test1 := createTestWithName(t, db, "1")
	test1.Name = "1 v2"

	_, err := db.UpdateTest(context.TODO(), test1)
	require.NoError(t, err)

	test2 := createTestWithName(t, db, "2")
	test2.Name = "2 v2"

	_, err = db.UpdateTest(context.TODO(), test2)
	require.NoError(t, err)

	tests, err := db.GetTests(context.TODO(), 20, 0)
	assert.NoError(t, err)
	assert.Len(t, tests, 2)

	for _, test := range tests {
		assert.Equal(t, 2, test.Version)
	}
}
