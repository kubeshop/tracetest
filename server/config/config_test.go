package config_test

import (
	"strings"
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
)

func TestFromFileSuccess(t *testing.T) {
	cases := []struct {
		name     string
		file     string
		expected config.Config
	}{
		{
			name: "Jaeger",
			file: "./testdata/jaeger.yaml",
			expected: config.Config{
				PoolingConfig: config.PoolingConfig{
					MaxWaitTimeForTrace: "1m",
					RetryDelay:          "3s",
				},
				PostgresConnString: "host=postgres user=postgres password=postgres port=5432 sslmode=disable",
				JaegerConnectionConfig: &configgrpc.GRPCClientSettings{
					Endpoint:   "jaeger-query:16685",
					TLSSetting: configtls.TLSClientSetting{Insecure: true},
				},
			},
		},
		{
			name: "Tempo",
			file: "./testdata/tempo.yaml",
			expected: config.Config{
				PoolingConfig: config.PoolingConfig{
					MaxWaitTimeForTrace: "1m",
					RetryDelay:          "1s",
				},
				PostgresConnString: "host=postgres user=postgres password=postgres port=5432 sslmode=disable",
				TempoConnectionConfig: &configgrpc.GRPCClientSettings{
					Endpoint:   "tempo:9095",
					TLSSetting: configtls.TLSClientSetting{Insecure: true},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cl := c
			t.Parallel()

			actual, err := config.FromFile(cl.file)

			require.NoError(t, err)
			assert.Equal(t, cl.expected, actual)

		})
	}
}

func TestFromFileError(t *testing.T) {
	cases := []struct {
		name     string
		file     string
		expected string
	}{
		{
			name:     "Jaeger",
			file:     "./testdata/notexists.yaml",
			expected: "read file: ",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cl := c
			t.Parallel()

			_, err := config.FromFile(cl.file)

			assert.True(t, strings.HasPrefix(err.Error(), cl.expected))

		})
	}
}
