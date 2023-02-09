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

func TestBasicConfiguration(t *testing.T) {
	expectedConfig := config.Config{
		PoolingConfig: config.PoolingConfig{
			MaxWaitTimeForTrace: "1m",
			RetryDelay:          "3s",
		},
		PostgresConnString: "host=postgres user=postgres password=postgres port=5432 sslmode=disable",
		Server: config.ServerConfig{
			PathPrefix: "/tracetest",
			HttpPort:   9999,
		},
	}

	actual, err := config.FromFile("./testdata/basic_config.yaml")
	require.NoError(t, err)
	assert.Equal(t, expectedConfig, actual)
}

func TestFromFileError(t *testing.T) {
	_, err := config.FromFile("./testdata/notexists.yaml")
	assert.True(t, strings.HasPrefix(err.Error(), "read file: "))
}

func TestJaegerDataStore(t *testing.T) {
	expectedDataStore := config.TracingBackendDataStoreConfig{
		Type: "jaeger",
		Jaeger: configgrpc.GRPCClientSettings{
			Endpoint:   "jaeger-query:16685",
			TLSSetting: configtls.TLSClientSetting{Insecure: true},
		},
	}

	config, err := config.FromFile("./testdata/jaeger_datastore.yaml")
	require.NoError(t, err)

	selectedDataStore, err := config.DataStore()
	assert.NoError(t, err)
	assert.Equal(t, expectedDataStore, *selectedDataStore)
}

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

	config, err := config.FromFile("./testdata/tempo_datastore.yaml")
	require.NoError(t, err)

	selectedDataStore, err := config.DataStore()
	assert.NoError(t, err)
	assert.Equal(t, expectedDataStore, *selectedDataStore)
}

func TestInexistentDataStore(t *testing.T) {
	config, err := config.FromFile("./testdata/inexistent_datastore.yaml")
	assert.NoError(t, err)

	dataStore, err := config.DataStore()
	assert.Error(t, err)
	assert.Nil(t, dataStore)
}

func TestEmptyDataStore(t *testing.T) {
	config, err := config.FromFile("./testdata/empty_datastore.yaml")
	assert.NoError(t, err)

	dataStore, err := config.DataStore()

	assert.NoError(t, err)
	assert.Nil(t, dataStore)
}

func TestExporterConfig(t *testing.T) {
	expectedExporter := config.TelemetryExporterOption{
		ServiceName: "tracetest",
		Sampling:    100,
		Exporter: config.ExporterConfig{
			Type: "collector",
			CollectorConfiguration: config.OTELCollectorConfig{
				Endpoint: "collector:8888",
			},
		},
	}

	config, err := config.FromFile("./testdata/exporter_config.yaml")
	require.NoError(t, err)

	selectedExporter, err := config.Exporter()
	assert.NoError(t, err)
	assert.Equal(t, expectedExporter, *selectedExporter)

	selectedAppExporter, err := config.ApplicationExporter()
	assert.NoError(t, err)
	assert.Equal(t, *selectedExporter, *selectedAppExporter)
}

func TestInexistentExporter(t *testing.T) {
	config, err := config.FromFile("./testdata/inexistent_exporter.yaml")
	require.NoError(t, err)

	exporter, err := config.Exporter()
	assert.Error(t, err)
	assert.Nil(t, exporter)

	appExporter, err := config.ApplicationExporter()
	assert.Error(t, err)
	assert.Nil(t, appExporter)
}

func TestEmptyExporter(t *testing.T) {
	config, err := config.FromFile("./testdata/empty_exporter.yaml")
	require.NoError(t, err)

	exporter, err := config.Exporter()
	assert.NoError(t, err)
	assert.Nil(t, exporter)

	appExporter, err := config.ApplicationExporter()
	assert.NoError(t, err)
	assert.Nil(t, appExporter)
}
