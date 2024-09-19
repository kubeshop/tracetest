# Running Node.js with OpenTelemetry and Cloud-based Managed Tracetest Locally

## Install deps

```bash
npm install serverless \
  @opentelemetry/api \
  @opentelemetry/auto-instrumentations-node \
  @opentelemetry/exporter-trace-otlp-grpc \
  @opentelemetry/instrumentation \
  @opentelemetry/sdk-trace-base \
  @opentelemetry/sdk-trace-node
```

## Run Tracetest Agent

1. [Sign in to Tracetest](https://app.tracetest.io/).
2. [Create an Organization](https://docs.tracetest.io/concepts/organizations).
3. Start a Tracetest Cloud Agent by selecting `Application is publicly accessible` in the Settings, and click `Launch Public Agent`.

Tracetest Cloud Agent will run on gRPC and HTTP ports.

- `https://agent-<ID>.tracetest.io:443` — gRPC
- `https://agent-<ID>.tracetest.io:443/v1/traces` — HTTP

## Configure Trace Ingestion

Go to the Trace Ingestion tab in the settings, select OpenTelemetry, and enable the toggle.

## Deploy Serverless App

1. Sign in to your AWS account.
2. Set AWS credentials in the `~.aws` file. [Follow this guide, here](https://www.serverless.com/framework/docs/providers/aws/guide/credentials#recommended-using-local-credentials).

```bash
npm run deploy
```

This will deploy the serverless app and give you an output similar to this:

```bash
> serverless@1.0.0 deploy
> sls deploy


Deploying "otel-serverless-node" to stage "dev" (us-east-1)

✔ Service deployed to stack otel-serverless-node-dev (50s)

endpoint: GET - https://<ID>.execute-api.us-east-1.amazonaws.com/dev/
functions:
  hello: otel-serverless-node-dev-hello (18 MB)
```

## Run Trace-based Tests

You can now run tests against your apps on `https://<ID>.execute-api.us-east-1.amazonaws.com/dev/` by going to Tracetest and creating a new HTTP test.

1. Click create a test and select HTTP.
2. Add `https://<ID>.execute-api.us-east-1.amazonaws.com/dev/` in the URL field.
3. Click Run. You’ll see the response data right away.
