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

func configWithFlagsE(t *testing.T, inputFlags []string) (*config.Config, error) {
	flags := pflag.NewFlagSet("fake", pflag.ExitOnError)
	config.SetupFlags(flags)

	err := flags.Parse(inputFlags)
	require.NoError(t, err)

	return config.New(flags)
}

func configWithFlags(t *testing.T, inputFlags []string) *config.Config {
	cfg, err := configWithFlagsE(t, inputFlags)
	require.NoError(t, err)

	return cfg
}

func configFromFile(t *testing.T, path string) *config.Config {
	return configWithFlags(t, []string{"--config", path})
}

func configWithEnv(t *testing.T, env map[string]string) *config.Config {
	for k, v := range env {
		os.Setenv(k, v)
	}

	cfg, err := config.New(nil)
	require.NoError(t, err)

	return cfg
}

func TestFlags(t *testing.T) {

	t.Run("ConfigFileOverrideNotExists", func(t *testing.T) {
		t.Parallel()

		flags := pflag.NewFlagSet("fake", pflag.ExitOnError)
		config.SetupFlags(flags)

		err := flags.Parse([]string{"--config", "notexists.yaml"})
		require.NoError(t, err)

		cfg, err := config.New(flags)
		assert.Nil(t, cfg)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("ConfigFileOK", func(t *testing.T) {
		t.Parallel()

		cfg := configFromFile(t, "./testdata/basic.yaml")

		assert.Equal(t, "host=postgres user=postgres password=postgres port=5432 dbname=tracetest sslmode=disable", cfg.PostgresConnString())

		assert.Equal(t, "/tracetest", cfg.ServerPathPrefix())
		assert.Equal(t, 9999, cfg.ServerPort())

		assert.Equal(t, 1*time.Minute, cfg.PoolingMaxWaitTimeForTraceDuration())
		assert.Equal(t, 3*time.Second, cfg.PoolingRetryDelay())
	})

}
