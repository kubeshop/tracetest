package tracedb_test

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/tracedb/datastoreresource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateClient(t *testing.T) {
	cases := []struct {
		name          string
		ds            datastoreresource.DataStore
		expectedType  string
		expectedError error
	}{
		{
			name: "Jaeger",
			ds: datastoreresource.DataStore{
				Type: datastoreresource.DataStoreTypeJaeger,
				Values: datastoreresource.DataStoreValues{
					Jaeger: &datastoreresource.GRPCClientSettings{
						Endpoint: "notexists:123",
					},
				},
			},
			expectedType: "*tracedb.jaegerTraceDB",
		},
		{
			name: "Tempo",
			ds: datastoreresource.DataStore{
				Type: datastoreresource.DataStoreTypeTempo,
				Values: datastoreresource.DataStoreValues{
					Tempo: &datastoreresource.MultiChannelClientConfig{},
				},
			},
			expectedType: "*tracedb.tempoTraceDB",
		},
		{
			name: "ElasticSearch",
			ds: datastoreresource.DataStore{
				Type: datastoreresource.DataStoreTypeElasticAPM,
				Values: datastoreresource.DataStoreValues{
					ElasticApm: &datastoreresource.ElasticSearchConfig{},
				},
			},
			expectedType: "*tracedb.elasticsearchDB",
		},
		{
			name: "OpenSearch",
			ds: datastoreresource.DataStore{
				Type: datastoreresource.DataStoreTypeOpenSearch,
				Values: datastoreresource.DataStoreValues{
					OpenSearch: &datastoreresource.ElasticSearchConfig{},
				},
			},
			expectedType: "*tracedb.opensearchDB",
		},
		{
			name: "SignalFX",
			ds: datastoreresource.DataStore{
				Type: datastoreresource.DataStoreTypeSignalFX,
				Values: datastoreresource.DataStoreValues{
					SignalFx: &datastoreresource.SignalFXConfig{},
				},
			},
			expectedType: "*tracedb.signalfxDB",
		},
		{
			name: "AWSXRay",
			ds: datastoreresource.DataStore{
				Type: datastoreresource.DataStoreTypeAwsXRay,
				Values: datastoreresource.DataStoreValues{
					AwsXRay: &datastoreresource.AWSXRayConfig{},
				},
			},
			expectedType: "*tracedb.awsxrayDB",
		},
		{
			name: "OTLP",
			ds: datastoreresource.DataStore{
				Type:   datastoreresource.DataStoreTypeOTLP,
				Values: datastoreresource.DataStoreValues{},
			},
			expectedType: "*tracedb.OTLPTraceDB",
		},
		{
			name: "NewRelic",
			ds: datastoreresource.DataStore{
				Type:   datastoreresource.DataStoreTypeNewRelic,
				Values: datastoreresource.DataStoreValues{},
			},
			expectedType: "*tracedb.OTLPTraceDB",
		},
		{
			name: "Lightstep",
			ds: datastoreresource.DataStore{
				Type:   datastoreresource.DataStoreTypeLighStep,
				Values: datastoreresource.DataStoreValues{},
			},
			expectedType: "*tracedb.OTLPTraceDB",
		},
		{
			name: "Honeycomb",
			ds: datastoreresource.DataStore{
				Type:   datastoreresource.DataStoreTypeHoneycomb,
				Values: datastoreresource.DataStoreValues{},
			},
			expectedType: "*tracedb.OTLPTraceDB",
		},
		{
			name: "DataDog",
			ds: datastoreresource.DataStore{
				Type:   datastoreresource.DataStoreTypeDataDog,
				Values: datastoreresource.DataStoreValues{},
			},
			expectedType: "*tracedb.OTLPTraceDB",
		},
		{
			name:         "EmptyConfig",
			ds:           datastoreresource.DataStore{},
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
