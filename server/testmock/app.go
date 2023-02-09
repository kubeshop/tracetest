package testmock

import (
	"github.com/kubeshop/tracetest/server/app"
	"github.com/kubeshop/tracetest/server/config"
)

type TestingAppOption func(config *config.Config)

func WithServerPrefix(prefix string) TestingAppOption {
	return func(config *config.Config) {
		config.SetServerPathPrefix(prefix)
	}
}

func WithHttpPort(port int) TestingAppOption {
	return func(config *config.Config) {
		config.SetServerPort(port)
	}
}

func GetTestingApp(options ...TestingAppOption) (*app.App, error) {
	db, err := GetRawTestingDatabase()

	if err != nil {
		return nil, err
	}

	cfg := config.New()
	for _, option := range options {
		option(cfg)
	}

	return app.New(app.Config{
		Config:     cfg,
		Migrations: "file://../migrations",
	}, db)
}
