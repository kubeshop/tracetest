package mappings

import (
	"fmt"
	"time"

	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
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

var dataStoreTypesMapping = map[datastore.DataStoreType]openapi.SupportedDataStores{
	datastore.DataStoreTypeJaeger:           openapi.JAEGER,
	datastore.DataStoreTypeTempo:            openapi.TEMPO,
	datastore.DataStoreTypeOpenSearch:       openapi.OPENSEARCH,
	datastore.DataStoreTypeSignalFX:         openapi.SIGNALFX,
	datastore.DataStoreTypeOTLP:             openapi.OTLP,
	datastore.DataStoreTypeNewRelic:         openapi.NEWRELIC,
	datastore.DataStoreTypeLighStep:         openapi.LIGHTSTEP,
	datastore.DataStoreTypeElasticAPM:       openapi.ELASTICAPM,
	datastore.DataStoreTypeDataDog:          openapi.DATADOG,
	datastore.DataStoreTypeAwsXRay:          openapi.AWSXRAY,
	datastore.DataStoreTypeHoneycomb:        openapi.HONEYCOMB,
	datastore.DataStoreTypeAzureAppInsights: openapi.AZUREAPPINSIGHTS,
	datastore.DataStoreTypeDynatrace:        openapi.DYNATRACE,
	datastore.DataStoreTypeSumoLogic:        openapi.SUMOLOGIC,
	datastore.DataStoreTypeInstana:          openapi.INSTANA,
}

func (m OpenAPI) DataStoreType(in datastore.DataStoreType) openapi.SupportedDataStores {
	dsd, exists := dataStoreTypesMapping[in]
	if !exists {
		// this should only happen during development,
		// so it's more an alert for devs than actual error handling
		panic(fmt.Errorf("trying to convert an undefined datastore.DataStoreType '%s'", in))
	}
	return dsd
}

func (m Model) DataStoreType(in openapi.SupportedDataStores) datastore.DataStoreType {
	for k, v := range dataStoreTypesMapping {
		if v == in {
			return k
		}
	}

	// this should only happen during development,
	// so it's more an alert for devs than actual error handling
	panic(fmt.Errorf("trying to convert an undefined datastore.DataStoreType '%s'", in))

}

func (m Model) DataStore(in openapi.DataStore) datastore.DataStore {
	dataStore := datastore.DataStore{
		ID:        "current",
		Name:      in.Name,
		Type:      m.DataStoreType(in.Type),
		Default:   in.Default,
		CreatedAt: in.CreatedAt.Format(time.RFC3339Nano),
		Values:    datastore.DataStoreValues{},
	}

	// Jaeger
	if dataStore.Type == datastore.DataStoreTypeJaeger {
		dataStore.Values.Jaeger = &datastore.GRPCClientSettings{
			TLS: &datastore.TLS{},
		}

		deepcopy.DeepCopy(in.Jaeger, &dataStore.Values.Jaeger)
		deepcopy.DeepCopy(in.Jaeger.Tls, &dataStore.Values.Jaeger.TLS)
	}

	// Tempo
	if dataStore.Type == datastore.DataStoreTypeTempo {
		dataStore.Values.Tempo = &datastore.MultiChannelClientConfig{
			Grpc: &datastore.GRPCClientSettings{
				TLS: &datastore.TLS{},
			},
			Http: &datastore.HttpClientConfig{
				TLS: &datastore.TLS{},
			},
		}

		deepcopy.DeepCopy(in.Tempo, &dataStore.Values.Tempo)
		deepcopy.DeepCopy(in.Tempo.Grpc.Tls, &dataStore.Values.Tempo.Grpc.TLS)
		deepcopy.DeepCopy(in.Tempo.Http.Tls, &dataStore.Values.Tempo.Http.TLS)
	}

	// AWS XRay
	if dataStore.Type == datastore.DataStoreTypeAwsXRay {
		dataStore.Values.AwsXRay = &datastore.AWSXRayConfig{}
		deepcopy.DeepCopy(in.Awsxray, &dataStore.Values.AwsXRay)
	}

	// OpenSearch
	if dataStore.Type == datastore.DataStoreTypeOpenSearch {
		dataStore.Values.OpenSearch = &datastore.ElasticSearchConfig{}
		deepcopy.DeepCopy(in.Opensearch, &dataStore.Values.OpenSearch)
	}

	// ElasticAPM
	if dataStore.Type == datastore.DataStoreTypeElasticAPM {
		dataStore.Values.ElasticApm = &datastore.ElasticSearchConfig{}
		deepcopy.DeepCopy(in.Elasticapm, &dataStore.Values.ElasticApm)
	}

	// SignalFX
	if dataStore.Type == datastore.DataStoreTypeSignalFX {
		dataStore.Values.SignalFx = &datastore.SignalFXConfig{}
		deepcopy.DeepCopy(in.Signalfx, &dataStore.Values.SignalFx)
	}

	// Azure App Insights
	if dataStore.Type == datastore.DataStoreTypeAzureAppInsights {
		dataStore.Values.AzureAppInsights = &datastore.AzureAppInsightsConfig{}
		deepcopy.DeepCopy(in.Azureappinsights, &dataStore.Values.AzureAppInsights)
	}

	// SumoLogic
	if dataStore.Type == datastore.DataStoreTypeSumoLogic {
		dataStore.Values.SumoLogic = &datastore.SumoLogicConfig{}
		deepcopy.DeepCopy(in.Sumologic, &dataStore.Values.SumoLogic)
	}

	return dataStore
}
