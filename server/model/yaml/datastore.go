package yaml

type DataStore struct {
	Id         string             `mapstructure:"id"`
	Name       string             `mapstructure:"name"`
	Type       string             `mapstructure:"type"`
	IsDefault  bool               `mapstructure:"isDefault"`
	Jaeger     GRPCClientSettings `mapstructure:"jaeger"`
	Tempo      BaseClientConfig   `mapstructure:"tempo"`
	OpenSearch OpenSearch         `mapstructure:"openSearch"`
	SignalFx   SignalFX           `mapstructure:"signalFx"`
	AwsXRay    AwsXRay            `mapstructure:"awsxray"`
}

type GRPCClientSettings struct {
	Endpoint        string             `mapstructure:"endpoint"`
	ReadBufferSize  float32            `mapstructure:"readBufferSize"`
	WriteBufferSize float32            `mapstructure:"writeBufferSize"`
	WaitForReady    bool               `mapstructure:"waitForReady"`
	Headers         map[string]string  `mapstructure:"headers"`
	BalancerName    string             `mapstructure:"balancerName"`
	Compression     string             `mapstructure:"compression"`
	Tls             TLS                `mapstructure:"tls"`
	Auth            HTTPAuthentication `mapstructure:"auth"`
}

type BaseClientConfig struct {
	Type string             `mapstructure:"type"`
	Grpc GRPCClientSettings `mapstructure:"grpc"`
	Http HttpClientConfig   `mapstructure:"http"`
}

type HttpClientConfig struct {
	Url     string            `yaml:",omitempty" mapstructure:"url"`
	Headers map[string]string `yaml:",omitempty" mapstructure:"headers"`
	Tls     TLS               `yaml:",omitempty" mapstructure:"tls"`
}

type TLS struct {
	Insecure           bool       `mapstructure:"insecure"`
	InsecureSkipVerify bool       `mapstructure:"insecureSkipVerify"`
	ServerName         string     `mapstructure:"serverName"`
	Settings           TLSSetting `mapstructure:"settings"`
}

type TLSSetting struct {
	CAFile     string `mapstructure:"cAFile"`
	CertFile   string `mapstructure:"certFile"`
	KeyFile    string `mapstructure:"keyFile"`
	MinVersion string `mapstructure:"minVersion"`
	MaxVersion string `mapstructure:"maxVersion"`
}

type SignalFX struct {
	Realm string `mapstructure:"realm"`
	Token string `mapstructure:"token"`
}

type AwsXRay struct {
	Region          string `mapstructure:"region"`
	AccessKeyId     string `mapstructure:"accessKeyId"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
}

type OpenSearch struct {
	Addresses []string `mapstructure:"addresses"`
	Username  string   `mapstructure:"username"`
	Password  string   `mapstructure:"password"`
	Index     string   `mapstructure:"index"`
}
