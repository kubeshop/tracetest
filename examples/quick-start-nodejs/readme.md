# Quick Start - Node.js app with OpenTelemetry and Tracetest

> [Read the detailed recipe for setting up OpenTelemetry Collector with Tractest in our documentation.](https://docs.tracetest.io/examples-tutorials/recipes/running-tracetest-without-a-trace-data-store)

This is a simple quick start on how to configure a Node.js app to use OpenTelemetry instrumentation with traces, and Tracetest for enhancing your e2e and integration tests with trace-based testing.

## Steps to run Tracetest

### 1. Start the Tracetest Agent locally

```bash
tracetest start
```

Once started, Tracetest Agent will:

- Expose OTLP ports 4317 (gRPC) and 4318 (HTTP) for trace ingestion.
- Be able to trigger test runs in the environment where it is running.
- Be able to connect to a trace data store that is not accessible outside of your environment. Eg. a Jaeger instance running in the cluster without an ingress controller.

### 2. Start Node.js App

You can run the example with Docker or locally.

#### Docker Compose

```bash
docker compose up --build
```

#### Locally

```bash
npm i 
OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://localhost:4317 npm run with-grpc-tracer
```

### 3. Run tests

Create and run a test against `http://localhost:8080` on [`https://app.tracetest.io/`](https://app.tracetest.io/). View the `./test-api.yaml` for reference.

## Steps to run Tracetest Core

### 1. Start Node.js App and Tracetest Core with Docker Compose

```bash
docker compose -f ./docker-compose.yaml -f ./tracetest/docker-compose.yaml up --build
```

### 2. Run tests

Once started, you will need to make sure to trigger tests with correct service names since both the Node.js app and Tracetest Core are running in the same Docker Network. In this example the Node.js app would be at `http://app:8080`. View the `./test-api.yaml` for reference.

---

Feel free to check out the [docs](https://docs.tracetest.io/), and join our [Slack Community](https://dub.sh/tracetest-community) for more info!
