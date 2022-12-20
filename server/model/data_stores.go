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
	switch dsc.Type {
	case jaeger:
		cType = openapi.JAEGER
	case tempo:
		cType = openapi.TEMPO
	case opensearch:
		cType = openapi.OPEN_SEARCH
	case signalfx:
		cType = openapi.SIGNAL_FX
	case otlp:
		cType = openapi.OTLP
	}

	return DataStore{
		Type: cType,
		Values: DataStoreValues{
			Jaeger:     &dsc.Jaeger,
			Tempo:      &dsc.Tempo,
			OpenSearch: &dsc.OpenSearch,
			SignalFx:   &dsc.SignalFX,
		},
	}
}
