# OpenTelemetry Collector Configuration File Reference

This page contains a reference for using the OpenTelemetry Collector to send trace data from your application to any of Tracetest's supported trace data stores.

:::tip
Examples of configuring Tracetest can be found in the [`examples` folder of the Tracetest GitHub repo](https://github.com/kubeshop/tracetest/tree/main/examples).
:::

## Supported Trace Data Stores

Tracetest is designed to work with different trace data stores. To enable Tracetest to run end-to-end tests against trace data, you need to configure Tracetest to access trace data.

Currently, Tracetest supports the following data stores. Click on the respective data store to view configuration examples:

- [Jaeger](./connecting-to-data-stores/jaeger)
- [OpenSearch](./connecting-to-data-stores/opensearch)
- [Elastic](./connecting-to-data-stores/elasticapm)
- [SignalFX](./connecting-to-data-stores/signalfx)
- [Grafana Tempo](./connecting-to-data-stores/tempo)
- [Lightstep](./connecting-to-data-stores/lightstep)
- [New Relic](./connecting-to-data-stores/new-relic)
- [Datadog](./connecting-to-data-stores/datadog)

Continue reading below to learn how to configure the OpenTelemetry Collector to send trace data from your application to any of the trace data stores above.

## Configure OpenTelemetry Collector to Send Traces to Jaeger

In your OpenTelemetry Collector config file, make sure to set the `exporter` to `jaeger`, with the `endpoint` pointing to your Jaeger instance on port `14250`. If you are running Tracetest with Docker, the endpoint might look like this `http://jaeger:14250`.

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

## Configure OpenTelemetry Collector to Send Traces to OpenSearch

You'll configure the OpenTelemetry Collector to receive traces from your system and then send them to OpenSearch via Data Prepper. And, you don't have to change your existing pipelines to do so.

In your OpenTelemetry Collector config file, make sure to set the `exporter` to `otlp`, with the `endpoint` pointing to the Data Prepper on port `21890`. If you are running Tracetest with Docker, the endpoint might look like this `data-prepper:21890`.

```yaml
# collector.config.yaml

# If you already have receivers declared, you can just ignore
# this one and use yours instead.
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
    loglevel: debug
  otlp/2:
    endpoint: data-prepper:21890
    tls:
      insecure: true
      insecure_skip_verify: true

service:
  pipelines:
    # You probably already have a traces pipeline, you don't have to change it.
    # Just add this one to your configuration. Just make sure to not have two
    # pipelines with the same name.
    traces/1:
      receivers: [otlp] # your receiver
      processors: [batch] # make sure to add the batch processor
      exporters: [otlp/2] # your exporter pointing to your Data Prepper instance

```

## Configure OpenTelemetry Collector to Send Traces to Elastic APM

You'll configure the OpenTelemetry Collector to receive traces from your system and then send them to Elasticsearch via Elastic APM. And, you don't have to change your existing pipelines to do so.

In your OpenTelemetry Collector config file, make sure to set the `exporter` to `otlp`, with the `endpoint` pointing to the Elastic APM server on port `8200`. If you are running Tracetest with Docker, the endpoint might look like this `apm-server:8200`.

```yaml
# collector.config.yaml

# If you already have receivers declared, you can just ignore
# this one and use yours instead.
receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:
    timeout: 100ms

exporters:
  otlp/elastic:
    endpoint: apm-server:8200
    tls:
      insecure: true
      insecure_skip_verify: true

service:
  pipelines:
    # You probably already have a traces pipeline, you don't have to change it.
    # Just add this one to your configuration. Just make sure to not have two
    # pipelines with the same name.
    traces/1:
      receivers: [otlp] # your receiver
      processors: [batch] # make sure to add the batch processor
      exporters: [otlp/elastic] # your exporter pointing to your Elastic APM server instance

```

## Configure OpenTelemetry Collector to Send Traces to SignalFx

You'll configure the OpenTelemetry Collector to receive traces from your system and then send them to SignalFx. And, you don't have to change your existing pipelines to do so.

In your OpenTelemetry Collector config file, make sure to set the `exporter` to `sapm`, with the `endpoint` pointing to the SignalFx trace ingestion endpoint. The endpoint might look like this `https://ingest.us1.signalfx.com/v2/trace`. Also make sure to add your SignalFx `access_token`.

```yaml
# collector.config.yaml

# If you already have receivers declared, you can just ignore
# this one and use yours instead.
receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:
    timeout: 100ms

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
      processors: [batch] # make sure to add the batch processor
      exporters: [sapm] # your exporter pointing to your SignalFx instance
```

## Configure OpenTelemetry Collector to Send Traces to Tempo

You'll configure the OpenTelemetry Collector to receive traces from your system and then send them to Tempo. And, you don't have to change your existing pipelines to do so.

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

