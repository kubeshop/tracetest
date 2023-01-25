# New Relic

If you want to use [New Relic](https://newrelic.com/) as the trace data store, you'll configure the OpenTelemetry Collector to receive traces from your system and then send them to both Tracetest and New Relic. And, you don't have to change your existing pipelines to do so.

:::tip
Examples of configuring Tracetest with New Relic can be found in the [`examples` folder of the Tracetest GitHub repo](https://github.com/kubeshop/tracetest/tree/main/examples). 
:::

## Configuring OpenTelemetry Collector to Send Traces to both New Relic and Tracetest

In your OpenTelemetry Collector config file, make sure to set the `exporter` to `otlp/tt`, with the `endpoint` pointing to your Tracetest instance on port `21321`. If you are running Tracetest with Docker, the endpoint might look like this `http://tracetest:21321`.

Additionally, set another `exporter` to `otlp/ls`, with the `endpoint` pointing to your New Relic account. Set the endpoint to the New Relic ingest API and add your New Relic access token as a header.

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
  logging:
    logLevel: debug
  # OTLP for Tracetest
  otlp/tt:
    endpoint: tracetest:21321 # Send traces to Tracetest. Read more in docs here:  https://docs.tracetest.io/configuration/connecting-to-data-stores/opentelemetry-collector
    tls:
      insecure: true
  # OTLP for New Relic
  otlp/nr:
    endpoint: otlp.nr-data.net:443
    headers:
      api-key: <new_relic_ingest_licence_key> # Send traces to New Relic.
      # Read more in docs here: https://docs.newrelic.com/docs/more-integrations/open-source-telemetry-integrations/opentelemetry/opentelemetry-setup/#collector-export
      # And here: https://docs.newrelic.com/docs/more-integrations/open-source-telemetry-integrations/opentelemetry/collector/opentelemetry-collector-basic/

service:
  pipelines:
    # Your probably already have a traces pipeline, you don't have to change it.
    # Add these two to your configuration. Just make sure to not have two
    # pipelines with the same name
    traces/tt:
      receivers: [otlp] # your receiver
      processors: [batch]
      exporters: [otlp/tt] # your exporter pointing to your tracetest instance
    traces/nr:
      receivers: [otlp] # your receiver
      processors: [batch]
      exporters: [logging, otlp/nr]  # your exporter pointing to your lighstep account

```

### Configure Tracetest to Use New Relic as a Trace Data Store

You also have to configure your Tracetest instance to expose an `otlp` endpoint to make it aware it will receive traces from the OpenTelemetry Collector.

### Web UI

In the Web UI, open settings, and select New Relic.

![](https://res.cloudinary.com/djwdcmwdz/image/upload/v1674643685/Blogposts/Docs/screely-1674643680615_de8fry.png)


### CLI

Or, if you prefer using the CLI, you can use this file config.

```yaml
type: DataStore
spec:
  name: Opentelemetry Collector pipeline
  type: otlp
  isDefault: true
```

Proceed to run this command in the terminal, and specify the file above.

```bash
tracetest datastore apply -f my/data-store/file/location.yaml
```
