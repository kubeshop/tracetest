package tracedb_test

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/configgrpc"
)

func TestCreateClient(t *testing.T) {
	cases := []struct {
		name          string
		ds            model.DataStore
		expectedType  string
		expectedError error
	}{
		{
			name: "Jaeger",
			ds: model.DataStore{
				Type: model.DataStoreTypeJaeger,
				Values: model.DataStoreValues{
					Jaeger: &configgrpc.GRPCClientSettings{},
				},
			},
			expectedType: "*tracedb.jaegerTraceDB",
		},
		{
			name: "Tempo",
			ds: model.DataStore{
				Type: model.DataStoreTypeTempo,
				Values: model.DataStoreValues{
					Tempo: &model.BaseClientConfig{},
				},
			},
			expectedType: "*tracedb.tempoTraceDB",
		},
		{
			name:         "EmptyConfig",
			ds:           model.DataStore{},
			expectedType: "*tracedb.noopTraceDB",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cl := c
			t.Parallel()

			newFn := tracedb.Factory(nil)

			actual, err := newFn(cl.ds)

			require.NoError(t, err)
			assert.Equal(t, cl.expectedType, fmt.Sprintf("%T", actual))
		})
	}

}
