package testmock

import (
	"github.com/kubeshop/tracetest/server/app"
	"github.com/kubeshop/tracetest/server/config"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
)

type TestingAppOption func(config *config.Config)

func WithServerPrefix(prefix string) TestingAppOption {
	return func(config *config.Config) {
		config.Server.PathPrefix = prefix
	}
}

func WithHttpPort(port int) TestingAppOption {
	return func(config *config.Config) {
		config.Server.HttpPort = port
	}
}

func GetTestingApp(options ...TestingAppOption) (*app.App, error) {
	db, err := GetRawTestingDatabase()

	if err != nil {
		return nil, err
	}

	config := config.Config{
		Telemetry: config.Telemetry{
			DataStores: map[string]config.TracingBackendDataStoreConfig{
				"jaeger": {
					Type: "jaeger",
					Jaeger: configgrpc.GRPCClientSettings{
						Endpoint:   "",
						TLSSetting: configtls.TLSClientSetting{Insecure: true},
					},
				},
			},
		},
		Server: config.ServerConfig{
			Telemetry: config.ServerTelemetryConfig{
				DataStore: "jaeger",
			},
		},
		PoolingConfig: config.PoolingConfig{
			RetryDelay: "5s",
		},
	}

	for _, option := range options {
		option(&config)
	}

	return app.New(app.Config{
		Config:     config,
		Migrations: "file://../migrations",
	}, db)
}
