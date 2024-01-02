# OpenTelemetry Demo with Tracetest and Dynatrace

> [Read the detailed recipe for setting up Dynatrace with Tractest in our documentation.](https://docs.tracetest.io/examples-tutorials/recipes/running-tracetest-with-dynatrace)

This example uses the OpenTelemetry Demo `v1.3.0`.

This is a simple sample app on how to configure the [OpenTelemetry Demo](https://github.com/open-telemetry/opentelemetry-demo) to use [Tracetest](https://tracetest.io/) for enhancing your E2E and integration tests with trace-based testing, and [Dynatrace](https://www.dynatrace.com/) as a trace data store.

Feel free to check out the [docs](https://docs.tracetest.io/), and join our [Slack Community](https://dub.sh/tracetest-community) for more info!

You can run it locally using the command:

```sh
docker compose -f ./docker-compose.yaml -f ./tracetest/docker-compose.yaml up
```
