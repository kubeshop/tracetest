package model

import (
	"fmt"
	"strings"
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

func (ds DataStore) Slug() string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(ds.Name), " ", "-"))
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
		Type:   cType,
		Values: DataStoreValues{},
	}
	switch dsc.Type {
	case jaeger:
		ds.Type = openapi.JAEGER
		ds.Values.Jaeger = &dsc.Jaeger
	case tempo:
		ds.Type = openapi.TEMPO
		ds.Values.Tempo = &dsc.Tempo
	case opensearch:
		ds.Type = openapi.OPEN_SEARCH
		ds.Values.OpenSearch = &dsc.OpenSearch
	case signalfx:
		ds.Type = openapi.SIGNAL_FX
		ds.Values.SignalFx = &dsc.SignalFX
	case otlp:
		ds.Type = openapi.OTLP
	}

	return ds
}
