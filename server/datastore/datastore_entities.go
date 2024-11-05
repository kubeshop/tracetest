package datastore

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"golang.org/x/exp/slices"
)

const ResourceName = "DataStore"
const ResourceNamePlural = "DataStores"

type DataStoreType string

type DataStore struct {
	ID        id.ID           `json:"id"`
	Name      string          `json:"name"`
	Type      DataStoreType   `json:"type"`
	Default   bool            `json:"default"`
	Values    DataStoreValues `json:"values"`
	CreatedAt string          `json:"createdAt"`
}

type DataStoreValues struct {
	AwsXRay          *AWSXRayConfig            `json:"awsxray,omitempty"`
	ElasticApm       *ElasticSearchConfig      `json:"elasticapm,omitempty"`
	Jaeger           *GRPCClientSettings       `json:"jaeger,omitempty"`
	OpenSearch       *ElasticSearchConfig      `json:"opensearch,omitempty"`
	SignalFx         *SignalFXConfig           `json:"signalfx,omitempty"`
	Tempo            *MultiChannelClientConfig `json:"tempo,omitempty"`
	AzureAppInsights *AzureAppInsightsConfig   `json:"azureappinsights,omitempty"`
	SumoLogic        *SumoLogicConfig          `json:"sumologic"`
}

type AWSXRayConfig struct {
	Region          string `json:"region"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
	UseDefaultAuth  bool   `json:"useDefaultAuth"`
}

type ElasticSearchConfig struct {
	Addresses          []string `json:"addresses"`
	Username           string   `json:"username"`
	Password           string   `json:"password"`
	Index              string   `json:"index"`
	Certificate        string   `json:"certificate"`
	InsecureSkipVerify bool     `json:"insecureSkipVerify"`
}

type GRPCClientSettings struct {
	Endpoint        string            `json:"endpoint,omitempty"`
	ReadBufferSize  int               `json:"readBufferSize,omitempty"`
	WriteBufferSize int               `json:"writeBufferSize,omitempty"`
	WaitForReady    bool              `json:"waitForReady,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	BalancerName    string            `json:"balancerName,omitempty"`
	Compression     GRPCCompression   `json:"compression,omitempty"`
	TLS             *TLS              `json:"tls,omitempty"`
}

type GRPCCompression string

const (
	GRPCCompressionGZip    GRPCCompression = "gzip"
	GRPCCompressionZlib    GRPCCompression = "zlib"
	GRPCCompressionDeflate GRPCCompression = "deflate"
	GRPCCompressionSnappy  GRPCCompression = "snappy"
	GRPCCompressionZstd    GRPCCompression = "zstd"
	GRPCCompressionNone    GRPCCompression = "none"
)

type TLS struct {
	Insecure           bool        `json:"insecure,omitempty"`
	InsecureSkipVerify bool        `json:"insecureSkipVerify,omitempty"`
	ServerName         string      `json:"serverName,omitempty"`
	Settings           *TLSSetting `json:"settings,omitempty"`
}

type TLSSetting struct {
	CAFile     string `json:"cAFile,omitempty"`
	CertFile   string `json:"certFile,omitempty"`
	KeyFile    string `json:"keyFile,omitempty"`
	MinVersion string `json:"minVersion,omitempty"`
	MaxVersion string `json:"maxVersion,omitempty"`
}

type MultiChannelClientType string

const (
	MultiChannelClientTypeGRPC MultiChannelClientType = "grpc"
	MultiChannelClientTypeHTTP MultiChannelClientType = "http"
)

type MultiChannelClientConfig struct {
	Type MultiChannelClientType `json:"type"`
	Grpc *GRPCClientSettings    `json:"grpc,omitempty"`
	Http *HttpClientConfig      `json:"http,omitempty"`
}

type HttpClientConfig struct {
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	TLS     *TLS              `json:"tls,omitempty"`
}

type SignalFXConfig struct {
	Realm string `json:"realm"`
	Token string `json:"token"`
}

type ConnectionTypes string

const (
	ConnectionTypesDirect    ConnectionTypes = "direct"
	ConnectionTypesCollector ConnectionTypes = "collector"
)

type AzureAppInsightsConfig struct {
	UseAzureActiveDirectoryAuth bool            `json:"useAzureActiveDirectoryAuth"`
	AccessToken                 string          `json:"accessToken"`
	ResourceArmId               string          `json:"resourceArmId"`
	ConnectionType              ConnectionTypes `json:"connectionType"`
}

type SumoLogicConfig struct {
	URL       string `json:"url"`
	AccessID  string `json:"accessID"`
	AccessKey string `json:"accessKey"`
}

const (
	DataStoreTypeAgent            DataStoreType = "agent"
	DataStoreTypeAwsXRay          DataStoreType = "awsxray"
	DataStoreTypeAzureAppInsights DataStoreType = "azureappinsights"
	DataStoreTypeDash0            DataStoreType = "dash0"
	DataStoreTypeDataDog          DataStoreType = "datadog"
	DataStoreTypeDynatrace        DataStoreType = "dynatrace"
	DataStoreTypeElasticAPM       DataStoreType = "elasticapm"
	DataStoreTypeHoneycomb        DataStoreType = "honeycomb"
	DataStoreTypeInstana          DataStoreType = "instana"
	DataStoreTypeJaeger           DataStoreType = "jaeger"
	DataStoreTypeLighStep         DataStoreType = "lightstep"
	DataStoreTypeNewRelic         DataStoreType = "newrelic"
	DataStoreTypeOpenSearch       DataStoreType = "opensearch"
	DataStoreTypeOTLP             DataStoreType = "otlp"
	DataStoreTypeSignalFX         DataStoreType = "signalfx"
	DataStoreTypeSignoz           DataStoreType = "signoz"
	DataStoreTypeSumoLogic        DataStoreType = "sumologic"
	DataStoreTypeTempo            DataStoreType = "tempo"
)

var validTypes = []DataStoreType{
	DataStoreTypeAgent,
	DataStoreTypeAwsXRay,
	DataStoreTypeAzureAppInsights,
	DataStoreTypeDash0,
	DataStoreTypeDataDog,
	DataStoreTypeDynatrace,
	DataStoreTypeElasticAPM,
	DataStoreTypeHoneycomb,
	DataStoreTypeInstana,
	DataStoreTypeJaeger,
	DataStoreTypeLighStep,
	DataStoreTypeNewRelic,
	DataStoreTypeOpenSearch,
	DataStoreTypeOTLP,
	DataStoreTypeSignalFX,
	DataStoreTypeSignoz,
	DataStoreTypeSumoLogic,
	DataStoreTypeTempo,
}

var otlpBasedDataStores = []DataStoreType{
	DataStoreTypeAgent,
	DataStoreTypeDash0,
	DataStoreTypeDataDog,
	DataStoreTypeDynatrace,
	DataStoreTypeHoneycomb,
	DataStoreTypeInstana,
	DataStoreTypeLighStep,
	DataStoreTypeNewRelic,
	DataStoreTypeOTLP,
	DataStoreTypeSignoz,
}

func (ds DataStore) Validate() error {
	if ds.Type == "" {
		return fmt.Errorf("data store should have a type")
	}

	if !slices.Contains(validTypes, ds.Type) {
		return fmt.Errorf("unsupported data store type %s", ds.Type)
	}

	if ds.Name == "" {
		return fmt.Errorf("data store should have a name")
	}

	if ds.CreatedAt != "" {
		if _, err := time.Parse(time.RFC3339Nano, ds.CreatedAt); err != nil {
			return fmt.Errorf("data store should have the createdAt field in a valid format")
		}
	}

	if ds.Type == DataStoreTypeAwsXRay && ds.Values.AwsXRay == nil {
		return fmt.Errorf("data store should have AWSXRay config values set up")
	}

	if ds.Type == DataStoreTypeElasticAPM && ds.Values.ElasticApm == nil {
		return fmt.Errorf("data store should have ElasticApm config values set up")
	}

	if ds.Type == DataStoreTypeJaeger && ds.Values.Jaeger == nil {
		return fmt.Errorf("data store should have Jaeger config values set up")
	}

	if ds.Type == DataStoreTypeOpenSearch && ds.Values.OpenSearch == nil {
		return fmt.Errorf("data store should have OpenSearch config values set up")
	}

	if ds.Type == DataStoreTypeSignalFX && ds.Values.SignalFx == nil {
		return fmt.Errorf("data store should have SignalFx config values set up")
	}

	if ds.Type == DataStoreTypeTempo && ds.Values.Tempo == nil {
		return fmt.Errorf("data store should have Tempo config values set up")
	}

	if ds.Type == DataStoreTypeAzureAppInsights && ds.Values.AzureAppInsights == nil {
		return fmt.Errorf("data store should have Azure Application Insights config values set up")
	}

	if ds.Type == DataStoreTypeSumoLogic && ds.Values.SumoLogic == nil {
		return fmt.Errorf("data store should have Sumo Logic config values set up")
	}

	return nil
}

func (ds DataStore) HasID() bool {
	return ds.ID.String() != ""
}

func (ds DataStore) GetID() id.ID {
	return ds.ID
}

func (ds DataStore) IsOTLPBasedProvider() bool {
	if ds.Type == DataStoreTypeAzureAppInsights {
		return ds.Values.AzureAppInsights.ConnectionType == ConnectionTypesCollector
	}

	return slices.Contains(otlpBasedDataStores, ds.Type)
}

type squashedDataStore struct {
	ID                     id.ID                     `json:"id"`
	Name                   string                    `json:"name"`
	Type                   DataStoreType             `json:"type"`
	Default                bool                      `json:"default"`
	CreatedAt              string                    `json:"createdAt"`
	AwsXRay                *AWSXRayConfig            `json:"awsxray,omitempty"`
	ElasticApm             *ElasticSearchConfig      `json:"elasticapm,omitempty"`
	Jaeger                 *GRPCClientSettings       `json:"jaeger,omitempty"`
	OpenSearch             *ElasticSearchConfig      `json:"opensearch,omitempty"`
	SumoLogic              *SumoLogicConfig          `json:"sumologic,omitempty"`
	SignalFx               *SignalFXConfig           `json:"signalfx,omitempty"`
	Tempo                  *MultiChannelClientConfig `json:"tempo,omitempty"`
	AzureAppInsightsConfig *AzureAppInsightsConfig   `json:"azureappinsights,omitempty"`
}

func (d squashedDataStore) populate(dataStore *DataStore) {
	if dataStore == nil {
		return
	}

	dataStore.ID = d.ID
	dataStore.Name = d.Name
	dataStore.Type = d.Type
	dataStore.Default = d.Default
	dataStore.CreatedAt = d.CreatedAt
	dataStore.Values = DataStoreValues{
		AwsXRay:          d.AwsXRay,
		ElasticApm:       d.ElasticApm,
		Jaeger:           d.Jaeger,
		OpenSearch:       d.OpenSearch,
		SignalFx:         d.SignalFx,
		SumoLogic:        d.SumoLogic,
		Tempo:            d.Tempo,
		AzureAppInsights: d.AzureAppInsightsConfig,
	}
}

func (d DataStore) MarshalJSON() ([]byte, error) {
	squashedObject := d.squashed()
	return json.Marshal(squashedObject)
}

func (d DataStore) MarshalYAML() ([]byte, error) {
	squashedObject := d.squashed()
	return yaml.Marshal(squashedObject)
}

func (d *DataStore) UnmarshalJSON(input []byte) error {
	squashedObject := squashedDataStore{}
	err := json.Unmarshal(input, &squashedObject)
	if err != nil {
		return err
	}

	squashedObject.populate(d)

	return nil
}

func (d *DataStore) UnmarshalYAML(input []byte) error {
	squashedObject := squashedDataStore{}
	err := yaml.Unmarshal(input, &squashedObject)
	if err != nil {
		return err
	}

	squashedObject.populate(d)

	return nil
}

func (d DataStore) squashed() squashedDataStore {
	return squashedDataStore{
		ID:                     d.ID,
		Name:                   d.Name,
		Type:                   d.Type,
		Default:                d.Default,
		CreatedAt:              d.CreatedAt,
		AwsXRay:                d.Values.AwsXRay,
		ElasticApm:             d.Values.ElasticApm,
		Jaeger:                 d.Values.Jaeger,
		OpenSearch:             d.Values.OpenSearch,
		SignalFx:               d.Values.SignalFx,
		Tempo:                  d.Values.Tempo,
		AzureAppInsightsConfig: d.Values.AzureAppInsights,
		SumoLogic:              d.Values.SumoLogic,
	}
}
