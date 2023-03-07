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
)

type mockPublisher struct {
	mock.Mock
}

func TestPublishing(t *testing.T) {
	updated := configresource.Config{
		AnalyticsEnabled: true,
	}

}

func TestIsAnalyticsEnabled(t *testing.T) {
	cleanEnv := func() {
		// make sure this env is empty
		os.Setenv("TRACETEST_DEV", "")
	}

	db := testmock.MustGetRawTestingDatabase()
	t.Run("DefaultValues", func(t *testing.T) {
		cleanEnv()
		repo := configresource.Repository(
			testmock.MustCreateRandomMigratedDatabase(db),
		)

		cfg := repo.Current(context.TODO())
		assert.True(t, cfg.IsAnalyticsEnabled())

	})

	t.Run("FromRepo", func(t *testing.T) {
		cleanEnv()
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
		Name:             "test",
		AnalyticsEnabled: true,
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
			}
		},
		SampleJSON: `{
			"type": "Config",
			"spec": {
				"id": "1",
				"name": "test",
				"analyticsEnabled": true
			}
		}`,

		SampleJSONUpdated: `{
			"type": "Config",
			"spec": {
				"id": "1",
				"name": "test updated",
				"analyticsEnabled": true
			}
		}`,
	})
}
