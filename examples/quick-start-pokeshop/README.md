# Quick Start - Tracetest + Pokeshop

> [Read the installation guide in our documentation.](https://docs.tracetest.io/getting-started/installation)

This examples' objective is to show how you can:

1. Configure your Tracetest instance to connect to receive traces from OpenTelemetry Collector.

## Steps

1. [Install the tracetest CLI](https://docs.tracetest.io/installing/)
2. Run `tracetest configure --endpoint http://localhost:11633` on a terminal
3. Run the project by using docker-compose: `docker-compose up -d` (Linux) or `docker compose up -d` (Mac)
4. Test if it works by running: `tracetest run test -f tests/test.yaml`. This would trigger a test that will send and retrieve spans from the OpenTelemetry Collector instance. View the test on `http://localhost:11633`.
