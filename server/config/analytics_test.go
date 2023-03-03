package config_test

import (
	"context"
	"os"
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/config/configresource"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/stretchr/testify/assert"
)

func TestAnalyticsEnabledConfig(t *testing.T) {
	cleanEnv := func() {
		// make sure this env is empty
		os.Setenv("TRACETEST_DEV", "")
	}

	db := testmock.MustGetRawTestingDatabase()

	setup := func(existingConfigs ...configresource.Config) config.Option {
		repo := configresource.Repository(
			testmock.MustCreateRandomMigratedDatabase(db),
		)
		for _, ec := range existingConfigs {
			_, err := repo.Create(context.TODO(), ec)
			if err != nil {
				panic(err)
			}
		}

		return config.WithConfigResource(repo)
	}
	t.Run("DefaultValues", func(t *testing.T) {
		cleanEnv()
		opt := setup()
		cfg, _ := config.New(opt)

		assert.True(t, cfg.AnalyticsEnabled())
	})

	t.Run("FromRepo", func(t *testing.T) {
		cleanEnv()
		opt := setup(configresource.Config{
			AnalyticsEnabled: false,
		})

		cfg, _ := config.New(opt)

		assert.False(t, cfg.AnalyticsEnabled())
	})

	t.Run("EnvOverride", func(t *testing.T) {
		opt := setup(configresource.Config{
			AnalyticsEnabled: true,
		})
		cfg, _ := config.New(opt)

		os.Setenv("TRACETEST_DEV", "yes")
		assert.False(t, cfg.AnalyticsEnabled())
	})
}
