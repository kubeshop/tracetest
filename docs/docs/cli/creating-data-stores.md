# Defining Data Stores as Text Files

You might have multiple Tracetest instances that need to be connected to the same data stores. An easy way of sharing the configuration is by using a configuration file that can be applied to your Tracetest instance.

### Jaeger

```yaml
type: DataStore
spec:
  name: development
  type: jaeger
  default: true
  jaeger:
    endpoint: 127.0.0.1:16685
    tls:
      insecure: true
```

### Tempo

```yaml
type: DataStore
spec:
  name: Grafana Tempo
  type: tempo
  default: true
  tempo:
    endpoint: tempo:9095
    tls:
      insecure: true
```

### OpenSearch

```yaml
type: DataStore
spec:
  name: OpenSearch Data Store
  type: openSearch
  default: true
  opensearch:
    addresses:
      - http://opensearch:9200
    index: traces
```

### SignalFX

```yaml
type: DataStore
spec:
  name: SignalFX
  type: signalFx
  default: true
  signalFx:
    realm: us1
    token: mytoken
```

### Using the OpenTelemetry Collector

```yaml
type: DataStore
spec:
  name: Opentelemetry Collector pipeline
  type: otlp
  default: true
```

> Consider reading about [how to use the OTEL collector](../configuration/connecting-to-data-stores/opentelemetry-collector.md) to send traces to your Tracetest instance.

## Apply Configuration

To apply the configuration, you need a [configured CLI](./configuring-your-cli.md) pointed to the instance you want to apply the data store. Then use the following command:

```
tracetest apply datastore -f my/data-store/file/location.yaml
```

## Additional Information

In the current version, you can only have one active data store at any given time. The flag `default` defines which data store should be used by your tests. So, if you want to add a new data store and make sure it will be used in future test runs, make sure to define `default` as `true` in the data store configuration file.

After a configuration is applied, you can export it using the CLI by using the following command:

```
tracetest export datastore -f my/file/location.yaml --id my-data-store-id
```
