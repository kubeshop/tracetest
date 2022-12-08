# SignalFx

If you want to use SignalFx as the trace data store, you can configure Tracetest to fetch trace data from SignalFx.

You'll configure the OpenTelemetry Collector to receive traces from your system and then send them to SignalFx. And, you don't have to change your existing pipelines to do so.

:::note
It is important to notice that this relies on the [probabilistic_sampler](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/probabilisticsamplerprocessor) processor, which, at the moment, is only available in the [contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib/) version of the collector.
:::

:::tip
Examples of configuring Tracetest can be found in the [`examples` folder of the Tracetest GitHub repo](https://github.com/kubeshop/tracetest/tree/main/examples). 
:::

## Configure OpenTelemetry Collector to send traces to SignalFx

In your OpenTelemetry Collector config file, make sure to set the `exporter` to `sapm`, with the `endpoint` pointing to the SignalFx trace ingestion endpoint. The endpoint might look like this `https://ingest.us1.signalfx.com/v2/trace`. Also make sure to add your SignalFx `access_token`.

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
  sapm:
    access_token: <YOUR_TOKEN> # UPDATE THIS
    access_token_passthrough: true
    endpoint: https://ingest.us1.signalfx.com/v2/trace # UPDATE THIS IF NEEDED
    max_connections: 100
    num_workers: 8

service:
  pipelines:
    # your probably already have a traces pipeline, you don't have to change it.
    # just add this one to your configuration. Just make sure to not have two
    # pipelines with the same name
    traces/1:
      receivers: [otlp] # your receiver
      processors: [probabilistic_sampler, batch] # make sure to have the probabilistic_sampler before your batch processor
      exporters: [sapm] # your exporter pointing to your SignalFx instance

```

## Configure Tracetest to use SignalFx as a trace data store

You also have to configure your Tracetest instance to make it aware that it has to fetch trace data from SignalFx.

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
    signalfx:
      type: signalfx
      signalfx:
        token: <YOUR_TOKEN> # UPDATE WITH YOUR TOKEN
        realm: <YOUR_REALM> # UPDATE WITH YOUR REALM

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
    dataStore: signalfx
    exporter: collector
    applicationExporter: collector

```

Examples of configuring Tracetest can be found in the [`examples` folder of the Tracetest GitHub repo](https://github.com/kubeshop/tracetest/tree/main/examples). 
