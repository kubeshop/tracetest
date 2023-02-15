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
		Tempo:      openapi.BaseClient{},
		OpenSearch: openapi.ElasticSearch{},
		ElasticApm: openapi.ElasticSearch{},
		SignalFx:   openapi.SignalFx{},
		Awsxray:    openapi.AwsXRay{},
		CreatedAt:  in.CreatedAt,
	}

	deepcopy.DeepCopy(in.Values.Jaeger, &dataStore.Jaeger)
	deepcopy.DeepCopy(in.Values.Jaeger.TLSSetting, &dataStore.Jaeger.Tls)
	deepcopy.DeepCopy(in.Values.Jaeger.TLSSetting.TLSSetting, &dataStore.Jaeger.Tls.Settings)

	deepcopy.DeepCopy(in.Values.Tempo, &dataStore.Tempo)
	deepcopy.DeepCopy(in.Values.Tempo.Grpc.TLSSetting, &dataStore.Tempo.Grpc.Tls)
	deepcopy.DeepCopy(in.Values.Tempo.Grpc.TLSSetting.TLSSetting, &dataStore.Tempo.Grpc.Tls.Settings)
	deepcopy.DeepCopy(in.Values.Tempo.Http.TLSSetting, &dataStore.Tempo.Http.Tls)
	deepcopy.DeepCopy(in.Values.Tempo.Http.TLSSetting, &dataStore.Tempo.Http.Tls.Settings)

	deepcopy.DeepCopy(in.Values.OpenSearch, &dataStore.OpenSearch)
	deepcopy.DeepCopy(in.Values.ElasticApm, &dataStore.ElasticApm)
	deepcopy.DeepCopy(in.Values.SignalFx, &dataStore.SignalFx)
	deepcopy.DeepCopy(in.Values.AwsXRay, &dataStore.Awsxray)

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
	deepcopy.DeepCopy(in.Tempo.Grpc.Tls, &dataStore.Values.Tempo.Grpc.TLSSetting)
	deepcopy.DeepCopy(in.Tempo.Grpc.Tls.Settings, &dataStore.Values.Tempo.Grpc.TLSSetting.TLSSetting)
	deepcopy.DeepCopy(in.Tempo.Http.Tls, &dataStore.Values.Tempo.Http.TLSSetting)
	deepcopy.DeepCopy(in.Tempo.Http.Tls.Settings, &dataStore.Values.Tempo.Grpc.TLSSetting.TLSSetting)

	deepcopy.DeepCopy(in.OpenSearch, &dataStore.Values.OpenSearch)
	deepcopy.DeepCopy(in.ElasticApm, &dataStore.Values.ElasticApm)
	deepcopy.DeepCopy(in.SignalFx, &dataStore.Values.SignalFx)
	deepcopy.DeepCopy(in.Awsxray, &dataStore.Values.AwsXRay)

	return dataStore
}
