# Jaeger

If you want to use Jaeger as the trace data store, you can configure Tracetest to fetch trace data from Jaeger.

You'll configure the OpenTelemetry Collector to receive traces from your system and then send them to Jaeger. And, you don't have to change your existing pipelines to do so.

:::note
It is important to notice that this relies on the [probabilistic_sampler](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/probabilisticsamplerprocessor) processor, which, at the moment, is only available in the [contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib/) version of the collector.
:::

:::tip
Examples of configuring Tracetest can be found in the [`examples` folder of the Tracetest GitHub repo](https://github.com/kubeshop/tracetest/tree/main/examples). 
:::

## Configure OpenTelemetry Collector to send traces to Jaeger

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
  probabilistic_sampler:
    hash_seed: 22
    sampling_percentage: 100

exporters:
  jaeger:
    endpoint: jaeger:14250
    tls:
      insecure: true

service:
  pipelines:
    # your probably already have a traces pipeline, you don't have to change it.
    # just add this one to your configuration. Just make sure to not have two
    # pipelines with the same name
    traces/1:
      receivers: [otlp] # your receiver
      processors: [tail_sampling, batch] # make sure to have the probabilistic_sampler before your batch processor
      exporters: [jaeger] # your exporter pointing to your Jaeger instance

```

## Configure Tracetest to use Jaeger as a trace data store

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

  exporters:
    collector:
      serviceName: tracetest
      sampling: 100 # 100%
      exporter:
        type: collector
        collector:
          endpoint: otel-collector:4317

server:
  telemetry:
    dataStore: jaeger
    exporter: collector
    applicationExporter: collector

```
