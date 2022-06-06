package testmock

import (
	"context"

	"github.com/kubeshop/tracetest/server/app"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/tracing"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
)

func GetTestingApp(demoApp *DemoApp) (*app.App, error) {
	ctx := context.Background()
	db, err := GetTestingDatabase("file://../migrations")

	if err != nil {
		return nil, err
	}

	config := config.Config{
		JaegerConnectionConfig: &configgrpc.GRPCClientSettings{
			Endpoint: demoApp.JaegerEndpoint(),
			TLSSetting: configtls.TLSClientSetting{
				Insecure: true,
			},
		},
		PoolingConfig: config.PoolingConfig{
			RetryDelay: "5s",
		},
		Telemetry: config.TelemetryConfig{
			Exporters:   []string{"console"},
			ServiceName: "tracetest",
		},
	}

	tracedb, err := tracedb.New(config)
	if err != nil {
		return nil, err
	}

	tracer, err := tracing.NewTracer(ctx, config)
	if err != nil {
		return nil, err
	}

	return app.New(config, db, tracedb, tracer)
}
