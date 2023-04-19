package testmock

import (
	"github.com/kubeshop/tracetest/server/app"
	"github.com/kubeshop/tracetest/server/config"
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
	cfg, _ := config.New()
	for _, option := range options {
		option(cfg)
	}

	ConfigureDB(cfg)

	return app.New(cfg)
}

func ConfigureDB(cfg *config.Config) {
	db := getTestDatabaseEnvironment()

	cfg.Set("postgres.host", db.container.Host)
	cfg.Set("postgres.user", "tracetest")
	cfg.Set("postgres.password", "tracetest")
	cfg.Set("postgres.dbname", "postgres")
	cfg.Set("postgres.port", db.container.DefaultPort())
}
