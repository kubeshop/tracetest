package testdb_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEnvironment(t *testing.T) {
	db, clean := getDB()
	defer clean()

	environment := model.Environment{
		Name:        "first environment",
		Description: "description",
		Values:      []model.EnvironmentValue{{Key: "key1", Value: "value1"}, {Key: "key2", Value: "value2"}},
	}

	updated, err := db.CreateEnvironment(context.TODO(), environment)
	require.NoError(t, err)

	actual, err := db.GetEnvironment(context.TODO(), updated.ID)
	require.NoError(t, err)
	assert.Equal(t, environment.Name, actual.Name)
	assert.Equal(t, environment.Description, actual.Description)
	assert.Equal(t, environment.Values, actual.Values)
	assert.False(t, actual.CreatedAt.IsZero())
}

func TestDeleteEnvironment(t *testing.T) {
	db, clean := getDB()
	defer clean()

	environment := createEnvironment(t, db, "env1")

	err := db.DeleteEnvironment(context.TODO(), environment)
	require.NoError(t, err)

	actual, err := db.GetEnvironment(context.TODO(), environment.ID)
	assert.ErrorIs(t, err, testdb.ErrNotFound)
	assert.Empty(t, actual)
}

func TestUpdateEnvironment(t *testing.T) {
	db, clean := getDB()
	defer clean()

	environment := createEnvironment(t, db, "env1")
	environment.Name = "1 v2"
	environment.Description = "1 v2 description"

	_, err := db.UpdateEnvironment(context.TODO(), environment)
	require.NoError(t, err)

	latestTest, err := db.GetEnvironment(context.TODO(), environment.Slug())
	assert.NoError(t, err)
	assert.Equal(t, "1 v2", latestTest.Name)
	assert.Equal(t, "1 v2 description", latestTest.Description)
}

func TestGetEnvironments(t *testing.T) {
	db, clean := getDB()
	defer clean()

	createEnvironment(t, db, "env1")
	createEnvironment(t, db, "env2")
	createEnvironment(t, db, "env3")

	t.Run("Order", func(t *testing.T) {
		actual, err := db.GetEnvironments(context.TODO(), 20, 0, "", "", "")
		require.NoError(t, err)

		assert.Len(t, actual.Items, 3)
		assert.Equal(t, actual.TotalCount, 3)

		// test order
		assert.Equal(t, actual.TotalCount, 3)
		assert.Equal(t, "env3", actual.Items[0].Name)
		assert.Equal(t, "env2", actual.Items[1].Name)
		assert.Equal(t, "env1", actual.Items[2].Name)
	})

	t.Run("Pagination", func(t *testing.T) {
		actual, err := db.GetEnvironments(context.TODO(), 20, 10, "", "", "")
		require.NoError(t, err)

		assert.Equal(t, actual.TotalCount, 3)
		assert.Len(t, actual.Items, 0)
	})

	t.Run("SortByCreated", func(t *testing.T) {
		actual, err := db.GetEnvironments(context.TODO(), 20, 0, "", "created", "")
		require.NoError(t, err)

		// test order
		assert.Equal(t, "env3", actual.Items[0].Name)
		assert.Equal(t, "env2", actual.Items[1].Name)
		assert.Equal(t, "env1", actual.Items[2].Name)
	})

	t.Run("SortByNameAsc", func(t *testing.T) {
		actual, err := db.GetEnvironments(context.TODO(), 20, 0, "", "name", "asc")
		require.NoError(t, err)

		// test order
		assert.Equal(t, "env1", actual.Items[0].Name)
		assert.Equal(t, "env2", actual.Items[1].Name)
		assert.Equal(t, "env3", actual.Items[2].Name)
	})

	t.Run("SortByNameDesc", func(t *testing.T) {
		actual, err := db.GetEnvironments(context.TODO(), 20, 0, "", "name", "desc")
		require.NoError(t, err)

		// test order
		assert.Equal(t, "env3", actual.Items[0].Name)
		assert.Equal(t, "env2", actual.Items[1].Name)
		assert.Equal(t, "env1", actual.Items[2].Name)
	})

	t.Run("SearchByName", func(t *testing.T) {
		createEnvironment(t, db, "VerySpecificName")

		actual, err := db.GetEnvironments(context.TODO(), 10, 0, "specif", "", "")
		require.NoError(t, err)
		assert.Len(t, actual.Items, 1)
		assert.Equal(t, actual.TotalCount, 1)

		assert.Equal(t, "VerySpecificName", actual.Items[0].Name)
	})

	t.Run("SearchByDescription", func(t *testing.T) {
		_, _ = db.CreateEnvironment(context.TODO(), model.Environment{Description: "VeryUniqueText"})

		actual, err := db.GetEnvironments(context.TODO(), 10, 0, "niqueText", "", "")
		require.NoError(t, err)
		assert.Len(t, actual.Items, 1)
		assert.Equal(t, actual.TotalCount, 1)

		assert.Equal(t, "VeryUniqueText", actual.Items[0].Description)
	})
}
