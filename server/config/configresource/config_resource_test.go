package configresource_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/config/configresource"
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
		ID:               "current",
		Name:             "Config",
		AnalyticsEnabled: true,
	}

	publisher := new(mockPublisher)
	publisher.Test(t)

	publisher.On("Publish", configresource.ResourceID, updated)

	repo := configresource.NewRepository(
		testmock.CreateMigratedDatabase(),
		configresource.WithPublisher(publisher),
	)

	_, err := repo.Update(context.TODO(), updated)
	require.NoError(t, err)

	publisher.AssertExpectations(t)

}

func TestIsAnalyticsEnabled(t *testing.T) {
	t.Run("DefaultValues", func(t *testing.T) {
		restore := cleanEnv()
		defer restore()

		repo := configresource.NewRepository(
			testmock.CreateMigratedDatabase(),
		)

		cfg := repo.Current(context.TODO())
		assert.True(t, cfg.IsAnalyticsEnabled())

	})

	t.Run("FromRepo", func(t *testing.T) {
		restore := cleanEnv()
		defer restore()
		repo := configresource.NewRepository(
			testmock.CreateMigratedDatabase(),
		)
		repo.Update(context.TODO(), configresource.Config{
			AnalyticsEnabled: false,
		})

		cfg := repo.Current(context.TODO())
		assert.False(t, cfg.IsAnalyticsEnabled())
	})

	t.Run("EnvOverride", func(t *testing.T) {
		restore := cleanEnv()
		defer restore()
		repo := configresource.NewRepository(
			testmock.CreateMigratedDatabase(),
		)
		repo.Update(context.TODO(), configresource.Config{
			AnalyticsEnabled: true,
		})

		cfg := repo.Current(context.TODO())

		os.Setenv("TRACETEST_DEV", "yes")
		assert.False(t, cfg.IsAnalyticsEnabled())

	})
}

func TestConfigResource(t *testing.T) {
	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: configresource.ResourceName,
		ResourceTypePlural:   configresource.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			configRepo := configresource.NewRepository(db)

			manager := resourcemanager.New[configresource.Config](
				configresource.ResourceName,
				configresource.ResourceNamePlural,
				configRepo,
				resourcemanager.WithOperations(configresource.Operations...),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		SampleJSON: `{
			"type": "Config",
			"spec": {
				"id": "current",
				"name": "Config",
				"analyticsEnabled": true
			}
		}`,
		SampleJSONUpdated: `{
			"type": "Config",
			"spec": {
				"id": "current",
				"name": "Config",
				"analyticsEnabled": false
			}
		}`,
	},
		rmtests.ExcludeOperations(
			rmtests.OperationGetNotFound,
			rmtests.OperationUpdateNotFound,
		),
	)
}

func cleanEnv() func() {
	val := os.Getenv("TRACETEST_DEV")
	// make sure this env is empty
	os.Setenv("TRACETEST_DEV", "")
	return func() {
		os.Setenv("TRACETEST_DEV", val)
	}
}
