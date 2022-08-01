package tracedb_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/otel/trace"
)

func TestJaegerGetTraceByIdentification(t *testing.T) {
	t.Skip("TODO: docker-compose jaeger")
	db, err := tracedb.New(config.Config{
		Telemetry: config.Telemetry{
			DataStores: map[string]config.TracingBackendDataStoreConfig{
				"jaeger": {
					Type: tracedb.JAEGER_BACKEND,
					Jaeger: configgrpc.GRPCClientSettings{
						Endpoint:   "localhost:16685",
						TLSSetting: configtls.TLSClientSetting{Insecure: true},
					},
				},
			},
		},
		Server: config.ServerConfig{
			Telemetry: config.ServerTelemetryConfig{
				DataStore: "jaeger",
			},
		},
	})
	assert.NoError(t, err)

	defer db.Close()
	traceId, _ := trace.TraceIDFromHex("0194fdc2fa2ffcc041d3ff12045b73c8")
	traceIdentification := traces.TraceIdentification{
		TraceID: traceId,
	}
	trace, err := db.GetTraceByIdentification(context.Background(), traceIdentification)
	assert.NoError(t, err)

	assert.NotEmpty(t, trace.ID)
}
