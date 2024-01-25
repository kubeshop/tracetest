# Tracetest env var provisoning

This repository objective is to show how you can provision a Tracetest instance using a environment variable

## Steps

1. [Install the Tracetest CLI](https://docs.tracetest.io/installing/)
2. Run `tracetest configure --server-url http://localhost:11633` on a terminal
3. Run the project by using docker-compose: `docker-compose up` (Linux) or `docker compose up` (Mac)
4. Test if it works by running: `tracetest run test -f tests/list-tests.yaml`. This would trigger a test that will send and retrieve spans from the Jaeger instance that is running on your machine.
