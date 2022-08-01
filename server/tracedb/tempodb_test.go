package tracedb_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/otel/trace"
)

func TestGetTraceByIdentification(t *testing.T) {
	t.Skip("TODO: docker-compose tempo")
	db, err := tracedb.New(config.Config{
		Telemetry: config.Telemetry{
			DataStores: map[string]config.TracingBackendDataStoreConfig{
				"tempo": {
					Type: tracedb.TEMPO_BACKEND,
					Tempo: configgrpc.GRPCClientSettings{
						Endpoint:   "localhost:9095",
						TLSSetting: configtls.TLSClientSetting{Insecure: true},
					},
				},
			},
		},
		Server: config.ServerConfig{
			Telemetry: config.ServerTelemetryConfig{
				DataStore: "tempo",
			},
		},
	})
	require.NoError(t, err)

	defer db.Close()
	traceId, _ := trace.TraceIDFromHex("0194fdc2fa2ffcc041d3ff12045b73c9")
	traceIdentification := traces.TraceIdentification{
		TraceID: traceId,
	}

	trace, err := db.GetTraceByIdentification(context.Background(), traceIdentification)
	assert.NoError(t, err)

	assert.NotEmpty(t, trace.ID)
}
