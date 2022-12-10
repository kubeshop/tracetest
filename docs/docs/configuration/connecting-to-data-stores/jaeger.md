# Jaeger

If you want to use Jaeger as the trace data store, you can configure Tracetest to fetch trace data from Jaeger.

You'll configure the OpenTelemetry Collector to receive traces from your system and then send them to Jaeger. And, you don't have to change your existing pipelines to do so.

:::tip
Examples of configuring Tracetest can be found in the [`examples` folder of the Tracetest GitHub repo](https://github.com/kubeshop/tracetest/tree/main/examples). 
:::

## Configure OpenTelemetry Collector to Send Traces to Jaeger

In your OpenTelemetry Collector config file, make sure to set the `exporter` to `jaeger`, with the `endpoint` pointing to your Jaeger's instance on port `14250`. If you are running Tracetest with Docker, the endpoint might look like this `http://jaeger:14250`.

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
  jaeger:
    endpoint: jaeger:14250
    tls:
      insecure: true

service:
  pipelines:
    # You probably already have a traces pipeline, you don't have to change it.
    # Just add this one to your configuration. Just make sure to not have two
    # pipelines with the same name.
    traces/1:
      receivers: [otlp] # your receiver
      processors: [batch] # make sure to add the batch processor
      exporters: [jaeger] # your exporter pointing to your Jaeger instance

```

## Configure Tracetest to Use Jaeger as a Trace Data Store

You also have to configure your Tracetest instance to make it aware that it has to fetch trace data from Jaeger. 

Edit your configuration file to include this configuration:

```yaml
# tracetest.config.yaml

postgresConnString: "host=postgres user=postgres password=postgres port=5432 sslmode=disable"

poolingConfig:
  maxWaitTimeForTrace: 10m
  retryDelay: 5s

googleAnalytics:
  enabled: true

demo:
  enabled: []

experimentalFeatures: []

telemetry:
  dataStores:
    jaeger:
      type: jaeger
      jaeger:
        endpoint: jaeger:16685
        tls:
          insecure: true

server:
  telemetry:
    dataStore: jaeger
    exporter: collector
    applicationExporter: collector

```
