package tracedb_test

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateClient(t *testing.T) {
	cases := []struct {
		name          string
		ds            datastore.DataStore
		expectedType  string
		expectedError error
	}{
		{
			name: "Jaeger",
			ds: datastore.DataStore{
				Type: datastore.DataStoreTypeJaeger,
				Values: datastore.DataStoreValues{
					Jaeger: &datastore.GRPCClientSettings{
						Endpoint: "notexists:123",
					},
				},
			},
			expectedType: "*tracedb.jaegerTraceDB",
		},
		{
			name: "Tempo",
			ds: datastore.DataStore{
				Type: datastore.DataStoreTypeTempo,
				Values: datastore.DataStoreValues{
					Tempo: &datastore.MultiChannelClientConfig{},
				},
			},
			expectedType: "*tracedb.tempoTraceDB",
		},
		{
			name: "ElasticSearch",
			ds: datastore.DataStore{
				Type: datastore.DataStoreTypeElasticAPM,
				Values: datastore.DataStoreValues{
					ElasticApm: &datastore.ElasticSearchConfig{},
				},
			},
			expectedType: "*tracedb.elasticsearchDB",
		},
		{
			name: "OpenSearch",
			ds: datastore.DataStore{
				Type: datastore.DataStoreTypeOpenSearch,
				Values: datastore.DataStoreValues{
					OpenSearch: &datastore.ElasticSearchConfig{},
				},
			},
			expectedType: "*tracedb.opensearchDB",
		},
		{
			name: "SignalFX",
			ds: datastore.DataStore{
				Type: datastore.DataStoreTypeSignalFX,
				Values: datastore.DataStoreValues{
					SignalFx: &datastore.SignalFXConfig{},
				},
			},
			expectedType: "*tracedb.signalfxDB",
		},
		{
			name: "AWSXRay",
			ds: datastore.DataStore{
				Type: datastore.DataStoreTypeAwsXRay,
				Values: datastore.DataStoreValues{
					AwsXRay: &datastore.AWSXRayConfig{},
				},
			},
			expectedType: "*tracedb.awsxrayDB",
		},
		{
			name: "OTLP",
			ds: datastore.DataStore{
				Type:   datastore.DataStoreTypeOTLP,
				Values: datastore.DataStoreValues{},
			},
			expectedType: "*tracedb.OTLPTraceDB",
		},
		{
			name: "NewRelic",
			ds: datastore.DataStore{
				Type:   datastore.DataStoreTypeNewRelic,
				Values: datastore.DataStoreValues{},
			},
			expectedType: "*tracedb.OTLPTraceDB",
		},
		{
			name: "Lightstep",
			ds: datastore.DataStore{
				Type:   datastore.DataStoreTypeLighStep,
				Values: datastore.DataStoreValues{},
			},
			expectedType: "*tracedb.OTLPTraceDB",
		},
		{
			name: "Honeycomb",
			ds: datastore.DataStore{
				Type:   datastore.DataStoreTypeHoneycomb,
				Values: datastore.DataStoreValues{},
			},
			expectedType: "*tracedb.OTLPTraceDB",
		},
		{
			name: "DataDog",
			ds: datastore.DataStore{
				Type:   datastore.DataStoreTypeDataDog,
				Values: datastore.DataStoreValues{},
			},
			expectedType: "*tracedb.OTLPTraceDB",
		},
		{
			name: "Dynatrace",
			ds: datastore.DataStore{
				Type:   datastore.DatastoreTypeDynatrace,
				Values: datastore.DataStoreValues{},
			},
			expectedType: "*tracedb.OTLPTraceDB",
		},
		{
			name:         "EmptyConfig",
			ds:           datastore.DataStore{},
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
