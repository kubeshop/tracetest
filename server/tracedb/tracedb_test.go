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
			name: "ElasticSearch",
			ds: model.DataStore{
				Type: model.DataStoreTypeElasticAPM,
				Values: model.DataStoreValues{
					ElasticApm: &model.ElasticSearchDataStoreConfig{},
				},
			},
			expectedType: "*tracedb.elasticsearchDB",
		},
		{
			name: "OpenSearch",
			ds: model.DataStore{
				Type: model.DataStoreTypeOpenSearch,
				Values: model.DataStoreValues{
					OpenSearch: &model.ElasticSearchDataStoreConfig{},
				},
			},
			expectedType: "*tracedb.opensearchDB",
		},
		{
			name: "SignalFX",
			ds: model.DataStore{
				Type: model.DataStoreTypeSignalFX,
				Values: model.DataStoreValues{
					SignalFx: &model.SignalFXDataStoreConfig{},
				},
			},
			expectedType: "*tracedb.signalfxDB",
		},
		{
			name: "AWSXRay",
			ds: model.DataStore{
				Type: model.DataStoreTypeAwsXRay,
				Values: model.DataStoreValues{
					AwsXRay: &model.AWSXRayDataStoreConfig{},
				},
			},
			expectedType: "*tracedb.awsxrayDB",
		},
		{
			name: "OTLP",
			ds: model.DataStore{
				Type:   model.DataStoreTypeOTLP,
				Values: model.DataStoreValues{},
			},
			expectedType: "*tracedb.OTLPTraceDB",
		},
		{
			name: "NewRelic",
			ds: model.DataStore{
				Type:   model.DataStoreTypeNewRelic,
				Values: model.DataStoreValues{},
			},
			expectedType: "*tracedb.OTLPTraceDB",
		},
		{
			name: "Lightstep",
			ds: model.DataStore{
				Type:   model.DataStoreTypeLighStep,
				Values: model.DataStoreValues{},
			},
			expectedType: "*tracedb.OTLPTraceDB",
		},
		{
			name: "DataDog",
			ds: model.DataStore{
				Type:   model.DataStoreTypeDataDog,
				Values: model.DataStoreValues{},
			},
			expectedType: "*tracedb.OTLPTraceDB",
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
