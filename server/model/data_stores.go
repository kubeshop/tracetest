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
		Tempo      *configgrpc.GRPCClientSettings
		OpenSearch *config.OpensearchDataStoreConfig
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
)

func DataStoreFromConfig(dsc config.TracingBackendDataStoreConfig) DataStore {
	var cType openapi.SupportedDataStores
	ds := DataStore{
		Name:   string(cType),
		Type:   cType,
		Values: DataStoreValues{},
	}
	switch dsc.Type {
	case jaeger:
		ds.Type = openapi.JAEGER
	case tempo:
		ds.Type = openapi.TEMPO
	case opensearch:
		ds.Type = openapi.OPEN_SEARCH
	case signalfx:
		ds.Type = openapi.SIGNAL_FX
	case otlp:
		ds.Type = openapi.OTLP
	}
	ds.Values.Jaeger = &dsc.Jaeger
	ds.Values.Tempo = &dsc.Tempo
	ds.Values.OpenSearch = &dsc.OpenSearch
	ds.Values.SignalFx = &dsc.SignalFX

	return ds
}
