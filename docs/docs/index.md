# Tracetest Documentation

Tracetest allows you to quickly build integration and end-to-end tests, powered by your OpenTelemetry traces.

## Prerequisites

Tracetest requires that you have [OpenTelemetry instrumentation](https://opentelemetry.io/docs/instrumentation/) added in your code.

:::tip Don't have OpenTelemetry installed?
[Follow these instructions to install OpenTelemetry in 5 minutes without any code changes!](./getting-started/no-otel.md)
:::

## In a Nutshell

Tracetest uses existing [OpenTelemetry](https://opentelemetry.io/docs/getting-started/) traces to power trace-based testing with assertions against your trace data at every point of a request transaction.

You only need to point Tracetest to your existing [trace data source](./configuration/connecting-to-data-stores/jaeger.md) or [send traces to Tracetest](./configuration/connecting-to-data-stores/opentelemetry-collector.md) directly!

With Tracetest you can:

- Define tests and assertions against every single microservice that a trace goes through.
- Work with your existing distributed tracing solution, allowing you to build tests based on your already instrumented system.
- Define multiple transaction triggers, such as a GET against an API endpoint, a GRPC request, etc.
- Define assertions against both the response and trace data, ensuring both your response and the underlying processes worked correctly, quickly, and without errors.
- Save and run the tests manually or via CI build jobs with

New to trace-based testing? Read more about the concepts [here](./concepts/what-is-trace-based-testing).

## Architecture

Here you can see how Tracetest interacts with a system under test.

1. Trigger a test and generate a trace response
2. Fetch traces to render and analyze them
3. Add assertions to traces
4. See test results
5. Run tests as part of CI/CD pipelines

![Marketechture](https://res.cloudinary.com/djwdcmwdz/image/upload/v1686654113/docs/tracetest-marketechture-jun12-v3_ffj2e2.png)

## Who Uses Tracetest?

Our users are typically developers or QA engineers building distributed systems with microservices using back-end languages like Go, Rust, Node.js and Python.

Tracetest enables you to write detailed trace-based tests, primarily:

- End-to-end tests
- Integration tests

## What Makes Tracetest Special?

Tracetest can be compared with Cypress or Selenium; however Tracetest is fundamentally different.

Cypress and Selenium are constrained by using the browser for testing. Tracetest bypasses this entirely by using your existing OpenTelemetry instrumentation and trace data to run tests and assertions against traces in every step of a request transaction.
