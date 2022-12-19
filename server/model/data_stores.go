package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/openapi"
	"go.opentelemetry.io/collector/config/configgrpc"
)

type (
	DataStore struct {
		ID        string
		Name      string
		Type      openapi.SupportedDataStores
		IsDefault bool
		Values    DataStoreValues
		CreatedAt time.Time
	}

	DataStoreValues struct {
		Jaeger     *configgrpc.GRPCClientSettings
		Tempo      *configgrpc.GRPCClientSettings
		OpenSearch *config.OpensearchDataStoreConfig
		SignalFx   *config.SignalFXDataStoreConfig
	}
)

func (e DataStore) Slug() string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(e.Name), " ", "-"))
}

func (e DataStore) Config() (config.TracingBackendDataStoreConfig, error) {
	switch e.Type {
	case openapi.JAEGER:
		return config.TracingBackendDataStoreConfig{
			Type:   "jaeger",
			Jaeger: *e.Values.Jaeger,
		}, nil

	case openapi.OPEN_SEARCH:
		return config.TracingBackendDataStoreConfig{
			Type:       "opensearch",
			OpenSearch: *e.Values.OpenSearch,
		}, nil

	case openapi.OTLP:
		return config.TracingBackendDataStoreConfig{
			Type: "otlp",
		}, nil

	case openapi.SIGNAL_FX:
		return config.TracingBackendDataStoreConfig{
			Type:     "signalfx",
			SignalFX: *e.Values.SignalFx,
		}, nil

	case openapi.TEMPO:
		return config.TracingBackendDataStoreConfig{
			Type:  "tempo",
			Tempo: *e.Values.Tempo,
		}, nil

	default:
		return config.TracingBackendDataStoreConfig{}, fmt.Errorf("unsupported data store")
	}
}
