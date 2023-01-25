# Running Tracetest with Jaeger

:::note
[Check out the source code on GitHub here.](https://github.com/kubeshop/tracetest/tree/main/examples/quick-start-jaeger-nodejs) 
:::

## Sample Node.js app with Jaeger, OpenTelemetry and Tracetest

This is a simple quick start on how to configure a Node.js app to use OpenTelemetry instrumentation with traces and Tracetest for enhancing your E2E and integration tests with trace-based testing. The infrastructure will use Jaeger as the trace data store, and OpenTelemetry Collector to receive traces from the Node.js app and send them to Jaeger.

## Prerequisites

You will need [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) installed on your machine to run this quick start app!

## Project Structure

The project is built with Docker Compose. It contains two distinct `docker-compose.yaml` files.

### 1. Node.js app
The `docker-compose.yaml` file and `Dockerfile` in the root directory are for the Node.js app.

### 2. Tracetest
The `docker-compose.yaml` file, `collector.config.yaml`, and `tracetest.config.yaml` in the `tracetest` directory are for the setting up Tracetest, Jaeger, and the OpenTelemetry Collector.

The `tracetest` directory is self-contained and will run all the prerequisites for enabling OpenTelemetry traces and trace-based testing with Tracetest.

### Docker Compose Network
All `services` in the `docker-compose.yaml` are on the same network and will be reachable by hostname from within other services. E.g. `jaeger:14250` in the `collector.config.yaml` will map to the `jaeger` service, where the port `14250` is the port where Jaeger accepts traces. And, `jaeger:16685` in the `tracetest.config.yaml` will map to the `jaeger` service and port `16685` where Tracetest will fetch trace data from Jaeger.

## Node.js app

The Node.js app is a simple Express app, contained in the `app.js` file.

The OpenTelemetry tracing is contained in the `tracing.otel.grpc.js` or `tracing.otel.http.js` files, respectively.
Traces will be sent to the OpenTelemetry Collector.

Here's the content of the `tracing.otel.grpc.js` file:

```js
const opentelemetry = require('@opentelemetry/sdk-node')
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node')
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-grpc');

const sdk = new opentelemetry.NodeSDK({
  traceExporter: new OTLPTraceExporter({ url: 'http://otel-collector:4317' }),
  instrumentations: [getNodeAutoInstrumentations()],
})
sdk.start()
```

Depending on which of these you choose, traces will be sent to either the `grpc` or `http` endpoint.

The hostnames and ports for these are:

- GRPC: `http://otel-collector:4317`
- HTTP: `http://otel-collector:4318/v1/traces`

Enabling the tracer is done by preloading the trace file.

```bash
node -r ./tracing.otel.grpc.js app.js
```

In the `package.json` you will see two npm script for running the respective tracers alongside the `app.js`.

```json
"scripts": {
  "with-grpc-tracer":"node -r ./tracing.otel.grpc.js app.js",
  "with-http-tracer":"node -r ./tracing.otel.http.js app.js"
},
```

To start the server you run this command.

```bash
npm run with-grpc-tracer
# or
npm run with-http-tracer
```

As you can see the `Dockerfile` uses the command above.

```Dockerfile
FROM node:slim
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 8080
CMD [ "npm", "run", "with-grpc-tracer" ]
```

And, the `docker-compose.yaml` contains just one service for the Node.js app.

```yaml
version: '3'
services:
  app:
    image: quick-start-nodejs
    build: .
    ports:
      - "8080:8080"
```

To start it, run this command:

```bash
docker compose build # optional if you haven't already built the image
docker compose up
```

This will start the Node.js app. But, you're not sending the traces anywhere.

Let's fix this by configuring Tracetest and OpenTelemetry Collector.

## Tracetest

The `docker-compose.yaml` in the `tracetest` directory is configured with four services.

- **Postgres** - Postgres is a prerequisite for Tracetest to work. It stores trace data when running the trace-based tests.
- [**OpenTelemetry Collector**](https://opentelemetry.io/docs/collector/) - A vendor-agnostic implementation of how to receive, process and export telemetry data.
- [**Jaeger**](https://www.jaegertracing.io/) - A trace data store.
- [**Tracetest**](https://tracetest.io/) - Trace-based testing that generates end-to-end tests automatically from traces.

```yaml
version: '3'
services:

  tracetest:
    image: kubeshop/tracetest
    platform: linux/amd64
    volumes:
      - ./tracetest/tracetest.config.yaml:/app/config.yaml
    ports:
      - 11633:11633
    depends_on:
      postgres:
        condition: service_healthy
      jaeger:
        condition: service_started
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

  jaeger:
    image: jaegertracing/all-in-one:latest
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--spider", "localhost:16686"]
      interval: 1s
      timeout: 3s
      retries: 60


```

Tracetest depends on Postgres, Jaeger and the OpenTelemetry Collector. Both Tracetest and the OpenTelemetry Collector require config files to be loaded via a volume. The volumes are mapped from the root directory into the `tracetest` directory and the respective config files.

**Why?** To start both the Node.js app and Tracetest we will run this command:

```bash
docker-compose -f docker-compose.yaml -f tracetest/docker-compose.yaml up # add --build if the images are not built already
```

The `tracetest.config.yaml` file contains the basic setup of connecting Tracetest to the Postgres instance, and defining the trace data store and exporter. The data store is set to Jaeger, meaning the traces will be stored in Jaeger and Tracetest will fetch them from Jaeger when running tests. The exporter is set to the OpenTelemetry Collector.

But how does Tracetest fetch traces?

Tracetest uses `jaeger.endpoint:jaeger:16685` to connect to Jaeger and fetch trace data.

```yaml
postgresConnString: "host=postgres user=postgres password=postgres port=5432 sslmode=disable"

poolingConfig:
  maxWaitTimeForTrace: 10m
  retryDelay: 5s

googleAnalytics:
  enabled: true

demo:
  enabled: []

experimentalFeatures: []

telemetry:
  dataStores:
    jaeger:
      type: jaeger
      jaeger:
        endpoint: jaeger:16685
        tls:
          insecure: true

server:
  telemetry:
    dataStore: jaeger
    exporter: collector
    applicationExporter: collector

```

How do traces reach Jaeger?

The `collector.config.yaml` explains that. It receives traces via either `grpc` or `http`. Then, exports them to Jaegers's `model.proto` endpoint `jaeger:14250`.

```yaml
receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:
    timeout: 100ms

exporters:
  logging:
    loglevel: debug
  jaeger:
    endpoint: jaeger:14250
    tls:
      insecure: true

service:
  pipelines:
    traces/1:
      receivers: [otlp]
      processors: [batch]
      exporters: [jaeger]

```

## Run Both the Node.js app and Tracetest

To start both the Node.js app and Tracetest, we will run this command:

```bash
docker-compose -f docker-compose.yaml -f tracetest/docker-compose.yaml up # add --build if the images are not built already
```

This will start your Tracetest instance on `http://localhost:11633/`. Go ahead and open it up.

Start creating tests! Make sure to use the `http://app:8080/` url in your test creation, because your Node.js app and Tracetest are in the same network.

## Learn more

Feel free to check out our [examples in GitHub](https://github.com/kubeshop/tracetest/tree/main/examples), and join our [Discord Community](https://discord.gg/8MtcMrQNbX) for more info!
