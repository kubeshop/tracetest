# Tracetest + Amazon X-Ray

This repository objective is to show how you can configure your Tracetest instance to connect to AWS X-Ray plus the ADOT collector and use it as its tracing backend.

## Steps

1. [Install the tracetest CLI](https://docs.tracetest.io/installing/)
2. Run `tracetest configure --endpoint http://localhost:11633` on a terminal
3. Update the `.env` file adding a valid set of AWS credentials
4. Run the project by using docker-compose: `docker-compose up -d` (Linux) or `docker compose up -d` (Mac)
5. Test if it works by running: `tracetest test run -d tests/test.yaml`. This would trigger a test that will send and retrieve spans from the X-Ray instance that is running on your machine.
