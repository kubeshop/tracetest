# Tracetest Agent with Pokeshop

This example shows how you can use [Tracetest agent](https://docs.tracetest.io/concepts/agent) to capture telemetry data locally and run tests on `app.tracetest.io`.

Note that the environment variable `COLLECTOR_ENDPOINT` is pointing to Tracetest Agent OTLP endpoint, sending all telemetry data directly to it.

To run this example, just set you Environment API key as the env var `TRACETEST_API_KEY` and run:

```sh
TRACETEST_API_KEY=my-api-key docker compose up
```
