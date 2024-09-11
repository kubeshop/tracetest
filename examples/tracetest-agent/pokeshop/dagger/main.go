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
	APIBase *dagger.Container

	// +private
	TraceTestAPIKey *dagger.Secret

	// +private
	TracetestEnvironment string

	// +private
	TracetestOrganization string
}

func New(
	ctx context.Context,
	// +optinal
	ttAPIKey *dagger.Secret,

	// +optional
	ttEnvironment string,

	// +optional
	ttOrganization string,
) *Pokeshop {
	p := &Pokeshop{
		TracetestEnvironment:  ttEnvironment,
		TraceTestAPIKey:       ttAPIKey,
		TracetestOrganization: ttOrganization,
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

func (m *Pokeshop) Tracetest(ctx context.Context, file *dagger.File) (string, error) {
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

	c, err := file.Contents(ctx)

	if err != nil {
		return "", err
	}

	updatedFile := strings.ReplaceAll(c, "localhost:8081", apiSvcEndpoint)

	return dag.Tracetest(dagger.TracetestOpts{
		APIKey:       m.TraceTestAPIKey,
		Environment:  m.TracetestEnvironment,
		Organization: m.TracetestOrganization,
	}).Cli().
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

const otelCollectorConfig = `
receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:
    timeout: 100ms

  # Data sources: traces
  probabilistic_sampler:
    hash_seed: 22
    sampling_percentage: 100

exporters:
  otlp/tracetestagent:
    endpoint: tracetest:4317
    tls:
      insecure: true

service:
  pipelines:
    traces/tracetestagent:
      receivers: [otlp]
      processors: [probabilistic_sampler, batch]
      exporters: [otlp/tracetestagent]
`

func (m *Pokeshop) otelCollector(ctx context.Context) *dagger.Service {
	c := dag.Container().
		From("otel/opentelemetry-collector:0.54.0")

	if m.TraceTestAPIKey != nil && m.TracetestEnvironment != "" {
		c = c.WithServiceBinding("tracetest", dag.Tracetest().Agent(m.TraceTestAPIKey, m.TracetestEnvironment).
			AsService())
	}

	return c.WithoutUser().
		WithNewFile("/etc/otelcol/config.yaml", otelCollectorConfig).
		WithExec([]string{"/otelcol", "--config", "/etc/otelcol/config.yaml"}).
		WithExposedPort(4317).AsService()
}
