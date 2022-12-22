package mappings

import (
	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
)

func (m OpenAPI) DataStore(in model.DataStore) openapi.DataStore {
	dataStore := openapi.DataStore{
		Id:         in.ID,
		Name:       in.Name,
		Type:       in.Type,
		IsDefault:  in.IsDefault,
		Jaeger:     openapi.GrpcClientSettings{},
		Tempo:      openapi.GrpcClientSettings{},
		OpenSearch: openapi.OpenSearch{},
		SignalFx:   openapi.SignalFx{},
		CreatedAt:  in.CreatedAt,
	}

	deepcopy.DeepCopy(in.Values, &dataStore)

	return dataStore
}

func (m OpenAPI) DataStores(in []model.DataStore) []openapi.DataStore {
	dataStores := make([]openapi.DataStore, len(in))
	for i, t := range in {
		dataStores[i] = m.DataStore(t)
	}

	return dataStores
}

func (m Model) DataStore(in openapi.DataStore) model.DataStore {
	dataStore := model.DataStore{
		ID:        in.Id,
		Name:      in.Name,
		Type:      in.Type,
		IsDefault: in.IsDefault,
		CreatedAt: in.CreatedAt,
	}

	deepcopy.DeepCopy(in.Jaeger, &dataStore.Values.Jaeger)
	deepcopy.DeepCopy(in.Jaeger.Tls, &dataStore.Values.Jaeger.TLSSetting)
	deepcopy.DeepCopy(in.Jaeger.Tls.Settings, &dataStore.Values.Jaeger.TLSSetting.TLSSetting)
	deepcopy.DeepCopy(in.Tempo, &dataStore.Values.Tempo)
	deepcopy.DeepCopy(in.Tempo.Tls, &dataStore.Values.Tempo.TLSSetting)
	deepcopy.DeepCopy(in.Tempo.Tls.Settings, &dataStore.Values.Tempo.TLSSetting.TLSSetting)
	deepcopy.DeepCopy(in.OpenSearch, &dataStore.Values.OpenSearch)
	deepcopy.DeepCopy(in.SignalFx, &dataStore.Values.SignalFx)

	return dataStore
}
