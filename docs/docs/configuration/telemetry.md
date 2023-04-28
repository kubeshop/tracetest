# Telemetry

The Tracetest server generates internal observability trace data. Use this data to track Tracetest test runs over time.

You can configure an exporter to send the trace data to an OpenTelemetry Collector and then store it safely in your trace data store for further historical analysis.

## Configuring Tracetest Server Internal Telemetry

In the `tracetest-config.yaml` file, alongside the [configuration](./server.md) of connecting Tracetest to the Postgres instance, you can also define a `telemetry` and `server` section.

With these two additional sections, you define an exporter where the Tracetest server's internal telemetry will be routed to. In the `telemetry` section, you define the endpoint of the OpenTelemetry Collector. And, in the `server` section you define which exporter the Tracetest server will use.

```yaml
# tracetest-config.yaml

postgres:
# [...]

telemetry:
  exporters:
    collector:
      serviceName: tracetest
      sampling: 100 # 100%
      exporter:
        type: collector
        collector:
          endpoint: otel-collector:4317
          # Replace with your OpenTelemetry Collector endpoint

server:
  telemetry:
    exporter: collector
```

:::note
Make sure to check what the service endpoint for the OpenTelemetry Collector in your infrastructure is. The example above is using `otel-collector` because that is the service name in Docker Compose. Your infrastructure might differ.
:::
