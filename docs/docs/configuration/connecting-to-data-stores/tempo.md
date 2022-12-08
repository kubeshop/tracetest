# Tempo

If you want to use Tempo as the trace data store, you can configure Tracetest to fetch trace data from Tempo.

You'll configure the OpenTelemetry Collector to receive traces from your system and then send them to Tempo. And, you don't have to change your existing pipelines to do so.

:::note
It is important to notice that this relies on the [probabilistic_sampler](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/probabilisticsamplerprocessor) processor, which, at the moment, is only available in the [contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib/) version of the collector.
:::

:::tip
Examples of configuring Tracetest can be found in the [`examples` folder of the Tracetest GitHub repo](https://github.com/kubeshop/tracetest/tree/main/examples). 
:::

## Configure OpenTelemetry Collector to send traces to Tempo

In your OpenTelemetry Collector config file, make sure to set the `exporter` to `tempo`, with the `endpoint` pointing to your Tempo's instance on port `4317`. If you are running Tracetest with Docker, the endpoint might look like this `http://tempo:4317`.

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
  otlp/2:
    endpoint: tempo:4317
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
      exporters: [otlp/2] # your exporter pointing to your Tempo instance

```

## Configure Tracetest to use Tempo as a trace data store

First, configure Tempo to run on port `9095`. Here's how your config file should look like:

```yaml
# tempo.config.yaml

auth_enabled: false
server:
  http_listen_port: 3100
  grpc_listen_port: 9095
distributor:
  receivers:                           # this configuration will listen on all ports and protocols that tempo is capable of.
    jaeger:                            # the receives all come from the OpenTelemetry collector.  more configuration information can
      protocols:                       # be found there: https://github.com/open-telemetry/opentelemetry-collector/tree/master/receiver
        thrift_http:                   #
        grpc:                          # for a production deployment you should only enable the receivers you need!
        thrift_binary:
        thrift_compact:
    zipkin:
    otlp:
      protocols:
        http:
        grpc:
    opencensus:
ingester:
  trace_idle_period: 10s               # the length of time after a trace has not received spans to consider it complete and flush it
  max_block_bytes: 1_000_000           # cut the head block when it hits this size or ...
  #traces_per_block: 1_000_000
  max_block_duration: 5m               #   this much time passes
compactor:
  compaction:
    compaction_window: 1h              # blocks in this time window will be compacted together
    max_compaction_objects: 1000000    # maximum size of compacted blocks
    block_retention: 1h
    compacted_block_retention: 10m
storage:
  trace:
    backend: local                     # backend configuration to use
    wal:
      path: /tmp/tempo/wal            # where to store the the wal locally
      #bloom_filter_false_positive: .05 # bloom filter false positive rate.  lower values create larger filters but fewer false positives
      #index_downsample: 10             # number of traces per index record
    local:
      path: /tmp/tempo/blocks
    pool:
      max_workers: 100                 # the worker pool mainly drives querying, but is also used for polling the blocklist
      queue_depth: 10000

```


You also have to configure your Tracetest instance to make it aware that it has to fetch trace data from Tempo. 

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
    tempo:
      type: tempo
      tempo:
        endpoint: tempo:9095
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
    dataStore: tempo
    exporter: collector
    applicationExporter: collector

```

Examples of configuring Tracetest can be found in the [`examples` folder of the Tracetest GitHub repo](https://github.com/kubeshop/tracetest/tree/main/examples). 
