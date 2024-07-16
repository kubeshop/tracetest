# Tracetest Playwright Engine Trigger

This repository's objective is to show how you can configure trace-based tests using the Tracetest Playwright Engine Trigger.

## Documentation Recipe

This example is part of the official Tracetest docs and can be found by following this [link](https://docs.tracetest.io/examples-tutorials/recipes/running-tests-with-tracetest-playwright-engine).

## Steps

1. Copy the `.env.template` file to `.env`.
2. Log into the [Tracetest app](https://app.tracetest.io/).
3. Fill out the [TRACETEST_API_TOKEN](https://docs.tracetest.io/concepts/environment-tokens) with an admin role token and the [TRACETEST_ENVIRONMENT_ID](https://docs.tracetest.io/concepts/environments) with the id of your environment.
4. Run `docker compose run tracetest-run`.
5. Follow the links in the log to view the test runs programmatically created by your Playwright script.
