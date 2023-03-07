package configresource_test

import (
	"context"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/config/configresource"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockPublisher struct {
	mock.Mock
}

func (m *mockPublisher) Publish(resourceID string, message any) {
	m.Called(resourceID, message)
}

func TestPublishing(t *testing.T) {
	restore := cleanEnv()
	defer restore()

	updated := configresource.Config{
		ID:               "123",
		AnalyticsEnabled: true,
	}

	publisher := new(mockPublisher)
	publisher.Test(t)

	publisher.On("Publish", configresource.ResourceID, updated)

	db := testmock.MustGetRawTestingDatabase()
	repo := configresource.Repository(
		testmock.MustCreateRandomMigratedDatabase(db),
		configresource.WithPublisher(publisher),
	)

	repo.Create(context.TODO(), configresource.Config{ID: "123"})

	_, err := repo.Update(context.TODO(), updated)
	require.NoError(t, err)

	publisher.AssertExpectations(t)

}

func TestIsAnalyticsEnabled(t *testing.T) {
	db := testmock.MustGetRawTestingDatabase()
	t.Run("DefaultValues", func(t *testing.T) {
		restore := cleanEnv()
		defer restore()

		repo := configresource.Repository(
			testmock.MustCreateRandomMigratedDatabase(db),
		)

		cfg := repo.Current(context.TODO())
		assert.True(t, cfg.IsAnalyticsEnabled())

	})

	t.Run("FromRepo", func(t *testing.T) {
		restore := cleanEnv()
		defer restore()
		repo := configresource.Repository(
			testmock.MustCreateRandomMigratedDatabase(db),
		)
		repo.Create(context.TODO(), configresource.Config{
			AnalyticsEnabled: false,
		})

		cfg := repo.Current(context.TODO())
		assert.False(t, cfg.IsAnalyticsEnabled())
	})

	t.Run("EnvOverride", func(t *testing.T) {
		repo := configresource.Repository(
			testmock.MustCreateRandomMigratedDatabase(db),
		)
		repo.Create(context.TODO(), configresource.Config{
			AnalyticsEnabled: true,
		})

		cfg := repo.Current(context.TODO())

		os.Setenv("TRACETEST_DEV", "yes")
		assert.False(t, cfg.IsAnalyticsEnabled())

	})
}

func TestConfigResource(t *testing.T) {

	db := testmock.MustGetRawTestingDatabase()
	sampleConfig := configresource.Config{
		ID:               "1",
		Name:             "test 1",
		AnalyticsEnabled: false,
	}
	secondSampleConfig := configresource.Config{
		ID:               "2",
		Name:             "test 2",
		AnalyticsEnabled: true,
	}
	thirdSampleConfig := configresource.Config{
		ID:               "3",
		Name:             "test 3",
		AnalyticsEnabled: false,
	}

	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceType: "Config",
		RegisterManagerFn: func(router *mux.Router) any {
			db := testmock.MustCreateRandomMigratedDatabase(db)
			configRepo := configresource.Repository(db)

			manager := resourcemanager.New[configresource.Config]("Config", configRepo, id.GenerateID)
			manager.RegisterRoutes(router)

			return configRepo
		},
		Prepare: func(t *testing.T, op rmtests.Operation, bridge any) {
			configRepo := bridge.(resourcemanager.ResourceHandler[configresource.Config])
			switch op {
			case rmtests.OperationGetSuccess,
				rmtests.OperationUpdateSuccess,
				rmtests.OperationDeleteSuccess,
				rmtests.OperationListSuccess:
				configRepo.Create(context.TODO(), sampleConfig)
			case rmtests.OperationListPaginatedSuccess:
				configRepo.Create(context.TODO(), sampleConfig)
				configRepo.Create(context.TODO(), secondSampleConfig)
				configRepo.Create(context.TODO(), thirdSampleConfig)
			}
		},
		SampleJSON: `{
			"type": "Config",
			"spec": {
				"id": "1",
				"name": "test 1",
				"analyticsEnabled": false
			}
		}`,
		SampleJSONUpdated: `{
			"type": "Config",
			"spec": {
				"id": "1",
				"name": "test updated",
				"analyticsEnabled": false
			}
		}`,
	})
}

func cleanEnv() func() {
	val := os.Getenv("TRACETEST_DEV")
	// make sure this env is empty
	os.Setenv("TRACETEST_DEV", "")
	return func() {
		os.Setenv("TRACETEST_DEV", val)
	}
}
