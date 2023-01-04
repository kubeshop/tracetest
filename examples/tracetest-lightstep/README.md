# Tracetest + Otel Demo

This repository objective is to show how you can configure your tracetest instance to connect to the [otel demo](https://github.com/open-telemetry/opentelemetry-demo).

## Steps

1. Run `tracetest configure --endpoint http://localhost:11633` on a terminal to configure the CLI to send all commands to that address
2. Run the project by using docker-compose: `docker-compose up` (Linux) or `docker compose up` (Mac)
3. Validate that the otel demo is running at <http://localhost:8084>
4. Log into the Tracetest front end and create a test using one of the otel demo examples
