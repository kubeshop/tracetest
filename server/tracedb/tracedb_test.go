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
				JaegerConnectionConfig: &configgrpc.GRPCClientSettings{},
			},
			expectedType: "*tracedb.jaegerTraceDB",
		},
		{
			name: "Tempo",
			config: config.Config{
				TempoConnectionConfig: &configgrpc.GRPCClientSettings{},
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

			actual, err := tracedb.New(cl.config)

			if cl.expectedType != "" {
				require.NoError(t, err)
				assert.Equal(t, cl.expectedType, fmt.Sprintf("%T", actual))
			} else if cl.expectedError != nil {
				assert.ErrorIs(t, err, cl.expectedError)
			}
		})
	}

}
