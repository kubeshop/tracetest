package mappings

import (
	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/openapi"
)

func (m OpenAPI) DataStore(name string, in config.TracingBackendDataStoreConfig) openapi.DataStore {
	config := openapi.DataStore{
		Type:       openapi.SupportedDataStores(in.Type),
		Name:       name,
		Jaeger:     openapi.GrpcClientSettings{},
		Tempo:      openapi.GrpcClientSettings{},
		OpenSearch: openapi.OpenSearch{},
		SignalFx:   openapi.SignalFx{},
	}

	deepcopy.DeepCopy(in.Jaeger, &config.Jaeger)
	deepcopy.DeepCopy(in.Tempo, &config.Tempo)
	deepcopy.DeepCopy(in.OpenSearch, &config.OpenSearch)
	deepcopy.DeepCopy(in.SignalFX, &config.SignalFx)

	return config
}

func (m Model) DataStore(in openapi.DataStore) config.TracingBackendDataStoreConfig {
	config := config.TracingBackendDataStoreConfig{
		Type: string(in.Type),
	}

	deepcopy.DeepCopy(in.Jaeger, &config.Jaeger)
	deepcopy.DeepCopy(in.Jaeger.Tls, &config.Jaeger.TLSSetting)
	deepcopy.DeepCopy(in.Tempo, &config.Tempo)
	deepcopy.DeepCopy(in.Tempo.Tls, &config.Tempo.TLSSetting)
	deepcopy.DeepCopy(in.OpenSearch, &config.OpenSearch)
	deepcopy.DeepCopy(in.SignalFx, &config.SignalFX)

	return config
}
