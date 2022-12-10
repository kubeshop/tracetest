# Getting Started with Tracetest Using Docker Compose

This guide will help you get Tracetest running using Docker Compose.

## Prerequisites

:::info
Make sure you have [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) installed.
:::

:::caution
Postgres is a prerequisite for Tracetest to work. It stores Tracetest's trace data. Make sure to have a Postgres service.
:::

:::info
In this quick start, OpenTelemetry Collector is used to send traces directly to Tracetest. If you have an existing trace data source, [read here](../configuration/overview.md).
:::

## 1. Create a Docker Compose config

Create a `docker-compose.yaml` file:

```yaml
version: '3'
services:
  tracetest:
    image: kubeshop/tracetest
    volumes:
      - ./tracetest.config.yaml:/app/config.yaml
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
      - ./collector.config.yaml:/otel-local-config.yaml
```

## 2. Create a Tracetest Config File

Create a `tracetest.config.yaml` file:

```yaml
postgresConnString: "host=postgres user=postgres password=postgres port=5432 sslmode=disable"

poolingConfig:
  maxWaitTimeForTrace: 10s
  retryDelay: 1s

googleAnalytics:
  enabled: true

telemetry:
  dataStores:
    otlp:
      type: otlp
  exporters:
    collector:
      serviceName: tracetest
      sampling: 100 # 100%
      exporter:
        type: collector
        collector:
          endpoint: otel-collector:4317

server:
  telemetry:
    dataStore: otlp
    exporter: collector
    applicationExporter: collector
```

The `postgresConnString` will configure Tracetest to connect to the Postgres service. The `telemetry.dataStores` defines that the trace data store will be through `otlp` because it's expecting to receive traces from the OpenTelemetry Collector.

## 3. Create an OpenTelemetry Collector Config File

Create a `collector.config.yaml` file:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  tail_sampling:
    decision_wait: 5s
    policies:
      - name: tracetest-spans
        type: trace_state
        trace_state: { key: tracetest, values: ["true"] }
  batch:
    timeout: 100ms

exporters:
  logging:
    loglevel: debug
  otlp/1:
    endpoint: tracetest:21321 # This sends traces to Tracetest
    tls:
      insecure: true

service:
  pipelines:
    traces/1:
      receivers: [otlp]
      processors: [tail_sampling, batch]
      exporters: [otlp/1]
```

## 4. Start Docker Compose

```bash
docker compose up
```

```bash title="Condensed expected output from the Tracetest container:"
Starting tracetest ...
...
2022/11/28 18:24:09 HTTP Server started
...
```

## 5. Open the Tracetest Web UI

Open your browser on [`http://localhost:11633`](http://localhost:11633).

Create a [test](../web-ui/creating-tests.md).

:::info
Running a test against `localhost` will resolve as the 127.0.0.1 inside the Tracetest container. To run tests against apps running on your local machine, add them to the same network and use service name mapping instead. Example: Instead of running an app on `localhost:8080`, add it to your Docker Compose file, connect it to the same network as your Tracetest service, and use `service-name:8080` in the URL field when creating an app.

You can reach services running on your local machine using:

- Linux (docker version < 20.10.0): `172.17.0.1:8080`
- MacOS (docker version >= 18.03) and Linux (docker version >= 20.10.0): `host.docker.internal:8080`
:::
