package provisioning_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/provisioning"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
)

var cases = []struct {
	name   string
	dsType model.DataStoreType
	file   string
	values model.DataStoreValues
}{
	{
		name:   "JaegerGRPC",
		file:   "./testdata/jaeger_grpc.yaml",
		dsType: model.DataStoreTypeJaeger,
		values: model.DataStoreValues{
			Jaeger: &configgrpc.GRPCClientSettings{
				Endpoint:   "jaeger-query:16685",
				TLSSetting: configtls.TLSClientSetting{Insecure: true},
			},
		},
	},
	{
		name:   "TempoGRPC",
		file:   "./testdata/tempo_grpc.yaml",
		dsType: model.DataStoreTypeTempo,
		values: model.DataStoreValues{
			Tempo: &model.BaseClientConfig{
				Grpc: configgrpc.GRPCClientSettings{
					Endpoint:   "tempo:9095",
					TLSSetting: configtls.TLSClientSetting{Insecure: true},
				},
			},
		},
	},
	{
		name:   "TempoHTTP",
		file:   "./testdata/tempo_http.yaml",
		dsType: model.DataStoreTypeTempo,
		values: model.DataStoreValues{
			Tempo: &model.BaseClientConfig{
				Http: model.HttpClientConfig{
					Url:        "tempo:80",
					TLSSetting: configtls.TLSClientSetting{Insecure: true},
				},
			},
		},
	},
	{
		name:   "OpenSearch",
		file:   "./testdata/opensearch.yaml",
		dsType: model.DataStoreTypeOpenSearch,
		values: model.DataStoreValues{
			OpenSearch: &model.ElasticSearchDataStoreConfig{
				Addresses: []string{"http://opensearch:9200"},
				Index:     "traces",
			},
		},
	},
	{
		name:   "SignalFX",
		file:   "./testdata/signalfx.yaml",
		dsType: model.DataStoreTypeSignalFX,
		values: model.DataStoreValues{
			SignalFx: &model.SignalFXDataStoreConfig{
				Token: "thetoken",
				Realm: "us1",
			},
		},
	},
	{
		name:   "ElasitcAPM",
		file:   "./testdata/elastic_apm.yaml",
		dsType: model.DataStoreTypeElasticAPM,
		values: model.DataStoreValues{
			ElasticApm: &model.ElasticSearchDataStoreConfig{
				Addresses:          []string{"https://es01:9200"},
				Username:           "elastic",
				Password:           "changeme",
				Index:              "traces-apm-default",
				InsecureSkipVerify: true,
			},
		},
	},
}

func TestDataStore(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Run("FromFile", func(t *testing.T) {
				c := c
				t.Parallel()

				expectedDataStore := model.DataStore{
					IsDefault: true,
					Name:      string(c.dsType),
					Type:      c.dsType,
					Values:    c.values,
				}

				mockRepo := &testdb.MockRepository{}
				provisioner := provisioning.New(mockRepo)

				mockRepo.
					On("CreateDataStore", expectedDataStore).
					Return(expectedDataStore, nil)

				err := provisioner.FromFile(c.file)
				assert.NoError(t, err)
				mockRepo.AssertExpectations(t)
			})
		})
	}
}
