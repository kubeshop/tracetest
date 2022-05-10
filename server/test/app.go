package test

import (
	"context"

	"github.com/kubeshop/tracetest/app"
	"github.com/kubeshop/tracetest/config"
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
	}

	return app.NewApp(ctx, config, app.WithDB(db))
}
