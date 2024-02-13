# Tracetest + Playwright

This repository's objective is to show how you can configure trace-based tests using Tracetest and Playwright.

## Steps

1. Copy the `.env.template` file to `.env`.
2. Log into the [Tracetest app](https://app.tracetest.io/).
3. This example is configured to use the OpenTelemetry Collector. Ensure the environment you will be utilizing to run this example is also configured to use the OpenTelemetry Tracing Backend by clicking on Settings, Tracing Backend, OpenTelemetry, and Save.
4. Fill out the [token](https://docs.tracetest.io/concepts/environment-tokens) and [agent API key](https://docs.tracetest.io/concepts/agent) details by editing your .env file. You can find these values in the Settings area for your environment.
5. Run `docker compose up -d`.
6. Look for the `tracetest-e2e` service in Docker and click on it to view the logs. It will show the results from the trace-based tests that are triggered by the Playwright tests.
7. Follow the links in the log to view the test runs programmatically created by your Playwright script.
