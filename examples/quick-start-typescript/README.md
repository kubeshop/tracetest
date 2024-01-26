# Tracetest + Typescript (using @tracetest/client NPM Package)

> [Read the detailed recipe for setting up Tracetest + Typescript (using @tracetest/client NPM Package) in our documentation.](https://docs.tracetest.io/tools-and-integrations/typescript)

This repository's objective is to show how you can run trace-based tests from your Javascript/Typescript environment, including setup stages and waiting for results to be ready.

## Steps

1. Copy the `.env.template` file to `.env`.
2. Log into the [Tracetest app](https://app.tracetest.io/).
3. Fill out the [token](https://docs.tracetest.io/concepts/environment-tokens) and [agent API key](https://docs.tracetest.io/concepts/agent) details.
4. Run `docker compose up -d`.
5. Look for the `tracetest-client` service logs to find out the results from the trace-based tests.
6. Follow the links to the [Tracetest app](https://app.tracetest.io/) to find more details.
