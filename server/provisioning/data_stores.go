package provisioning

import (
	"fmt"

	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
)

var ErrInvalidTraceDBProvider = fmt.Errorf("invalid traceDB provider")

type (
	dataStore struct {
		Type       string                         `mapstructure:"type"`
		Jaeger     *configgrpc.GRPCClientSettings `mapstructure:"jaeger"`
		Tempo      *baseClientConfig              `mapstructure:"tempo"`
		OpenSearch *elasticSearch                 `mapstructure:"opensearch"`
		SignalFX   *signalFX                      `mapstructure:"signalfx"`
		ElasticApm *elasticSearch                 `mapstructure:"elasticapm"`
		AwsXRay    *aWSXRayDataStoreConfig        `mapstructure:"awsxray"`
	}

	baseClientConfig struct {
		Type string                        `mapstructure:"type"`
		Grpc configgrpc.GRPCClientSettings `mapstructure:"grpc"`
		Http httpClientConfig              `mapstructure:"http"`
	}

	httpClientConfig struct {
		Url        string                     `mapstructure:"url"`
		Headers    map[string]string          `mapstructure:"headers"`
		TLSSetting configtls.TLSClientSetting `mapstructure:"tls"`
	}

	elasticSearch struct {
		Addresses          []string
		Username           string
		Password           string
		Index              string
		Certificate        string
		InsecureSkipVerify bool
	}

	signalFX struct {
		Realm string
		Token string
	}

	aWSXRayDataStoreConfig struct {
		Region          string
		AccessKeyID     string
		SecretAccessKey string
		SessionToken    string
	}
)

func (ds dataStore) model() model.DataStore {
	m := model.DataStore{
		Name: ds.Type,
		Type: model.DataStoreType(ds.Type),
	}

	if ds.Jaeger != nil {
		deepcopy.DeepCopy(ds.Jaeger, &m.Values.Jaeger)
	}
	if ds.Tempo != nil {
		deepcopy.DeepCopy(ds.Tempo, &m.Values.Tempo)
	}
	if ds.OpenSearch != nil {
		deepcopy.DeepCopy(ds.OpenSearch, &m.Values.OpenSearch)
	}
	if ds.ElasticApm != nil {
		deepcopy.DeepCopy(ds.ElasticApm, &m.Values.ElasticApm)
	}
	if ds.SignalFX != nil {
		deepcopy.DeepCopy(ds.SignalFX, &m.Values.SignalFx)
	}

	if ds.AwsXRay != nil {
		deepcopy.DeepCopy(ds.AwsXRay, &m.Values.AwsXRay)
	}

	return m
}
