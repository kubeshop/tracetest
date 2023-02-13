package config_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
)

func TestDataStoreExceptions(t *testing.T) {
	t.Run("Inexistent", func(t *testing.T) {
		t.Parallel()

		cfg := configFromFile(t, "./testdata/inexistent_datastore.yaml")

		dataStore, err := cfg.DataStore()
		assert.Error(t, err)
		assert.Nil(t, dataStore)
	})

	t.Run("Empty", func(t *testing.T) {
		t.Parallel()

		cfg := configFromFile(t, "./testdata/empty_datastore.yaml")

		dataStore, err := cfg.DataStore()
		assert.NoError(t, err)
		assert.Nil(t, dataStore)
	})
}

func TestDataStore(t *testing.T) {
	t.Run("JaegerGRPC", func(t *testing.T) {
		t.Parallel()

		expectedDataStore := &config.TracingBackendDataStoreConfig{
			Type: "jaeger",
			Jaeger: configgrpc.GRPCClientSettings{
				Endpoint:   "jaeger-query:16685",
				TLSSetting: configtls.TLSClientSetting{Insecure: true},
			},
		}

		cfg := configFromFile(t, "./testdata/jaeger_grpc.yaml")

		selectedDataStore, err := cfg.DataStore()
		assert.NoError(t, err)
		assert.Equal(t, expectedDataStore, selectedDataStore)
	})

	t.Run("TempoGRPC", func(t *testing.T) {
		t.Parallel()

		expectedDataStore := &config.TracingBackendDataStoreConfig{
			Type: "tempo",
			Tempo: config.BaseClientConfig{
				Grpc: configgrpc.GRPCClientSettings{
					Endpoint:   "tempo:9095",
					TLSSetting: configtls.TLSClientSetting{Insecure: true},
				},
			},
		}

		cfg := configFromFile(t, "./testdata/tempo_grpc.yaml")

		selectedDataStore, err := cfg.DataStore()
		assert.NoError(t, err)
		assert.Equal(t, expectedDataStore, selectedDataStore)
	})

	t.Run("TempoHTTP", func(t *testing.T) {
		t.Parallel()

		expectedDataStore := &config.TracingBackendDataStoreConfig{
			Type: "tempo",
			Tempo: config.BaseClientConfig{
				Http: config.HttpClientConfig{
					Url:        "tempo:80",
					TLSSetting: configtls.TLSClientSetting{Insecure: true},
				},
			},
		}

		cfg := configFromFile(t, "./testdata/tempo_http.yaml")

		selectedDataStore, err := cfg.DataStore()
		assert.NoError(t, err)
		assert.Equal(t, expectedDataStore, selectedDataStore)
	})

	t.Run("OpenSearch", func(t *testing.T) {
		t.Parallel()

		expectedDataStore := &config.TracingBackendDataStoreConfig{
			Type: "opensearch",
			OpenSearch: config.ElasticSearchDataStoreConfig{
				Addresses: []string{"http://opensearch:9200"},
				Index:     "traces",
			},
		}

		cfg := configFromFile(t, "./testdata/opensearch.yaml")

		selectedDataStore, err := cfg.DataStore()
		assert.NoError(t, err)
		assert.Equal(t, expectedDataStore, selectedDataStore)
	})

	t.Run("SignalFX", func(t *testing.T) {
		t.Parallel()

		expectedDataStore := &config.TracingBackendDataStoreConfig{
			Type: "signalfx",
			SignalFX: config.SignalFXDataStoreConfig{
				Token: "thetoken",
				Realm: "us1",
			},
		}

		cfg := configFromFile(t, "./testdata/signalfx.yaml")

		selectedDataStore, err := cfg.DataStore()
		assert.NoError(t, err)
		assert.Equal(t, expectedDataStore, selectedDataStore)
	})

	t.Run("ElasitcAPM", func(t *testing.T) {
		t.Parallel()

		expectedDataStore := &config.TracingBackendDataStoreConfig{
			Type: "elasticapm",
			ElasticApm: config.ElasticSearchDataStoreConfig{
				Addresses:          []string{"https://es01:9200"},
				Username:           "elastic",
				Password:           "changeme",
				Index:              "traces-apm-default",
				InsecureSkipVerify: true,
			},
		}

		cfg := configFromFile(t, "./testdata/elastic_apm.yaml")

		selectedDataStore, err := cfg.DataStore()
		assert.NoError(t, err)
		assert.Equal(t, expectedDataStore, selectedDataStore)
	})
}
