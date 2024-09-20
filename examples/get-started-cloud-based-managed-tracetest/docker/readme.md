# Running Node.js with OpenTelemetry and Cloud-based Managed Tracetest in Docker

## Run Docker Compose (both `app` and `tracetest-agent`)

1. [Sign in to Tracetest](https://app.tracetest.io/).
2. [Create an Organization](https://docs.tracetest.io/concepts/organizations).
3. Retrieve your [Tracetest Organization API Key/Token and Environment ID](https://app.tracetest.io/retrieve-token).
4. Copy the `.env.template` into `.env` and set the [Tracetest Agent](https://docs.tracetest.io/concepts/agent) env vars `<TRACETEST_API_KEY>` and `<TRACETEST_ENVIRONMENT_ID>` from step 2.

```bash
docker compose up -d --build
```

Tracetest Agent will run on gRPC and HTTP ports and use the Docker network where you can access it via its service name.

- `http://tracetest-agent:4317` — gRPC
- `http://tracetest-agent:4318/v1/traces` — HTTP

## Configure Trace Ingestion for Localhost

Go to the Trace Ingestion tab in the settings, select OpenTelemetry, and enable the toggle.

## Run Trace-based Tests

You can now run tests against your apps on `app:8080` by going to Tracetest and creating a new HTTP test.

1. Click create a test and select HTTP.
2. Add `http://app:8080` in the URL field.
3. Click Run. You’ll see the response and trace data right away.
