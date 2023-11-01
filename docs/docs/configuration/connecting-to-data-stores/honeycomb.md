# Honeycomb

If you want to use [Honeycomb](https://honeycomb.io/) as the trace data store, you'll configure the OpenTelemetry Collector to receive traces from your system and then send them to both Tracetest and Honeycomb. And, you don't have to change your existing pipelines to do so.

:::tip
Examples of configuring Tracetest with Honeycomb can be found in the [`examples` folder of the Tracetest GitHub repo](https://github.com/kubeshop/tracetest/tree/main/examples).
:::

## Configuring OpenTelemetry Collector to Send Traces to both Honeycomb and Tracetest

In your OpenTelemetry Collector config file:

- Set the `exporter` to `otlp/tracetest`
- Set the `endpoint` to your Tracetest instance on port `4317`

:::tip
If you are running Tracetest with Docker, and Tracetest's service name is `tracetest`, then the endpoint might look like this `http://tracetest:4317`
:::

Additionally, add another config:

- Set the `exporter` to `otlp/honeycomb`
- Set the `endpoint` pointing to the Honeycomb API and using Honeycomb API KEY

```yaml
# collector.config.yaml

# If you already have receivers declared, you can just ignore
# this one and still use yours instead.
receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:
    timeout: 100ms

exporters:
  logging:
    logLevel: debug
  # OTLP for Tracetest
  otlp/tracetest:
    endpoint: tracetest:4317 # Send traces to Tracetest. Read more in docs here:  https://docs.tracetest.io/configuration/connecting-to-data-stores/opentelemetry-collector
    tls:
      insecure: true
  # OTLP for Honeycomb
  otlp/honeycomb:
    endpoint: "api.honeycomb.io:443"
    headers:
      "x-honeycomb-team": "YOUR_API_KEY"
      # Read more in docs here: https://docs.honeycomb.io/getting-data-in/otel-collector/

service:
  pipelines:
    traces/tracetest:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/tracetest]
    traces/honeycomb:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, otlp/honeycomb]
```

## Configure Tracetest to Use Honeycomb as a Trace Data Store

Configure your Tracetest instance to expose an `otlp` endpoint to make it aware it will receive traces from the OpenTelemetry Collector. This will expose Tracetest's trace receiver on port `4317`.

## Connect Tracetest to Honeycomb with the Web UI

In the Web UI, (1) open Settings, and, on the (2) Configure Data Store tab, select (3) Honeycomb.

![Honeycomb](../img/honeycomb-settings.png)

## Connect Tracetest to Honeycomb with the CLI

Or, if you prefer using the CLI, you can use this file config.

```yaml
type: DataStore
spec:
  name: Honeycomb pipeline
  type: honeycomb
  default: true
```

Proceed to run this command in the terminal and specify the file above.

```bash
tracetest apply datastore -f my/data-store/file/location.yaml
```

:::tip
To learn more, [read the recipe on running a sample app with Honeycomb and Tracetest](../../examples-tutorials/recipes/running-tracetest-with-honeycomb.md).
:::