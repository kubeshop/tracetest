package model

import (
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/openapi"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
	"golang.org/x/exp/slices"
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

	GRPCClientSettings struct {
		configgrpc.GRPCClientSettings
	}

	DataStoreValues struct {
		Jaeger     *GRPCClientSettings
		Tempo      *BaseClientConfig
		OpenSearch *ElasticSearchDataStoreConfig
		ElasticApm *ElasticSearchDataStoreConfig
		SignalFx   *SignalFXDataStoreConfig
	}

	TracingBackendDataStoreConfig struct {
		Type       string
		Jaeger     GRPCClientSettings
		Tempo      BaseClientConfig
		OpenSearch ElasticSearchDataStoreConfig
		SignalFX   SignalFXDataStoreConfig
		ElasticApm ElasticSearchDataStoreConfig
	}

	BaseClientConfig struct {
		Type string
		Grpc GRPCClientSettings
		Http HttpClientConfig
	}

	HttpClientConfig struct {
		Url        string
		Headers    map[string]string
		TLSSetting configtls.TLSClientSetting
	}

	OTELCollectorConfig struct {
		Endpoint string
	}

	ElasticSearchDataStoreConfig struct {
		Addresses          []string
		Username           string
		Password           string
		Index              string
		Certificate        string
		InsecureSkipVerify bool
	}

	SignalFXDataStoreConfig struct {
		Realm string
		Token string
	}
)

func (ds DataStore) IsZero() bool {
	return ds.Type == ""
}

var validTypes = []openapi.SupportedDataStores{
	openapi.JAEGER,
	openapi.OPEN_SEARCH,
	openapi.TEMPO,
	openapi.SIGNAL_FX,
	openapi.OTLP,
	openapi.ELASTIC_APM,
}

func (ds DataStore) Validate() error {
	if !slices.Contains(validTypes, ds.Type) {
		return fmt.Errorf("unsupported data store")
	}

	return nil
}

const (
	jaeger     string = "jaeger"
	tempo      string = "tempo"
	opensearch string = "opensearch"
	signalfx   string = "signalfx"
	otlp       string = "otlp"
	newrelic   string = "newrelic"
	lighstep   string = "lighstep"
	elasticapm string = "elasticapm"
	datadog    string = "datadog"
)

func DataStoreFromConfig(dsc TracingBackendDataStoreConfig) DataStore {
	var cType openapi.SupportedDataStores
	ds := DataStore{
		Name:   dsc.Type,
		Type:   cType,
		Values: DataStoreValues{},
	}

	switch dsc.Type {
	case jaeger:
		ds.Type = openapi.JAEGER
	case tempo:
		ds.Type = openapi.TEMPO
	case elasticapm:
		ds.Type = openapi.ELASTIC_APM
	case opensearch:
		ds.Type = openapi.OPEN_SEARCH
	case signalfx:
		ds.Type = openapi.SIGNAL_FX
	case newrelic:
		ds.Type = openapi.NEW_RELIC
	case lighstep:
		ds.Type = openapi.LIGHTSTEP
	case datadog:
		ds.Type = openapi.DATADOG
	case otlp:
		ds.Type = openapi.OTLP
	}
	ds.Values.Jaeger = &dsc.Jaeger
	ds.Values.Tempo = &dsc.Tempo
	ds.Values.OpenSearch = &dsc.OpenSearch
	ds.Values.SignalFx = &dsc.SignalFX
	ds.Values.ElasticApm = &dsc.ElasticApm

	return ds
}
