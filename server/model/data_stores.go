package model

import (
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/openapi"
	"go.opentelemetry.io/collector/config/configgrpc"
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

	DataStoreValues struct {
		Jaeger     *configgrpc.GRPCClientSettings
		Tempo      *config.BaseClientConfig
		OpenSearch *config.ElasticSearchDataStoreConfig
		ElasticApm *config.ElasticSearchDataStoreConfig
		SignalFx   *config.SignalFXDataStoreConfig
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
)

func DataStoreFromConfig(dsc config.TracingBackendDataStoreConfig) DataStore {
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
