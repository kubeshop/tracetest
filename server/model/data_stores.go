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
	openapi.OPENSEARCH,
	openapi.TEMPO,
	openapi.SIGNALFX,
	openapi.OTLP,
}

func (ds DataStore) Validate() error {
	if !slices.Contains(validTypes, ds.Type) {
		return fmt.Errorf("unsupported data store")
	}

	return nil
}

func DataStoreFromConfig(dsc config.TracingBackendDataStoreConfig) DataStore {
	return DataStore{
		Type: openapi.SupportedDataStores(dsc.Type),
		Values: DataStoreValues{
			Jaeger:     &dsc.Jaeger,
			Tempo:      &dsc.Tempo,
			OpenSearch: &dsc.OpenSearch,
			SignalFx:   &dsc.SignalFX,
		},
	}
}
