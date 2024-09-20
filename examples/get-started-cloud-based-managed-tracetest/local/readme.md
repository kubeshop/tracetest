# Running Node.js with OpenTelemetry and Cloud-based Managed Tracetest Locally

## Install Dependencies

```bash
npm install express \
  @opentelemetry/sdk-node \
  @opentelemetry/auto-instrumentations-node \
  @opentelemetry/exporter-trace-otlp-grpc \
  @opentelemetry/sdk-trace-node
```

## Run Node.js

```bash
npm start
```

## Run Tracetest Agent

1. [Sign in to Tracetest](https://app.tracetest.io/).
2. [Create an Organization](https://docs.tracetest.io/concepts/organizations).
3. Retrieve your [Tracetest Organization API Key/Token and Environment ID](https://app.tracetest.io/retrieve-token).
4. [Install the Tracetest CLI](https://docs.tracetest.io/install/cli#install-the-tracetest-cli).
5. Start [Tracetest Agent](https://docs.tracetest.io/concepts/agent) as a standalone process with the `<TRACETEST_API_KEY>` and `<TRACETEST_ENVIRONMENT_ID>` from step 2.

```bash
tracetest start --api-key <TRACETEST_API_KEY> --environment <TRACETEST_ENVIRONMENT_ID>
```

Tracetest Agent will run on gRPC and HTTP ports.

- `http://localhost:4317` — gRPC
- `http://localhost:4318/v1/traces` — HTTP

## Configure Trace Ingestion for Localhost

Go to the Trace Ingestion tab in the settings, select OpenTelemetry, and enable the toggle.

## Run Trace-based Tests

You can now run tests against your apps on `http://localhost:8080` by going to Tracetest and creating a new HTTP test.

1. Click create a test and select HTTP.
2. Add `http://localhost:8080` in the URL field.
3. Click Run. You’ll see the response and trace data right away.
