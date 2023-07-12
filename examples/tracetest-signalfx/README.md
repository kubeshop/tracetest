# Tracetest + SignalFX

This repository objective is to show how you can configure your tracetest instance to connect to a signalFX account.

## Steps

1. [Install the tracetest CLI](https://docs.tracetest.io/installing/)
2. Run `tracetest configure --endpoint http://localhost:11633` on a terminal to configure the CLI to send all commands to that address
3. Update the `collector.config.yaml` and `tracetest-config.yaml` with the `token` and `realm` of your SignalFX account.
4. Run the project by using docker-compose: `docker-compose up` (Linux) or `docker compose up` (Mac)
5. Test if it works by running: `tracetest run test -f tests/list-tests.yaml`. This would trigger a test that will send and retrieve spans from the opensearch instance that is running on your machine.
