package config_test

import (
	"context"
	"os"
	"testing"

	"github.com/kubeshop/tracetest/server/config"
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

	updated := config.Config{
		ID:               "current",
		Name:             "Config",
		AnalyticsEnabled: true,
	}

	publisher := new(mockPublisher)
	publisher.Test(t)

	publisher.On("Publish", config.ResourceID, updated)

	repo := config.NewRepository(
		testmock.CreateMigratedDatabase(),
		config.WithPublisher(publisher),
	)

	_, err := repo.Update(context.TODO(), updated)
	require.NoError(t, err)

	publisher.AssertExpectations(t)

}

func TestIsAnalyticsEnabled(t *testing.T) {
	t.Run("DefaultValues", func(t *testing.T) {
		restore := cleanEnv()
		defer restore()

		repo := config.NewRepository(
			testmock.CreateMigratedDatabase(),
		)

		cfg := repo.Current(context.TODO())
		assert.True(t, cfg.IsAnalyticsEnabled())

	})

	t.Run("FromRepo", func(t *testing.T) {
		restore := cleanEnv()
		defer restore()
		repo := config.NewRepository(
			testmock.CreateMigratedDatabase(),
		)
		repo.Update(context.TODO(), config.Config{
			AnalyticsEnabled: false,
		})

		cfg := repo.Current(context.TODO())
		assert.False(t, cfg.IsAnalyticsEnabled())
	})

	t.Run("EnvOverride", func(t *testing.T) {
		restore := cleanEnv()
		defer restore()
		repo := config.NewRepository(
			testmock.CreateMigratedDatabase(),
		)
		repo.Update(context.TODO(), config.Config{
			AnalyticsEnabled: true,
		})

		cfg := repo.Current(context.TODO())

		os.Setenv("TRACETEST_DEV", "yes")
		assert.False(t, cfg.IsAnalyticsEnabled())

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
