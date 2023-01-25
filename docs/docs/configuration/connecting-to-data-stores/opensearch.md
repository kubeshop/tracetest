# OpenSearch

If you want to use OpenSearch as the trace data store, you can configure Tracetest to fetch trace data from OpenSearch.

You'll configure the OpenTelemetry Collector to receive traces from your system and then send them to OpenSearch via Data Prepper. And, you don't have to change your existing pipelines to do so.

:::tip
Examples of configuring Tracetest can be found in the [`examples` folder of the Tracetest GitHub repo](https://github.com/kubeshop/tracetest/tree/main/examples). 
:::

## Configure OpenTelemetry Collector to Send Traces to OpenSearch

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

## Configure Tracetest to Use OpenSearch as a Trace Data Store

You also have to configure your Tracetest instance to make it aware that it has to fetch trace data from OpenSearch. 

Make sure you know which Index name and Address you are using. In the screenshot below, the Index name is `traces`, the Address is `http://opensearch:9200`.

### Web UI

In the Web UI, open settings, and select OpenSearch.

![](https://res.cloudinary.com/djwdcmwdz/image/upload/v1674644099/Blogposts/Docs/screely-1674644094600_svcwp6.png)


### CLI

Or, if you prefer using the CLI, you can use this file config.

```yaml
type: DataStore
spec:
  name: OpenSearch Data Store
  type: openSearch
  isDefault: true
  opensearch:
    addresses:
      - http://opensearch:9200
    index: traces
```

Proceed to run this command in the terminal, and specify the file above.

```bash
tracetest datastore apply -f my/data-store/file/location.yaml
```
