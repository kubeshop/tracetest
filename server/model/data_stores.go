package model

import (
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

func (ds DataStore) IsZero() bool {
	return ds.Type == ""
}

func (ds DataStore) Slug() string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(ds.Name), " ", "-"))
}

func DataStoreFromConfig(dsc config.TracingBackendDataStoreConfig) DataStore {
	return DataStore{
		Type: openapi.SupportedDataStores(dsc.Type),
		Values: DataStoreValues{
			Jaeger:     &dsc.Jaeger,
			Tempo:      &dsc.Tempo,
			OpenSearch: &dsc.OpenSearch,
			SignalFx:   &dsc.SignalFX,
		},
	}
}
