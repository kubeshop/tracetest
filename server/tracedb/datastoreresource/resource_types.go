package datastoreresource

import (
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/id"
	"golang.org/x/exp/slices"
)

const ResourceName = "DataStore"
const ResourceNamePlural = "DataStores"

type DataStoreType string

type DataStore struct {
	ID        id.ID           `mapstructure:"id"`
	Name      string          `mapstructure:"name"`
	Type      DataStoreType   `mapstructure:"type"`
	Default   bool            `mapstructure:"default"`
	Values    DataStoreValues `mapstructure:"values,squash"`
	CreatedAt string          `mapstructure:"createdAt"`
}

type DataStoreValues struct {
	AwsXRay    *AWSXRayConfig           `mapstructure:"awsXRay,omitempty"`
	ElasticApm *ElasticSearchConfig     `mapstructure:"elasticAPM,omitempty"`
	Jaeger     *GRPCClientSettings      `mapstructure:"jaeger,omitempty"`
	OpenSearch *ElasticSearchConfig     `mapstructure:"openSearch,omitempty"`
	SignalFx   *SignalFXDataStoreConfig `mapstructure:"signalFX,omitempty"`
	Tempo      *TempoClientConfig       `mapstructure:"tempo,omitempty"`
}

type AWSXRayConfig struct {
	Region          string `mapstructure:"region"`
	AccessKeyID     string `mapstructure:"accessKeyID"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
	SessionToken    string `mapstructure:"sessionToken"`
	UseDefaultAuth  bool   `mapstructure:"useDefaultAuth"`
}

type ElasticSearchConfig struct {
	Addresses          []string `mapstructure:"addresses"`
	Username           string   `mapstructure:"username"`
	Password           string   `mapstructure:"password"`
	Index              string   `mapstructure:"index"`
	Certificate        string   `mapstructure:"certificate"`
	InsecureSkipVerify bool     `mapstructure:"insecureSkipVerify"`
}

type GRPCClientSettings struct {
	Endpoint        string            `mapstructure:"endpoint,omitempty"`
	ReadBufferSize  int               `mapstructure:"readBufferSize,omitempty"`
	WriteBufferSize int               `mapstructure:"writeBufferSize,omitempty"`
	WaitForReady    bool              `mapstructure:"waitForReady,omitempty"`
	Headers         map[string]string `mapstructure:"headers,omitempty"`
	BalancerName    string            `mapstructure:"balancerName,omitempty"`
	Compression     *GRPCCompression  `mapstructure:"compression,omitempty"`
	TLS             *TLS              `mapstructure:"tls,omitempty"`
	Auth            *HttpAuth         `mapstructure:"auth,omitempty"`
}

type HttpAuth struct {
	Type   HttpAuthType    `mapstructure:"type,omitempty"`
	ApiKey *HttpAuthApiKey `mapstructure:"apiKey,omitempty"`
	Basic  *HttpAuthBasic  `mapstructure:"basic,omitempty"`
	Bearer *HttpAuthBearer `mapstructure:"bearer,omitempty"`
}

type HttpAuthApiKey struct {
	Key   string `mapstructure:"key,omitempty"`
	Value string `mapstructure:"value,omitempty"`
	In    string `mapstructure:"in,omitempty"`
}

type HttpAuthBasic struct {
	Username string `mapstructure:"username,omitempty"`
	Password string `mapstructure:"password,omitempty"`
}

type HttpAuthBearer struct {
	Token string `mapstructure:"token,omitempty"`
}

type HttpAuthType string

const (
	HttpAuthTypeApiKey HttpAuthType = "apiKey"
	HttpAuthTypeBasic  HttpAuthType = "basic"
	HttpAuthTypeBearer HttpAuthType = "bearer"
)

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
	Insecure           bool        `mapstructure:"insecure,omitempty"`
	InsecureSkipVerify bool        `mapstructure:"insecureSkipVerify,omitempty"`
	ServerName         string      `mapstructure:"serverName,omitempty"`
	Settings           *TLSSetting `mapstructure:"settings,omitempty"`
}

type TLSSetting struct {
	CAFile     string `mapstructure:"cAFile,omitempty"`
	CertFile   string `mapstructure:"certFile,omitempty"`
	KeyFile    string `mapstructure:"keyFile,omitempty"`
	MinVersion string `mapstructure:"minVersion,omitempty"`
	MaxVersion string `mapstructure:"maxVersion,omitempty"`
}

type TempoClientType string

const (
	TempoClientTypeGRPC TempoClientType = "grpc"
	TempoClientTypeHTTP TempoClientType = "http"
)

type TempoClientConfig struct {
	Type TempoClientType     `mapstructure:"type"`
	Grpc *GRPCClientSettings `mapstructure:"grpc,omitempty"`
	Http *HttpClientConfig   `mapstructure:"http,omitempty"`
}

type HttpClientConfig struct {
	Url     string            `mapstructure:"url"`
	Headers map[string]string `mapstructure:"headers,omitempty"`
	TLS     *TLS              `mapstructure:"tls,omitempty"`
}

type SignalFXDataStoreConfig struct {
	Realm string `mapstructure:"realm"`
	Token string `mapstructure:"token"`
}

const (
	DataStoreTypeJaeger     DataStoreType = "jaeger"
	DataStoreTypeTempo      DataStoreType = "tempo"
	DataStoreTypeOpenSearch DataStoreType = "opensearch"
	DataStoreTypeSignalFX   DataStoreType = "signalfx"
	DataStoreTypeOTLP       DataStoreType = "otlp"
	DataStoreTypeNewRelic   DataStoreType = "newrelic"
	DataStoreTypeLighStep   DataStoreType = "lightstep"
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

func (ds DataStore) Validate() error {
	if ds.Type != "" {
		return fmt.Errorf("data store should have a type")
	}

	if !slices.Contains(validTypes, ds.Type) {
		return fmt.Errorf("unsupported data store type %s", ds.Type)
	}

	if ds.Name != "" {
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

	return nil
}

func (ds DataStore) HasID() bool {
	return ds.ID.String() != ""
}

func (ds DataStore) IsOTLPBasedProvider() bool {
	return slices.Contains(otlpBasedDataStores, ds.Type)
}
