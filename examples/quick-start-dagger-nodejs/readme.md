# Quick Start - Dagger CI for a Node.js app with OpenTelemetry and Tracetest

> [Read the detailed recipe for setting up OpenTelemetry Collector with Tractest in our documentation.](https://docs.tracetest.io/examples-tutorials/recipes/running-tracetest-without-a-trace-data-store)

This is a simple quick start on how to configure a Node.js app to use OpenTelemetry instrumentation with traces, and Tracetest for enhancing your e2e and integration tests with trace-based testing.

## Steps to run Tracetest with Dagger

### Install Dagger CLI

[Follow the guide in the Dagger docs.](https://docs.dagger.io/quickstart/cli)

### Install Dagger Tracetest Module

```bash
dagger install github.com/kubeshop/tracetest@8c44d01e33e677518a555d1b958cf7be0d70f940
```

### Run Dagger

```bash
export TRACETEST_API_KEY=<your-api-key>
export TRACETEST_ENVIRONMENT_ID=<your-env-id>
export TRACETEST_ORGANIZATION_ID=<your-org-id>

dagger call --api-key=$TRACETEST_API_KEY --environment=$TRACETEST_ENVIRONMENT_ID --organization=$TRACETEST_ORGANIZATION_ID tracetest --source=.
```
