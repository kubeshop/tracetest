# Tracetest + Opensearch

This repository objective is to show how you can configure your Tracetest instance to connect to an OpenSearch instance and use it as its tracing backend.

## Steps

1. [Install Tracetest CLI](https://docs.tracetest.io/installing/)
2. Run `tracetest configure --endpoint http://localhost:11633` on a terminal to configure the CLI to send all commands to that address
3. Run the project by using docker-compose: `docker-compose up` (Linux) or `docker compose up` (Mac)
4. Test if it works by running: `tracetest test run -d tests/list-tests.yaml`. This would trigger a test that will send and retrieve spans from the OpenSearch instance that is running on your machine.

> :warning: Note: The OpenSearch configuration used for this example is not meant to be used in production.

## Project Structure

- `opensearch` is a folder that contains all configuration files used to configure the local OpenSearch instance;
- `collector.config.yaml` is the configuration of the Open Telemetry Collector that will receive traces and send them to OpenSearch
- `tracetest-config.yaml` is the internal configuration of the Tracetest instance like database connection, trace polling and observability parameters
- `tracetest-provision.yaml` is a file containing the Data Store setup, showing how Tracetest will connect to the OpenSearch instance
