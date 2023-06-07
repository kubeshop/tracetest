# Defining Data Stores as Text Files

You might have multiple Tracetest instances that need to be connected to the same data stores. An easy way of sharing the configuration is by using a configuration file that can be applied to your Tracetest instance.

## Supported Trace Data Stores

### Jaeger

```yaml
type: DataStore
spec:
  name: Jaeger
  type: jaeger
  default: true
  jaeger:
    endpoint: jaeger:16685
    tls:
      insecure: true
```

### OpenSearch

```yaml
type: DataStore
spec:
  name: OpenSearch
  type: opensearch
  default: true
  opensearch:
    addresses:
      - http://opensearch:9200
    index: traces
```

### Elastic APM

```yaml
type: DataStore
spec:
  name: Elastic APM
  type: elasticapm
  default: true
  elasticapm:
    addresses:
      - https://es01:9200
    username: elastic
    password: changeme
    index: traces-apm-default
    insecureSkipVerify: true
```

### SignalFX

```yaml
type: DataStore
spec:
  name: SignalFX
  type: signalfx
  default: true
  signalfx:
    realm: us1
    token: mytoken
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

### Lightstep

```yaml
type: DataStore
spec:
  name: Lightstep pipeline
  type: lightstep
  default: true
```

### New Relic

```yaml
type: DataStore
spec:
  name: New Relic pipeline
  type: newrelic
  default: true
```

### AWS X-Ray

```yaml
type: DataStore
spec:
  name: AWS X-Ray
  type: awsxray
  default: true
  awsxray:
    accessKeyId: <your-accessKeyId>
    secretAccessKey: <your-secretAccessKey>
    sessionToken: <your-session-token>
    region: "us-west-2"
```

### Datadog

```yaml
type: DataStore
spec:
  name: Datadog pipeline
  type: datadog
  default: true
```

### Honeycomb

```yaml
type: DataStore
spec:
  name: Honeycomb pipeline
  type: honeycomb
  default: true
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
