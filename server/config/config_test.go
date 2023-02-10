package config_test

import (
	"os"
	"testing"

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

}

// func TestBasicConfiguration(t *testing.T) {
// 	actual, err := config.FromFile("./testdata/basic_config.yaml")
// 	require.NoError(t, err)

// 	assert.Equal(t, "host=postgres user=postgres password=postgres port=5432 sslmode=disable", actual.PostgresConnString())

// 	assert.Equal(t, "/tracetest", actual.ServerPathPrefix())
// 	assert.Equal(t, 9999, actual.ServerPort())

// 	expectedPoolingConfig := config.PoolingConfig{
// 		MaxWaitTimeForTrace: "1m",
// 		RetryDelay:          "3s",
// 	}
// 	assert.Equal(t, expectedPoolingConfig, actual.PoolingConfig())
// }

func TestTempoDataStore(t *testing.T) {
	expectedDataStore := config.TracingBackendDataStoreConfig{
		Type: "tempo",
		Tempo: config.BaseClientConfig{
			Grpc: configgrpc.GRPCClientSettings{
				Endpoint:   "tempo:9095",
				TLSSetting: configtls.TLSClientSetting{Insecure: true},
			},
		},
	}

// func TestJaegerDataStore(t *testing.T) {
// 	expectedDataStore := config.TracingBackendDataStoreConfig{
// 		Type: "jaeger",
// 		Jaeger: configgrpc.GRPCClientSettings{
// 			Endpoint:   "jaeger-query:16685",
// 			TLSSetting: configtls.TLSClientSetting{Insecure: true},
// 		},
// 	}

// 	config, err := config.FromFile("./testdata/jaeger_datastore.yaml")
// 	require.NoError(t, err)

// 	selectedDataStore, err := config.DataStore()
// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedDataStore, *selectedDataStore)
// }

// func TestTempoDataStore(t *testing.T) {
// 	expectedDataStore := config.TracingBackendDataStoreConfig{
// 		Type: "tempo",
// 		Tempo: configgrpc.GRPCClientSettings{
// 			Endpoint:   "tempo:9095",
// 			TLSSetting: configtls.TLSClientSetting{Insecure: true},
// 		},
// 	}

// 	config, err := config.FromFile("./testdata/tempo_datastore.yaml")
// 	require.NoError(t, err)

// 	selectedDataStore, err := config.DataStore()
// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedDataStore, *selectedDataStore)
// }

// func TestInexistentDataStore(t *testing.T) {
// 	config, err := config.FromFile("./testdata/inexistent_datastore.yaml")
// 	assert.NoError(t, err)

// 	dataStore, err := config.DataStore()
// 	assert.Error(t, err)
// 	assert.Nil(t, dataStore)
// }

// func TestEmptyDataStore(t *testing.T) {
// 	config, err := config.FromFile("./testdata/empty_datastore.yaml")
// 	assert.NoError(t, err)

// 	dataStore, err := config.DataStore()

// 	assert.NoError(t, err)
// 	assert.Nil(t, dataStore)
// }

// func TestExporterConfig(t *testing.T) {
// 	expectedExporter := config.TelemetryExporterOption{
// 		ServiceName: "tracetest",
// 		Sampling:    100,
// 		Exporter: config.ExporterConfig{
// 			Type: "collector",
// 			CollectorConfiguration: config.OTELCollectorConfig{
// 				Endpoint: "collector:8888",
// 			},
// 		},
// 	}

// 	config, err := config.FromFile("./testdata/exporter_config.yaml")
// 	require.NoError(t, err)

// 	selectedExporter, err := config.Exporter()
// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedExporter, *selectedExporter)

// 	selectedAppExporter, err := config.ApplicationExporter()
// 	assert.NoError(t, err)
// 	assert.Equal(t, *selectedExporter, *selectedAppExporter)
// }

// func TestInexistentExporter(t *testing.T) {
// 	config, err := config.FromFile("./testdata/inexistent_exporter.yaml")
// 	require.NoError(t, err)

// 	exporter, err := config.Exporter()
// 	assert.Error(t, err)
// 	assert.Nil(t, exporter)

// 	appExporter, err := config.ApplicationExporter()
// 	assert.Error(t, err)
// 	assert.Nil(t, appExporter)
// }

// func TestEmptyExporter(t *testing.T) {
// 	config, err := config.FromFile("./testdata/empty_exporter.yaml")
// 	require.NoError(t, err)

// 	exporter, err := config.Exporter()
// 	assert.NoError(t, err)
// 	assert.Nil(t, exporter)

// 	appExporter, err := config.ApplicationExporter()
// 	assert.NoError(t, err)
// 	assert.Nil(t, appExporter)
// }
