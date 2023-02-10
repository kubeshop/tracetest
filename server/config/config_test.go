package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlags(t *testing.T) {

	t.Run("config", func(t *testing.T) {
		t.Parallel()

		flags := pflag.NewFlagSet("fake", pflag.ExitOnError)
		config.SetupFlags(flags)

		err := flags.Parse([]string{"--config", "notexists.yaml"})
		require.NoError(t, err)

		cfg, err := config.New(flags)
		assert.Nil(t, cfg)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	configFromFile := func(t *testing.T, path string) *config.Config {
		flags := pflag.NewFlagSet("fake", pflag.ExitOnError)
		config.SetupFlags(flags)

		err := flags.Parse([]string{"--config", path})
		require.NoError(t, err)

		cfg, err := config.New(flags)
		require.NoError(t, err)

		return cfg
	}

	t.Run("BasicConfig", func(t *testing.T) {
		t.Parallel()

		cfg := configFromFile(t, "./testdata/basic_config.yaml")

		assert.Equal(t, "host=postgres user=postgres password=postgres port=5432 dbname=tracetest sslmode=disable", cfg.PostgresConnString())

		assert.Equal(t, "/tracetest", cfg.ServerPathPrefix())
		assert.Equal(t, 9999, cfg.ServerPort())

		assert.Equal(t, 1*time.Minute, cfg.PoolingMaxWaitTimeForTraceDuration())
		assert.Equal(t, 3*time.Second, cfg.PoolingRetryDelay())
	})

}
