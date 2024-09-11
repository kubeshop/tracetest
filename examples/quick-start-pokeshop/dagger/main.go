// Tracetest pokeshop example demo
package main

import (
	"context"
	"dagger/pokeshop/internal/dagger"
	"strings"
	"time"
)

type Pokeshop struct {
	// +private
	TracetestConfig *dagger.File

	// +private
	TracetestProvision *dagger.File

	// +private
	CollectorConfig *dagger.File

	// +private
	APIBase *dagger.Container
}

func New(
	ctx context.Context,

	// +defaultPath="/examples/quick-start-pokeshop/tracetest.config.yaml"
	ttConfig *dagger.File,

	// +defaultPath="/examples/quick-start-pokeshop/tracetest.provision.yaml"
	ttProvision *dagger.File,

	// +defaultPath="/examples/quick-start-pokeshop/collector.config.yaml"
	ttCollectorConfig *dagger.File,

) *Pokeshop {
	p := &Pokeshop{
		TracetestConfig:    ttConfig,
		TracetestProvision: ttProvision,
		CollectorConfig:    ttCollectorConfig,
	}
	apiBase := dag.Container().
		From("kubeshop/demo-pokemon-api:latest").
		WithServiceBinding("postgres", p.postgres(ctx)).
		WithServiceBinding("cache", p.cache(ctx)).
		WithServiceBinding("queue", p.queue(ctx)).
		WithServiceBinding("otel-collector", p.otelCollector(ctx)).
		WithEnvVariable("REDIS_URL", "cache").
		WithEnvVariable("DATABASE_URL", "postgresql://postgres:postgres@postgres:5432/postgres?schema=public").
		WithEnvVariable("RABBITMQ_HOST", "queue").
		WithEnvVariable("POKE_API_BASE_URL", "https://pokeapi.co/api/v2").
		WithEnvVariable("COLLECTOR_ENDPOINT", "http://otel-collector:4317")

	p.APIBase = apiBase
	return p
}

// returns the Pokehop API service
func (m *Pokeshop) API(ctx context.Context) (*dagger.Service, error) {

	_, err := m.worker(ctx).Start(ctx)
	if err != nil {
		return nil, err
	}

	// TODO is rpc service necessary?
	_, err = m.rpc(ctx).Start(ctx)
	if err != nil {
		return nil, err
	}

	apiSvc := m.APIBase.
		WithEnvVariable("NPM_RUN_COMMAND", "api").
		WithExposedPort(8081).
		AsService()

	return apiSvc, nil
}

func (m *Pokeshop) Tracetest(ctx context.Context, testFile *dagger.File) (string, error) {
	apiSvc, err := m.API(ctx)
	if err != nil {
		return "", err
	}

	apiSvc, err = apiSvc.Start(ctx)

	if err != nil {
		return "", err
	}

	apiSvcEndpoint, err := apiSvc.Endpoint(ctx)

	if err != nil {
		return "", err
	}

	c, err := testFile.Contents(ctx)

	if err != nil {
		return "", err
	}

	updatedFile := strings.ReplaceAll(c, "demo-api:8081", apiSvcEndpoint)

	return dag.Tracetest(dagger.TracetestOpts{Server: "http://tracetest:11633"}).
		WithService("tracetest", m.TracetestSvc(ctx)).
		Cli().
		WithNewFile("test.yaml", updatedFile).
		WithEnvVariable("CACHE", time.Now().String()).
		WithExec([]string{"run", "test", "--file", "test.yaml"},
			dagger.ContainerWithExecOpts{UseEntrypoint: true}).Stdout(ctx)

}

func (m *Pokeshop) worker(ctx context.Context) *dagger.Service {
	return m.APIBase.
		WithEnvVariable("NPM_RUN_COMMAND", "worker").
		AsService()
}

func (m *Pokeshop) rpc(ctx context.Context) *dagger.Service {
	return m.APIBase.
		WithEnvVariable("NPM_RUN_COMMAND", "rpc").
		AsService()
}

func (m *Pokeshop) cache(ctx context.Context) *dagger.Service {
	return dag.Container().
		From("redis:6").AsService()
}

func (m *Pokeshop) postgres(ctx context.Context) *dagger.Service {
	return dag.Container().
		From("postgres:14").
		WithEnvVariable("POSTGRES_PASSWORD", "postgres").
		WithEnvVariable("POSTGRES_USER", "postgres").
		WithExposedPort(5432).
		AsService()
}

func (m *Pokeshop) queue(ctx context.Context) *dagger.Service {
	return dag.Container().
		From("rabbitmq:3.8-management").AsService()
}

func (m *Pokeshop) TracetestSvc(ctx context.Context) *dagger.Service {
	return dag.Tracetest().Agent().
		WithExposedPort(11633).
		WithServiceBinding("postgres", m.postgres(ctx)).
		WithMountedFile("/app/tracetest.yaml", m.TracetestConfig).
		WithMountedFile("/app/provision.yaml", m.TracetestProvision).
		WithExec([]string{"--provisioning-file", "/app/provision.yaml"}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		AsService()
}

func (m *Pokeshop) otelCollector(ctx context.Context) *dagger.Service {
	c := dag.Container().
		From("otel/opentelemetry-collector:0.59.0").
		WithServiceBinding("tracetest", m.TracetestSvc(ctx))
	return c.WithoutUser().
		WithFile("/etc/otelcol/config.yaml", m.CollectorConfig).
		WithExec([]string{"/otelcol", "--config", "/etc/otelcol/config.yaml"}).
		WithExposedPort(4317).AsService()
}
