package tracedb_test

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/configgrpc"
)

func TestCreateClient(t *testing.T) {
	cases := []struct {
		name          string
		config        config.Config
		expectedType  string
		expectedError error
	}{
		{
			name: "Jaeger",
			config: config.Config{
				Telemetry: config.Telemetry{
					DataStores: map[string]config.TracingBackendDataStoreConfig{
						"jaeger": {
							Type:   tracedb.JAEGER_BACKEND,
							Jaeger: configgrpc.GRPCClientSettings{},
						},
					},
				},
				Server: config.ServerConfig{
					Telemetry: config.ServerTelemetryConfig{
						DataStore: "jaeger",
					},
				},
			},
			expectedType: "*tracedb.jaegerTraceDB",
		},
		{
			name: "Tempo",
			config: config.Config{
				Telemetry: config.Telemetry{
					DataStores: map[string]config.TracingBackendDataStoreConfig{
						"tempo": {
							Type:   tracedb.TEMPO_BACKEND,
							Jaeger: configgrpc.GRPCClientSettings{},
						},
					},
				},
				Server: config.ServerConfig{
					Telemetry: config.ServerTelemetryConfig{
						DataStore: "tempo",
					},
				},
			},
			expectedType: "*tracedb.tempoTraceDB",
		},
		{
			name:          "InvalidConfig",
			config:        config.Config{},
			expectedError: tracedb.ErrInvalidTraceDBProvider,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cl := c
			t.Parallel()

			actual, err := tracedb.New(cl.config, nil)

			if cl.expectedType != "" {
				require.NoError(t, err)
				assert.Equal(t, cl.expectedType, fmt.Sprintf("%T", actual))
			} else if cl.expectedError != nil {
				assert.ErrorIs(t, err, cl.expectedError)
			}
		})
	}

}
