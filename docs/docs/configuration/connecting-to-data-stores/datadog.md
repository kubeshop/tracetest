# Datadog

If you want to use [Datadog](https://www.datadoghq.com/) as the trace data store, you'll configure the OpenTelemetry Collector to receive traces from your system and then send them to both Tracetest and Datadog. And, you don't have to change your existing pipelines to do so.

:::tip
Examples of configuring Tracetest with Datadog can be found in the [`examples` folder of the Tracetest GitHub repo](https://github.com/kubeshop/tracetest/tree/main/examples). 
:::

## Configuring OpenTelemetry Collector to Send Traces to both Datadog and Tracetest

In your OpenTelemetry Collector config file, make sure to set the `exporter` to `otlp/tt`, with the `endpoint` pointing to your Tracetest instance on port `21321`. If you are running Tracetest with Docker, the endpoint might look like this `http://tracetest:21321`.

Additionally, set another `exporter` to `datadog`, with the `endpoint` pointing to your Datadog account. Set the site to the Datadog API `datadoghq.com` and add your API key.

```yaml
# collector.config.yaml

# If you already have receivers declared, you can just ignore
# this one and still use yours instead.
receivers:
  otlp:
    protocols:
      http:
      grpc:

processors:
  batch:
    send_batch_max_size: 100
    send_batch_size: 10
    timeout: 10s

exporters:
  # OTLP for Tracetest
  otlp/tt:
    endpoint: tracetest:21321 # Send traces to Tracetest. Read more in docs here:  https://docs.tracetest.io/configuration/connecting-to-data-stores/opentelemetry-collector
    tls:
      insecure: true
  # Datadog exporter
  datadog:
    api:
      site: datadoghq.com
      key: <datadog_API_key> # Add here you API key for Datadog
      # Read more in docs here: https://docs.datadoghq.com/opentelemetry/otel_collector_datadog_exporter
service:
  pipelines:
    traces/tt:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/tt] # exporter sending traces to your Tracetest instance
    traces/dd:
      receivers: [otlp]
      processors: [batch]
      exporters: [datadog] # exporter sending traces to directly to Datadog
```

### Configure Tracetest to Use Lightstep as a Trace Data Store

You also have to configure your Tracetest instance to expose an `otlp` endpoint to make it aware it will receive traces from the OpenTelemetry Collector.

### Web UI

In the Web UI, open settings, and select Datadog.

![](../img/configure-datadog.png)

### CLI

Or, if you prefer using the CLI, you can use this file config.

```yaml
type: DataStore
spec:
  name: OpenTelemetry Collector pipeline
  type: otlp
  isDefault: true
```

Proceed to run this command in the terminal, and specify the file above.

```bash
tracetest datastore apply -f my/data-store/file/location.yaml
```
