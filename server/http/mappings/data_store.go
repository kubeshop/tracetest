package mappings

import (
	"fmt"

	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
)

var dataStoreTypesMapping = map[model.DataStoreType]openapi.SupportedDataStores{
	model.DataStoreTypeJaeger:     openapi.JAEGER,
	model.DataStoreTypeTempo:      openapi.TEMPO,
	model.DataStoreTypeOpenSearch: openapi.OPEN_SEARCH,
	model.DataStoreTypeSignalFX:   openapi.SIGNAL_FX,
	model.DataStoreTypeOTLP:       openapi.OTLP,
	model.DataStoreTypeNewRelic:   openapi.NEW_RELIC,
	model.DataStoreTypeLighStep:   openapi.LIGHTSTEP,
	model.DataStoreTypeElasticAPM: openapi.ELASTIC_APM,
	model.DataStoreTypeDataDog:    openapi.DATADOG,
	model.DataStoreTypeAwsXRay:    openapi.AWSXRAY,
}

func (m OpenAPI) DataStoreType(in model.DataStoreType) openapi.SupportedDataStores {
	dsd, exists := dataStoreTypesMapping[in]
	if !exists {
		// this should only happen during development,
		// so it's more an alert for devs than actual error handling
		panic(fmt.Errorf("trying to convert an undefined model.DataStoreType '%s'", in))
	}
	return dsd
}

func (m OpenAPI) DataStore(in model.DataStore) openapi.DataStore {
	dataStore := openapi.DataStore{
		Id:         in.ID,
		Name:       in.Name,
		Type:       m.DataStoreType(in.Type),
		IsDefault:  in.IsDefault,
		Jaeger:     openapi.GrpcClientSettings{},
		Tempo:      openapi.BaseClient{},
		OpenSearch: openapi.ElasticSearch{},
		ElasticApm: openapi.ElasticSearch{},
		SignalFx:   openapi.SignalFx{},
		Awsxray:    openapi.AwsXRay{},
		CreatedAt:  in.CreatedAt,
	}

	if in.Values.Jaeger != nil {
		deepcopy.DeepCopy(in.Values.Jaeger, &dataStore.Jaeger)
		deepcopy.DeepCopy(in.Values.Jaeger.TLSSetting, &dataStore.Jaeger.Tls)
		deepcopy.DeepCopy(in.Values.Jaeger.TLSSetting.TLSSetting, &dataStore.Jaeger.Tls.Settings)
	}

	if in.Values.Tempo != nil {
		deepcopy.DeepCopy(in.Values.Tempo, &dataStore.Tempo)
		deepcopy.DeepCopy(in.Values.Tempo.Grpc.TLSSetting, &dataStore.Tempo.Grpc.Tls)
		deepcopy.DeepCopy(in.Values.Tempo.Grpc.TLSSetting.TLSSetting, &dataStore.Tempo.Grpc.Tls.Settings)
		deepcopy.DeepCopy(in.Values.Tempo.Http.TLSSetting, &dataStore.Tempo.Http.Tls)
		deepcopy.DeepCopy(in.Values.Tempo.Http.TLSSetting, &dataStore.Tempo.Http.Tls.Settings)
	}

	if in.Values.OpenSearch != nil {
		deepcopy.DeepCopy(in.Values.OpenSearch, &dataStore.OpenSearch)
	}
	if in.Values.ElasticApm != nil {
		deepcopy.DeepCopy(in.Values.ElasticApm, &dataStore.ElasticApm)
	}
	if in.Values.SignalFx != nil {
		deepcopy.DeepCopy(in.Values.SignalFx, &dataStore.SignalFx)
	}

	if in.Values.AwsXRay != nil {
		deepcopy.DeepCopy(in.Values.AwsXRay, &dataStore.Awsxray)
	}

	return dataStore
}

func (m OpenAPI) DataStores(in []model.DataStore) []openapi.DataStore {
	dataStores := make([]openapi.DataStore, len(in))
	for i, t := range in {
		dataStores[i] = m.DataStore(t)
	}

	return dataStores
}

func (m Model) DataStoreType(in openapi.SupportedDataStores) model.DataStoreType {
	for k, v := range dataStoreTypesMapping {
		if v == in {
			return k
		}
	}

	// this should only happen during development,
	// so it's more an alert for devs than actual error handling
	panic(fmt.Errorf("trying to convert an undefined model.DataStoreType '%s'", in))

}

func (m Model) DataStore(in openapi.DataStore) model.DataStore {
	dataStore := model.DataStore{
		ID:        in.Id,
		Name:      in.Name,
		Type:      m.DataStoreType(in.Type),
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
