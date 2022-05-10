package test

import (
	"context"

	"github.com/kubeshop/tracetest/app"
	"github.com/kubeshop/tracetest/config"
)

func GetTestingApp() (*app.App, error) {
	ctx := context.Background()
	db, err := GetTestingDatabase("file://../migrations")

	if err != nil {
		return nil, err
	}

	config := config.Config{}
	return app.NewApp(ctx, config, app.WithDB(db))
}
