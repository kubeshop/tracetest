# Tempo

If you want to use Tempo as the trace data store, you can configure Tracetest to fetch trace data from Tempo.

You'll configure the OpenTelemetry Collector to receive traces from your system and then send them to Tempo. And, you don't have to change your existing pipelines to do so.

:::tip
Examples of configuring Tracetest can be found in the [`examples` folder of the Tracetest GitHub repo](https://github.com/kubeshop/tracetest/tree/main/examples). 
:::

## Configure OpenTelemetry Collector to Send Traces to Tempo

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
      processors: [batch] # make sure to have the probabilistic_sampler before your batch processor
      exporters: [otlp/2] # your exporter pointing to your Tempo instance

```

## Configure Tracetest to Use Tempo as a Trace Data Store

First, configure Tempo to run on port `9095`. Here is an example config file:

```yaml
# tempo.config.yaml

auth_enabled: false
server:
  http_listen_port: 3100
  grpc_listen_port: 9095
distributor:
  receivers:                           # This configuration will listen on all ports and protocols that Tempo is capable of.
    jaeger:                            # the receives all come from the OpenTelemetry collector.  more configuration information can
      protocols:                       # be found here: https://github.com/open-telemetry/opentelemetry-collector/tree/master/receiver.
        thrift_http:                   #
        grpc:                          # For a production deployment you should only enable the receivers you need!
        thrift_binary:
        thrift_compact:
    zipkin:
    otlp:
      protocols:
        http:
        grpc:
    opencensus:
ingester:
  trace_idle_period: 10s               # The length of time after a trace has not received spans to consider it complete and flush it.
  max_block_bytes: 1_000_000           # Cut the head block when it hits this size or ...
  #traces_per_block: 1_000_000
  max_block_duration: 5m               #   this much time passes.
compactor:
  compaction:
    compaction_window: 1h              # Blocks in this time window will be compacted together.
    max_compaction_objects: 1000000    # Maximum size of compacted blocks.
    block_retention: 1h
    compacted_block_retention: 10m
storage:
  trace:
    backend: local                     # Backend configuration to use.
    wal:
      path: /tmp/tempo/wal            # Where to store the the wal locally.
      #bloom_filter_false_positive: .05 # Bloom filter false positive rate.  Lower values create larger filters but fewer false positives.
      #index_downsample: 10             # Number of traces per index record.
    local:
      path: /tmp/tempo/blocks
    pool:
      max_workers: 100                 # The worker pool mainly drives querying, but is also used for polling the blocklist.
      queue_depth: 10000

```

## Configure Tracetest to Use Tempo as a Trace Data Store

You also have to configure your Tracetest instance to make it aware that it has to fetch trace data from Tempo. 

Make sure you know what your Tempo endpoint for fetching traces is. In the screenshot below, the endpoint is `tempo:9095`.

### Web UI

In the Web UI, open settings, and select Tempo.

![](https://res.cloudinary.com/djwdcmwdz/image/upload/v1674644545/Blogposts/Docs/screely-1674644541618_ly8ur3.png)

### CLI

Or, if you prefer using the CLI, you can use this file config.

```yaml
type: DataStore
spec:
  name: Grafana Tempo
  type: tempo
  isDefault: true
  tempo:
    endpoint: tempo:9095
    tls:
      insecure: true
```

Proceed to run this command in the terminal, and specify the file above.

```bash
tracetest datastore apply -f my/data-store/file/location.yaml
```
