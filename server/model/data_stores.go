package model

import (
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/openapi"
	"go.opentelemetry.io/collector/config/configgrpc"
	"strings"
	"time"
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

func (e DataStore) Slug() string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(e.Name), " ", "-"))
}
