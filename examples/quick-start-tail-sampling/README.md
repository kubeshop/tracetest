# Quick Start - Trace-based Tests with Tail Sampling Configuration

> [Read the detailed recipe for setting up OpenTelemetry Collector with Tracetest in our documentation.](https://docs.tracetest.io/examples-tutorials/recipes/running-tracetest-without-a-trace-data-store)

This is a simple quick start example on how to set up `tail_sampling` into OTel Collector, allowing Tracetest to run tests in environments where we have a probabilistic sampler enabled and a percentage of the traces are not sent to the final data store.

## Scenario

In this scenario, we have [Go API](./simple-go-service/) that sends Trace data to one instance of [OTel Collector](https://opentelemetry.io/docs/collector/), which samples 50% of the traces and sends it to the [Jaeger](https://www.jaegertracing.io/) data store.

```mermaid
---
title: Initial observability architecture
---
flowchart LR
    GoAPI["Go API"]
    OTelCol["OTel Collector"]
    DataStore[("Jaeger")]

    GoAPI -- Send traces --> OTelCol
    OTelCol -- Forward 50% of Traces --> DataStore
``````

The collector configuration is defined by a `traces/jaeger` pipeline with one receiver, two processors (one to send traces in batch and the other set up the probabilistic sampling) and one exporter to Jaeger:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:
    timeout: 100ms
  probabilistic_sampler:
    hash_seed: 22
    sampling_percentage: 50.0 # 50%

exporters:
  otlp/jaeger:
    endpoint: jaeger:4317 # Send traces to Jaeger
    tls:
      insecure: true

service:
  pipelines:
    traces/jaeger:
      receivers: [otlp]
      processors: [batch, probabilistic_sampler]
      exporters: [otlp/jaeger]

```

However, to run a Trace-based test on this scenario and validate if the API is ok, we need to have the testing trace collected by our OTelCollector. How can we do that?

Since Tracetest always starts a trace by adding a `tracetest=true` key/value to [TraceState](https://opentelemetry.io/docs/specs/otel/trace/api/#tracestate) that is [propagated](https://opentelemetry.io/docs/instrumentation/js/propagation/), we can filter these traces on the collector and send them to Tracetest.

To do this we need to set up Tracetest to use an OTLP Data Store on its provision ([here](./tracetest/tracetest.provision.yaml)) and also set up another pipeline with a [Tail Sampling](https://opentelemetry.io/docs/concepts/sampling/#tail-sampling) processor on the OTel Collector, filtering by TraceState `tracetest=true` value.

We can do that by adding the following configuration to the OTel Collector (the entire collector config YAML is [here](./tracetest/collector.config.yaml)):
```yaml
#...

processors:
  #...
  tail_sampling:
    decision_wait: 10s
    num_traces: 100
    expected_new_traces_per_sec: 10
    policies:
      [
        {
          name: Accept only traces that started from tracetest,
          type: trace_state,
          trace_state: {
            key: tracetest,
            values: ["true"]
          }
        }
      ]

exporters:
  #...
  otlp/tracetest:
    endpoint: tracetest:4317 # Send traces to Tracetest.
    tls:
      insecure: true

service:
  pipelines:
    #...
    traces/tracetest:
      receivers: [otlp]
      processors: [batch, tail_sampling]
      exporters: [otlp/tracetest, logging]

```

With this configuration, we can run Trace-based tests, even with a sampling, into the main data store, having the following structure:

```mermaid
---
title: Observability architecture with Trace-based testing support
---
flowchart LR
    GoAPI["Go API"]
    OTelCol["OTel Collector"]
    DataStore[("Jaeger")]
    Tracetest

    GoAPI -- Send traces --> OTelCol
    OTelCol -- Forward 50% of Traces --> DataStore
    OTelCol -- Forward only `tracetest=true` traces --> Tracetest
``````

## Running the example

If you want to run this example, execute `docker compose up` on this folder.

To execute a Trace-based test with Tracetest against this structure, run `tracetest run test -f test-api-working.yaml`.

Feel free to check out the [docs](https://docs.tracetest.io/), and join our [Slack Community](https://dub.sh/tracetest-community) for more info!
