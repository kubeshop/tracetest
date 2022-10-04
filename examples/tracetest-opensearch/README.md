# Tracetest + Opensearch

This repository objective is to show how you can configure your tracetest instance to connect to an opensearch instance and use it as its tracing backend.

## Steps

1. [Install the tracetest CLI](https://github.com/kubeshop/tracetest/blob/main/docs/installing.md#cli-installation)
2. Run `tracetest configure --endpoint http://localhost:11633` on a terminal to configure the CLI to send all commands to that address
3. Run the project by using docker-compose: `docker-compose up` (Linux) or `docker compose up` (Mac)
4. Test if it works by running: `tracetest test run -d tests/list-tests.yaml`. This would trigger a test that will send and retrieve spans from the opensearch instance that is running on your machine.

> :warning: Note: The OpenSearch configuration used for this example is not meant to be used in production.

## Project structure

- `opensearch` is a folder that contains all configuration files used to configure the local opensearch instance;
- `collector.config.yaml` is the configuration of the opentelemetry collector that will receive traces and send them to opensearch
- `tracetest-config.yaml` is the configuration of the tracetest instance, including how to connect to the opensearch instance
