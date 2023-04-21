package mappings

import (
	"fmt"
	"time"

	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/tracedb/datastoreresource"
)

func (m *OpenAPI) ConnectionTestResult(in model.ConnectionResult) openapi.ConnectionResult {
	result := openapi.ConnectionResult{}

	if in.PortCheck.IsSet() {
		result.PortCheck = m.ConnectionTestStep(in.PortCheck)
	}

	if in.Connectivity.IsSet() {
		result.Connectivity = m.ConnectionTestStep(in.Connectivity)
	}

	if in.Authentication.IsSet() {
		result.Authentication = m.ConnectionTestStep(in.Authentication)
	}

	if in.FetchTraces.IsSet() {
		result.FetchTraces = m.ConnectionTestStep(in.FetchTraces)
	}

	return result
}

func (m *OpenAPI) ConnectionTestStep(in model.ConnectionTestStep) openapi.ConnectionTestStep {
	errMessage := ""
	if in.Error != nil {
		errMessage = in.Error.Error()
	}

	return openapi.ConnectionTestStep{
		Passed:  in.Status != model.StatusFailed,
		Message: in.Message,
		Status:  string(in.Status),
		Error:   errMessage,
	}
}

var dataStoreTypesMapping = map[datastoreresource.DataStoreType]openapi.SupportedDataStores{
	datastoreresource.DataStoreTypeJaeger:     openapi.JAEGER,
	datastoreresource.DataStoreTypeTempo:      openapi.TEMPO,
	datastoreresource.DataStoreTypeOpenSearch: openapi.OPEN_SEARCH,
	datastoreresource.DataStoreTypeSignalFX:   openapi.SIGNAL_FX,
	datastoreresource.DataStoreTypeOTLP:       openapi.OTLP,
	datastoreresource.DataStoreTypeNewRelic:   openapi.NEW_RELIC,
	datastoreresource.DataStoreTypeLighStep:   openapi.LIGHTSTEP,
	datastoreresource.DataStoreTypeElasticAPM: openapi.ELASTIC_APM,
	datastoreresource.DataStoreTypeDataDog:    openapi.DATADOG,
	datastoreresource.DataStoreTypeAwsXRay:    openapi.AWSXRAY,
}

func (m OpenAPI) DataStoreType(in datastoreresource.DataStoreType) openapi.SupportedDataStores {
	dsd, exists := dataStoreTypesMapping[in]
	if !exists {
		// this should only happen during development,
		// so it's more an alert for devs than actual error handling
		panic(fmt.Errorf("trying to convert an undefined datastoreresource.DataStoreType '%s'", in))
	}
	return dsd
}

func (m Model) DataStoreType(in openapi.SupportedDataStores) datastoreresource.DataStoreType {
	for k, v := range dataStoreTypesMapping {
		if v == in {
			return k
		}
	}

	// this should only happen during development,
	// so it's more an alert for devs than actual error handling
	panic(fmt.Errorf("trying to convert an undefined datastoreresource.DataStoreType '%s'", in))

}

func (m Model) DataStore(in openapi.DataStore) datastoreresource.DataStore {
	dataStore := datastoreresource.DataStore{
		ID:        "current",
		Name:      in.Name,
		Type:      m.DataStoreType(in.Type),
		Default:   in.Default,
		CreatedAt: in.CreatedAt.Format(time.RFC3339Nano),
	}

	deepcopy.DeepCopy(in.Jaeger, &dataStore.Values.Jaeger)
	deepcopy.DeepCopy(in.Jaeger.Tls, &dataStore.Values.Jaeger.TLS)

	deepcopy.DeepCopy(in.Tempo, &dataStore.Values.Tempo)
	deepcopy.DeepCopy(in.Tempo.Grpc.Tls, &dataStore.Values.Tempo.Grpc.TLS)
	deepcopy.DeepCopy(in.Tempo.Http.Tls, &dataStore.Values.Tempo.Http.TLS)

	deepcopy.DeepCopy(in.OpenSearch, &dataStore.Values.OpenSearch)
	deepcopy.DeepCopy(in.ElasticApm, &dataStore.Values.ElasticApm)
	deepcopy.DeepCopy(in.SignalFx, &dataStore.Values.SignalFx)
	deepcopy.DeepCopy(in.Awsxray, &dataStore.Values.AwsXRay)

	return dataStore
}
