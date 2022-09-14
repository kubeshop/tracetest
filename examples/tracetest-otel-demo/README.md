# Tracetest + Otel Demo

This repository objective is to show how you can configure your tracetest instance to connect to the [otel demo](https://github.com/open-telemetry/opentelemetry-demo).

## Steps

1. Run `tracetest configure` on a terminal and type: <http://localhost:8080> to make your CLI send all requests to that address
2. Run the project by using docker-compose: `docker-compose up` (Linux) or `docker compose up` (Mac)
3. Validate that the otel demo is running at <http://localhost:8084>
4. Log into the Tracetest front end and create a test using one of the otel demo examples
