package testdb_test

import (
	"context"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/configgrpc"
)

func TestCreateDataStore(t *testing.T) {
	db, clean := getDB()
	defer clean()

	dataStore := model.DataStore{
		Name:      "datastore",
		Type:      "jaeger",
		IsDefault: true,
		Values: model.DataStoreValues{
			Jaeger:     &configgrpc.GRPCClientSettings{},
			Tempo:      &model.BaseClientConfig{},
			SignalFx:   &model.SignalFXDataStoreConfig{},
			OpenSearch: &model.ElasticSearchDataStoreConfig{},
			ElasticApm: &model.ElasticSearchDataStoreConfig{},
			AwsXRay:    &model.AWSXRayDataStoreConfig{},
		},
	}

	created, err := db.CreateDataStore(context.TODO(), dataStore)
	require.NoError(t, err)

	actual, err := db.GetDataStore(context.TODO(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, dataStore.Name, actual.Name)
	assert.Equal(t, dataStore.Type, actual.Type)
	assert.Equal(t, dataStore.IsDefault, actual.IsDefault)
	assert.Equal(t, dataStore.Values, actual.Values)
	assert.False(t, actual.CreatedAt.IsZero())
}

func TestCreateMultipleDataStores(t *testing.T) {
	db, clean := getDB()
	defer clean()

	createDataStore(t, db, "datastore1")
	createDataStore(t, db, "datastore2")
	createDataStore(t, db, "datastore3")

	actual, err := db.GetDataStores(context.TODO(), 20, 0, "", "", "")
	require.NoError(t, err)

	assert.Len(t, actual.Items, 3)
	assert.Equal(t, actual.TotalCount, 3)

	// test one default data store
	assert.Equal(t, actual.TotalCount, 3)
	assert.Equal(t, "datastore3", actual.Items[0].Name)
	assert.True(t, actual.Items[0].IsDefault)
	assert.Equal(t, "datastore2", actual.Items[1].Name)
	assert.False(t, actual.Items[1].IsDefault)
	assert.Equal(t, "datastore1", actual.Items[2].Name)
	assert.False(t, actual.Items[2].IsDefault)
}

func TestDeleteDataStore(t *testing.T) {
	db, clean := getDB()
	defer clean()

	dataStore := createDataStore(t, db, "datastore1")

	err := db.DeleteDataStore(context.TODO(), dataStore)
	require.NoError(t, err)

	actual, err := db.GetDataStore(context.TODO(), dataStore.ID)
	assert.ErrorIs(t, err, testdb.ErrNotFound)
	assert.Empty(t, actual)
}

func TestUpdateDataStore(t *testing.T) {
	db, clean := getDB()
	defer clean()

	dataStore := createDataStore(t, db, "datastore1")
	dataStore.Name = "1 v2"
	dataStore.Type = "openSearch"

	_, err := db.UpdateDataStore(context.TODO(), dataStore)
	require.NoError(t, err)

	latestDataStore, err := db.GetDataStore(context.TODO(), dataStore.ID)
	assert.NoError(t, err)
	assert.Equal(t, "1 v2", latestDataStore.Name)
	assert.Equal(t, "openSearch", string(latestDataStore.Type))
}

func TestGetDataStores(t *testing.T) {
	db, clean := getDB()
	defer clean()

	createDataStore(t, db, "datastore1")
	createDataStore(t, db, "datastore2")
	createDataStore(t, db, "datastore3")

	t.Run("Order", func(t *testing.T) {
		actual, err := db.GetDataStores(context.TODO(), 20, 0, "", "", "")
		require.NoError(t, err)

		assert.Len(t, actual.Items, 3)
		assert.Equal(t, actual.TotalCount, 3)

		// test order
		assert.Equal(t, actual.TotalCount, 3)
		assert.Equal(t, "datastore3", actual.Items[0].Name)
		assert.Equal(t, "datastore2", actual.Items[1].Name)
		assert.Equal(t, "datastore1", actual.Items[2].Name)
	})

	t.Run("Pagination", func(t *testing.T) {
		actual, err := db.GetDataStores(context.TODO(), 20, 10, "", "", "")
		require.NoError(t, err)

		assert.Equal(t, actual.TotalCount, 3)
		assert.Len(t, actual.Items, 0)
	})

	t.Run("SortByCreated", func(t *testing.T) {
		actual, err := db.GetDataStores(context.TODO(), 20, 0, "", "created", "")
		require.NoError(t, err)

		// test order
		assert.Equal(t, "datastore3", actual.Items[0].Name)
		assert.Equal(t, "datastore2", actual.Items[1].Name)
		assert.Equal(t, "datastore1", actual.Items[2].Name)
	})

	t.Run("SortByNameAsc", func(t *testing.T) {
		actual, err := db.GetDataStores(context.TODO(), 20, 0, "", "name", "asc")
		require.NoError(t, err)

		// test order
		assert.Equal(t, "datastore1", actual.Items[0].Name)
		assert.Equal(t, "datastore2", actual.Items[1].Name)
		assert.Equal(t, "datastore3", actual.Items[2].Name)
	})

	t.Run("SortByNameDesc", func(t *testing.T) {
		actual, err := db.GetDataStores(context.TODO(), 20, 0, "", "name", "desc")
		require.NoError(t, err)

		// test order
		assert.Equal(t, "datastore3", actual.Items[0].Name)
		assert.Equal(t, "datastore2", actual.Items[1].Name)
		assert.Equal(t, "datastore1", actual.Items[2].Name)
	})

	t.Run("SearchByName", func(t *testing.T) {
		createDataStore(t, db, "VerySpecificName")

		actual, err := db.GetDataStores(context.TODO(), 10, 0, "specif", "", "")
		require.NoError(t, err)
		assert.Len(t, actual.Items, 1)
		assert.Equal(t, actual.TotalCount, 1)

		assert.Equal(t, "VerySpecificName", actual.Items[0].Name)
	})

}

func TestDataStoreProvisioner(t *testing.T) {

	db := testmock.MustGetRawTestingDatabase()

	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: testdb.DataStoreResourceName,
		ResourceTypePlural:   testdb.DataStoreResourceName,
		RegisterManagerFn: func(router *mux.Router) resourcemanager.Manager {
			db := testmock.MustCreateRandomMigratedDatabase(db)
			dsRepo, err := testdb.Postgres(testdb.WithDB(db))
			require.NoError(t, err)

			manager := resourcemanager.New[testdb.DataStoreResource](
				testdb.DataStoreResourceName,
				testdb.DataStoreResourceNamePlural,
				testdb.NewDataStoreResourceProvisioner(dsRepo),
				// this resource exists only for provisiooning at the moment``
				resourcemanager.WithOperations(resourcemanager.OperationNoop),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		SampleJSON: `{
			"type": "DataStore",
			"spec": {
				"id": "signalfx",
				"name": "SignalFX",
				"type": "signalfx",
				"signalfx": {
					"token": "thetoken",
    			"realm": "us1"
				}
			}
		}`,
	})
}
