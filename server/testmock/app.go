package testmock

import (
	"github.com/kubeshop/tracetest/server/app"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/logger"
)

type TestingAppOption func(*config.Config)

func WithServerPrefix(prefix string) TestingAppOption {
	return func(cfg *config.Config) {
		cfg.Set("server.pathPrefix", prefix)
	}
}

func WithHttpPort(port int) TestingAppOption {
	return func(cfg *config.Config) {
		cfg.Set("server.httpPort", port)
	}
}

func GetTestingApp(options ...TestingAppOption) (*app.App, error) {
	cfg, _ := config.New(nil, logger.Default())
	for _, option := range options {
		option(cfg)
	}
	err := ConfigureDB(cfg)
	if err != nil {
		panic(err)
	}

	return app.New(cfg)
}
