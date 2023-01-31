# Quick Start - Node.js app with OpenTelemetry and Tracetest that includes manual instrumentation

This is a simple quick start on how to configure a Node.js app to use OpenTelemetry instrumentation with traces, and Tracetest for enhancing your e2e and integration tests with trace-based testing. This example includes manual instrumentation and a sample bookstore array that simulates fetching data from a database.

## Prerequisites

You will need [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) installed on your machine to run this quick start app!

## Project structure

The project is built with Docker Compose. It contains two distinct `docker-compose.yaml` files.

### 1. Node.js app
The `docker-compose.yaml` file and `Dockerfile` in the root directory are for the Node.js app.

### 2. Tracetest
The `docker-compose.yaml` file, `collector.config.yaml`, and `tracetest.config.yaml` in the `tracetest` directory are for the setting up Tracetest and the OpenTelemetry Collector.

The `tracetest` directory is self-contained and will run all the prerequisites for enabling OpenTelemetry traces and trace-based testing with Tracetest.

### Docker Compose Network
All `services` in the `docker-compose.yaml` are on the same network and will be reachable by hostname from within other services. E.g. `tracetest:21321` in the `collector.config.yaml` will map to the `tracetest` service, where the port `21321` is the port where Tracetest accepts traces.

## Node.js app

The Node.js app is a simple Express app with two microservices, contained in the `app.js` and `availability.js` files.

The OpenTelemetry tracing is contained in the `tracing.otel.grpc.js` or `tracing.otel.http.js` files, respectively.
Traces will be sent to the OpenTelemetry Collector.

Here's the content of the `tracing.otel.grpc.js` file:

```js
const opentelemetry = require("@opentelemetry/sdk-node")
const { getNodeAutoInstrumentations } = require("@opentelemetry/auto-instrumentations-node")
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-grpc')
const { Resource } = require("@opentelemetry/resources")
const { SemanticResourceAttributes } = require("@opentelemetry/semantic-conventions")
const { NodeTracerProvider } = require("@opentelemetry/sdk-trace-node")
const { BatchSpanProcessor } = require("@opentelemetry/sdk-trace-base")

const resource = Resource.default().merge(
  new Resource({
    [SemanticResourceAttributes.SERVICE_NAME]: "quick-start-nodejs-manual-instrumentation",
    [SemanticResourceAttributes.SERVICE_VERSION]: "0.0.1",
  })
)

const provider = new NodeTracerProvider({ resource: resource })
const exporter = new OTLPTraceExporter({ url: 'http://otel-collector:4317' })
const processor = new BatchSpanProcessor(exporter)
provider.addSpanProcessor(processor)
provider.register()

const sdk = new opentelemetry.NodeSDK({
  traceExporter: exporter,
  instrumentations: [getNodeAutoInstrumentations()],
  serviceName: 'quick-start-nodejs-manual-instrumentation'
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

In the `package.json` you will see two npm script for running the respective tracers alongside the `app.js` or `availability.js`.

```json
"scripts": {
  "app-with-grpc-tracer": "node -r ./tracing.otel.grpc.js app.js",
  "app-with-http-tracer": "node -r ./tracing.otel.http.js app.js",
  "availability-with-grpc-tracer": "node -r ./tracing.otel.grpc.js availability.js",
  "availability-with-http-tracer": "node -r ./tracing.otel.http.js availability.js"
},
```

To start the `app.js` Express server you run this command.

```bash
npm run app-with-grpc-tracer
# or
npm run app-with-http-tracer
```

To start the `availability.js` Express server you run this command.

```bash
npm run availability-with-grpc-tracer
# or
npm run availability-with-http-tracer
```

As you can see the `Dockerfile` does not have a `CMD` section.

```Dockerfile
FROM node:slim
WORKDIR /usr/src/app/
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 8080
```

Instead, the `docker-compose.yaml` contains the `CMD` section for both services.

```yaml
version: '3'
services:
  app:
    image: quick-start-nodejs
    extra_hosts:
      - "host.docker.internal:host-gateway"
    build: .
    command: npm run app-with-grpc-tracer
    ports:
      - "8080:8080"
  availability:
    image: quick-start-nodejs-availability
    extra_hosts:
      - "host.docker.internal:host-gateway"
    build: .
    command: npm run availability-with-grpc-tracer
    ports:
      - "8080"
```

To start it, run this command:

```bash
docker compose build # optional if you haven't already built the image
docker compose up
```

This will start the Node.js services. But, you're not sending the traces anywhere.

Let's fix this by configuring Tracetest and OpenTelemetry Collector.

## Tracetest

The `docker-compose.yaml` in the `tracetest` directory is configured with three services.

- **Postgres** - Postgres is a prerequisite for Tracetest to work. It stores trace data when running the trace-based tests.
- [**OpenTelemetry Collector**](https://opentelemetry.io/docs/collector/) - A vendor-agnostic implementation of how to receive, process and export telemetry data.
- [**Tracetest**](https://tracetest.io/) - Trace-based testing that generates end-to-end tests automatically from traces.

```yaml
version: '3'
services:

  tracetest:
    image: kubeshop/tracetest:v0.9.3
    platform: linux/amd64
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

Tracetest depends on both Postgres and the OpenTelemetry Collector. Both Tracetest and the OpenTelemetry Collector require config files to be loaded via a volume. The volumes are mapped from the root directory into the `tracetest` directory and the respective config files.

**Why?** To start both the Node.js services and Tracetest we will run this command:

```bash
docker-compose -f docker-compose.yaml -f tracetest/docker-compose.yaml up # add --build if the images are not built already
```

The `tracetest.config.yaml` file contains the basic setup of connecting Tracetest to the Postgres instance, and defining the trace data store and exporter. The data store is set to OTLP meaning the traces will be stored in Tracetest itself. The exporter is set to the OpenTelemetry Collector.

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

But how are traces sent to Tracetest?

The `collector.config.yaml` explains that. It receives traces via either `grpc` or `http`. Then, exports them to Tracetest's otlp endpoint `tracetest:21321`.

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
  otlp/1:
    endpoint: tracetest:21321 # Send traces to Tracetest. Read more in docs here. (TODO ADD LINK)
    tls:
      insecure: true

service:
  pipelines:
    traces/1:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/1]

```

**Important!** Take a closer look at the sampling configs in both the `collector.config.yaml` and `tracetest.config.yaml`. They both set sampling to 100%. This is crucial when running trace-based e2e and integration tests.

## Run both the Node.js app and Tracetest

To start both the Node.js services and Tracetest we will run this command:

```bash
docker-compose -f docker-compose.yaml -f tracetest/docker-compose.yaml up # add --build if the images are not built already
```

This will start your Tracetest instance on `http://localhost:11633/`. Go ahead and open it up.

[Start creating tests in the Web UI](https://docs.tracetest.io/web-ui/creating-tests)! Make sure to use the `http://app:8080/books` URL in your test creation, because your Node.js app and Tracetest are in the same network.

Here's a sample of a failed test run, which happens if you add this assertion:

```
attr:books.list.count = 4
```

![](https://res.cloudinary.com/djwdcmwdz/image/upload/v1673808310/screely-1673808287031_sol4it.png)

It fails because of the `books.length` is equal to `3`.

## Run Tracetest tests with the Tracetest CLI

First, [install the CLI](https://docs.tracetest.io/getting-started/installation#install-the-tracetest-cli).
Then, configure the CLI:

```bash
tracetest configure --endpoint http://localhost:11633 --analytics
```

Once configure, you can run a test against the Tracetest instance via the terminal.

Check out the `test-api.yaml` file.

```yaml
type: Test
spec:
  id: W656Q0c4g
  name: Books List
  description: List of books
  trigger:
    type: http
    httpRequest:
      url: http://app:8080/books
      method: GET
      headers:
      - key: Content-Type
        value: application/json
  specs:
  - selector: span[tracetest.span.type="http" name="GET /books" http.target="/books" http.method="GET"]
    assertions:
    - attr:http.status_code = 200
  - selector: span[tracetest.span.type="general" name="Books List"]
    assertions:
    - attr:books.list.count = 4
```

This file defines the a test the same way you would through the Web UI.

To run the test, run this command in the terminal:

```bash
tracetest test run -d ./test-api.yaml -w
```

This test will fail just like the sample above due to the `attr:books.list.count = 4` assertion.

```
✘ http://app:8080 (http://localhost:11633/test/W656Q0c4g/run/5/test)
	✔ span[tracetest.span.type="http" name="GET /books" http.target="/books" http.method="GET"]
		✔ #994c63e0ea35e632
			✔ attr:http.status_code = 200 (200)
	✘ span[tracetest.span.type="general" name="Books List"]
		✘ #5ab1856c32b0d5c8
			✘ attr:books.list.count = 4 (3) (http://localhost:11633/test/W656Q0c4g/run/5/test?selectedAssertion=1&selectedSpan=5ab1856c32b0d5c8)
```

The tests will pass if you change the assertion to:

```
attr:books.list.count = 3
```

Feel free to check out our [docs](https://docs.tracetest.io/), and join our [Discord Community](https://discord.gg/8MtcMrQNbX) for more info!
