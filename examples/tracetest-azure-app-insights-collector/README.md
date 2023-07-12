# Tracetest + OTel Collector + Azure Application Insights (using the OpenTelemetry Collector)

> [Read the detailed recipe for setting up Tracetest + OTel Collector + Azure Application Insights (using the OpenTelemetry Collector) in our documentation.](https://docs.tracetest.io/examples-tutorials/recipes/running-tracetest-with-azure-app-insights-collector)

This repository objective is to show how you can configure your Tracetest instance using the OpenTelemetry collector to send telemetry data to both Azure App Insights and the Tracetest.

## Steps

1. [Install the tracetest CLI](https://docs.tracetest.io/installing/)
2. Run `tracetest configure --endpoint http://localhost:11633` on a terminal
3. Update the `.env` file adding a valid set the valid App Insights Instrumentation Key
4. Run the project by using docker-compose: `docker compose -f ./docker-compose.yaml -f ./tracetest/docker-compose.yaml up -d`
5. Test if it works by running: `tracetest run test -f tests/test.yaml`. This would trigger a test that will send spans to Azure Monitor API and directly to Tracetest that is running on your machine.
