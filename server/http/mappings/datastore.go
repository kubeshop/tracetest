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
	datastoreresource.DataStoreTypeOpenSearch: openapi.OPENSEARCH,
	datastoreresource.DataStoreTypeSignalFX:   openapi.SIGNALFX,
	datastoreresource.DataStoreTypeOTLP:       openapi.OTLP,
	datastoreresource.DataStoreTypeNewRelic:   openapi.NEWRELIC,
	datastoreresource.DataStoreTypeLighStep:   openapi.LIGHTSTEP,
	datastoreresource.DataStoreTypeElasticAPM: openapi.ELASTICAPM,
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
		Values:    datastoreresource.DataStoreValues{},
	}

	// Jaeger
	if dataStore.Type == datastoreresource.DataStoreTypeJaeger {
		dataStore.Values.Jaeger = &datastoreresource.GRPCClientSettings{
			TLS: &datastoreresource.TLS{},
		}

		deepcopy.DeepCopy(in.Jaeger, &dataStore.Values.Jaeger)
		deepcopy.DeepCopy(in.Jaeger.Tls, &dataStore.Values.Jaeger.TLS)
	}

	// Tempo
	if dataStore.Type == datastoreresource.DataStoreTypeTempo {
		dataStore.Values.Tempo = &datastoreresource.MultiChannelClientConfig{
			Grpc: &datastoreresource.GRPCClientSettings{
				TLS: &datastoreresource.TLS{},
			},
			Http: &datastoreresource.HttpClientConfig{
				TLS: &datastoreresource.TLS{},
			},
		}

		deepcopy.DeepCopy(in.Tempo, &dataStore.Values.Tempo)
		deepcopy.DeepCopy(in.Tempo.Grpc.Tls, &dataStore.Values.Tempo.Grpc.TLS)
		deepcopy.DeepCopy(in.Tempo.Http.Tls, &dataStore.Values.Tempo.Http.TLS)
	}

	// AWS XRay
	if dataStore.Type == datastoreresource.DataStoreTypeAwsXRay {
		dataStore.Values.AwsXRay = &datastoreresource.AWSXRayConfig{}
		deepcopy.DeepCopy(in.Awsxray, &dataStore.Values.AwsXRay)
	}

	// OpenSearch
	if dataStore.Type == datastoreresource.DataStoreTypeOpenSearch {
		dataStore.Values.OpenSearch = &datastoreresource.ElasticSearchConfig{}
		deepcopy.DeepCopy(in.OpenSearch, &dataStore.Values.OpenSearch)
	}

	// ElasticAPM
	if dataStore.Type == datastoreresource.DataStoreTypeElasticAPM {
		dataStore.Values.OpenSearch = &datastoreresource.ElasticSearchConfig{}
		deepcopy.DeepCopy(in.OpenSearch, &dataStore.Values.ElasticApm)
	}

	// SignalFX
	if dataStore.Type == datastoreresource.DataStoreTypeSignalFX {
		dataStore.Values.SignalFx = &datastoreresource.SignalFXConfig{}
		deepcopy.DeepCopy(in.SignalFx, &dataStore.Values.SignalFx)
	}

	return dataStore
}
