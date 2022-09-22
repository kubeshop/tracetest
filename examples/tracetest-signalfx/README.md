# Tracetest + SignalFX

This repository objective is to show how you can configure your tracetest instance to connect to a signalFX account.

## Steps

1. [Install the tracetest CLI](https://github.com/kubeshop/tracetest/blob/main/docs/installing.md#cli-installation)
2. Run `tracetest configure` on a terminal and type: `http://localhost:11633` to make your CLI send all requests to that address
3. Update the `collector.config.yaml` and `tracetest-config.yaml` with the `token` and `realm` of your SignalFX account.
4. Run the project by using docker-compose: `docker-compose up` (Linux) or `docker compose up` (Mac)
5. Test if it works by running: `tracetest test run -d tests/list-tests.yaml`. This would trigger a test that will send and retrieve spans from the opensearch instance that is running on your machine.
