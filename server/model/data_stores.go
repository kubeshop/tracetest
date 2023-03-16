package model

import (
	"fmt"
	"time"

	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
	"golang.org/x/exp/slices"
)

type (
	DataStore struct {
		ID        string
		Name      string
		Type      DataStoreType
		IsDefault bool
		Values    DataStoreValues
		CreatedAt time.Time
	}

	DataStoreValues struct {
		Jaeger     *configgrpc.GRPCClientSettings
		Tempo      *BaseClientConfig
		OpenSearch *ElasticSearchDataStoreConfig
		ElasticApm *ElasticSearchDataStoreConfig
		SignalFx   *SignalFXDataStoreConfig
		AwsXRay    *AWSXRayDataStoreConfig
	}

	BaseClientConfig struct {
		Type string
		Grpc configgrpc.GRPCClientSettings
		Http HttpClientConfig
	}

	HttpClientConfig struct {
		Url        string
		Headers    map[string]string
		TLSSetting configtls.TLSClientSetting
	}

	OTELCollectorConfig struct {
		Endpoint string
	}

	ElasticSearchDataStoreConfig struct {
		Addresses          []string
		Username           string
		Password           string
		Index              string
		Certificate        string
		InsecureSkipVerify bool
	}

	SignalFXDataStoreConfig struct {
		Realm string
		Token string
	}

	AWSXRayDataStoreConfig struct {
		Region          string
		AccessKeyID     string
		SecretAccessKey string
		SessionToken    string
	}
)

func (ds DataStore) IsZero() bool {
	return ds.Type == ""
}

type DataStoreType string

const (
	DataStoreTypeJaeger     DataStoreType = "jaeger"
	DataStoreTypeTempo      DataStoreType = "tempo"
	DataStoreTypeOpenSearch DataStoreType = "opensearch"
	DataStoreTypeSignalFX   DataStoreType = "signalfx"
	DataStoreTypeOTLP       DataStoreType = "otlp"
	DataStoreTypeNewRelic   DataStoreType = "newrelic"
	DataStoreTypeLighStep   DataStoreType = "lighstep"
	DataStoreTypeElasticAPM DataStoreType = "elasticapm"
	DataStoreTypeDataDog    DataStoreType = "datadog"
	DataStoreTypeAwsXRay    DataStoreType = "awsxray"
)

var validTypes = []DataStoreType{
	DataStoreTypeJaeger,
	DataStoreTypeTempo,
	DataStoreTypeOpenSearch,
	DataStoreTypeSignalFX,
	DataStoreTypeOTLP,
	DataStoreTypeNewRelic,
	DataStoreTypeLighStep,
	DataStoreTypeElasticAPM,
	DataStoreTypeDataDog,
	DataStoreTypeAwsXRay,
}

var otlpBasedDataStores = []DataStoreType{
	DataStoreTypeOTLP,
	DataStoreTypeNewRelic,
	DataStoreTypeLighStep,
	DataStoreTypeDataDog,
}

func (ds DataStore) HasID() bool {
	return ds.ID != ""
}

func (ds DataStore) Validate() error {
	if !slices.Contains(validTypes, ds.Type) {
		return fmt.Errorf("unsupported data store")
	}

	return nil
}

func (ds DataStore) IsOTLPBasedProvider() bool {
	return slices.Contains(otlpBasedDataStores, ds.Type)
}
