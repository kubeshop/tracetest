package config_test

import (
	"os"
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func configWithFlagsE(t *testing.T, inputFlags []string, opts ...config.Option) (*config.Config, error) {
	flags := pflag.NewFlagSet("fake", pflag.ExitOnError)
	config.SetupFlags(flags)

	err := flags.Parse(inputFlags)
	require.NoError(t, err)

	return config.New(append(opts, config.WithFlagSet(flags))...)
}

func configWithFlags(t *testing.T, inputFlags []string, opts ...config.Option) *config.Config {
	cfg, err := configWithFlagsE(t, inputFlags, opts...)
	require.NoError(t, err)

	return cfg
}

func configFromFile(t *testing.T, path string, opts ...config.Option) *config.Config {
	return configWithFlags(t, []string{"--config", path}, opts...)
}

func configWithEnv(t *testing.T, env map[string]string) *config.Config {
	for k, v := range env {
		os.Setenv(k, v)
	}

	cfg, err := config.New()
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

		cfg, err := config.New(config.WithFlagSet(flags))
		assert.Nil(t, cfg)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("ConfigFileOK", func(t *testing.T) {
		t.Parallel()

		cfg := configFromFile(t, "./testdata/basic.yaml")

		assert.Equal(t, "host=postgres user=postgres password=postgres port=5432 dbname=tracetest sslmode=disable", cfg.PostgresConnString())

		assert.Equal(t, "/tracetest", cfg.ServerPathPrefix())
		assert.Equal(t, 9999, cfg.ServerPort())
	})

	t.Run("ConfigFileDefault", func(t *testing.T) {
		t.Parallel()

		// These tests are not great, since they rely on file writing and removing.
		// Becase of this, the tests cannot be run in parallel because they might interfer with each other.

		t.Run("OK", func(t *testing.T) {
			// copy an example config file to the default location
			err := copyFile("./testdata/basic.yaml", "./tracetest.yaml")
			defer os.Remove("./tracetest.yaml")

			require.NoError(t, err)

			cfg, err := config.New()
			require.NoError(t, err)

			// this one assertion is enough to guarantee we're not using the defaults
			assert.Equal(t, 9999, cfg.ServerPort())
		})

		t.Run("MustHaveExtension", func(t *testing.T) {
			// copy an example config file to the default location
			err := copyFile("./testdata/basic.yaml", "./tracetest")
			defer os.Remove("./tracetest")

			require.NoError(t, err)

			cfg, err := config.New()
			require.NoError(t, err)

			// the config file would change this value to 9999, but we want to make sure
			// the file is NOT being read, since it doesn't have the .yaml extension
			assert.Equal(t, 11633, cfg.ServerPort())
		})
	})

}

func copyFile(in, out string) error {
	b, err := os.ReadFile(in)
	if err != nil {
		return err
	}

	return os.WriteFile(out, b, 0644)
}
