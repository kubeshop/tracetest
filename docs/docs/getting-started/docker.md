# Getting Started with Tracetest using Docker

Make sure you have Docker and Docker Compose installed.

## TL;DR

Create a `docker-compose.yaml` file:

:::caution
Postgres is a prerequisite for Tracetest to work. It stores Tracetest's trace data. Make sure to have a Postgres service.
:::

:::info
OpenTelemetry Collector is used to send traces directly to Tracetest. If you have an existing trace data source, [read here](../configuration/overview.md).
:::

```yaml
version: '3'
services:
  tracetest:
    image: kubeshop/tracetest
    volumes:
      - ./tracetest/tracetest.config.yaml:/app/config.yaml
    ports:
      - 11633:11633
    depends_on:
      postgres:
        condition: service_healthy
      otel-collector:
        condition: service_started
    healthcheck:
      test: ["CMD", "wget", "--spider", "localhost:11633"]
      interval: 1s
      timeout: 3s
      retries: 60

  postgres:
    image: postgres:14
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_USER: postgres
    healthcheck:
      test: pg_isready -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"
      interval: 1s
      timeout: 5s
      retries: 60

  otel-collector:
    image: otel/opentelemetry-collector-contrib:0.59.0
    command:
      - "--config"
      - "/otel-local-config.yaml"
    volumes:
      - ./tracetest/collector.config.yaml:/otel-local-config.yaml
```



## How It Works

Getting started with observability and OpenTelemetry can be complex and overwhelming. It involves different interconected services working together.

Our CLI offers an **install wizard** that helps with the process. It helps not only to setup tracetest itself, but all the tools needed 
to observe your application.

Use the install wizard to install Tracetest locally using Docker Compose or to a local or remote Kubernetes cluster.
It installs all the tools required to set up the desired environment and creates all the configurations, tailored to your case.

## CLI Installation