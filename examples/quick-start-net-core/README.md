# Quick Start - .NET Core API with Jaeger, OpenTelemetry and Tracetest

> [Read the detailed recipe for setting up Jaeger with Tractest in our documentation.](https://docs.tracetest.io/examples-tutorials/recipes/running-tracetest-with-jaeger)

This is a simple quick start on how to configure a .NET Core API to use OpenTelemetry instrumentation with traces and Tracetest for enhancing your E2E and integration tests with trace-based testing. The infrastructure will use Jaeger as the trace data store, and OpenTelemetry Collector to receive traces from the API and send them to Jaeger.

## Steps

1. [Install the tracetest CLI](https://github.com/kubeshop/tracetest/blob/main/docs/installing.md#cli-installation)
2. Run `tracetest configure --endpoint http://localhost:11633` on a terminal to configure the CLI to send all commands to that address
3. Run the project by using docker-compose: `docker-compose up -d` (Linux) or `docker compose up -d` (Mac)
4. Test if it works by running: `tracetest test run -d tests/test.yaml`. This would execute a test against the .NET Core API that will send spans to Jaeger to be fetched from the Tracetest server.

Feel free to check out the [docs](https://docs.tracetest.io/), and join our [Discord Community](https://discord.gg/8MtcMrQNbX) for more info!
