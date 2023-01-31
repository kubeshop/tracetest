# Elastic APM

If you want to use Elastic APM as the trace data store, you can configure Tracetest to fetch trace data from Elasticsearch.

You'll configure the OpenTelemetry Collector to receive traces from your system and then send them to Elasticsearch via Elastic APM. And, you don't have to change your existing pipelines to do so.

:::tip
Examples of configuring Tracetest can be found in the [`examples` folder of the Tracetest GitHub repo](https://github.com/kubeshop/tracetest/tree/main/examples). 
:::

## Configure OpenTelemetry Collector to Send Traces to Elastic APM

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

## Configure Tracetest to Use Elastic APM as a Trace Data Store

You also have to configure your Tracetest instance to make it aware that it has to fetch trace data from Elastic APM. 

Make sure you know which Index name, Address, and credentials you are using. In the screenshot below, the Index name is `traces-apm-default`, the Address is `https://es01:9200`, and the Username and Password are to set to `elastic` and `changeme`.

To configure Elastic APM you will need to download the CA certificate from the docker image and upload it to the config under "Upload CA file".

- The command to download the `ca.crt` file is:
`docker cp tracetest-elasticapm-with-elastic-agent-es01-1:/usr/share/elasticsearch/config/certs/ca/ca.crt .`
- Alternatively, you can skip CA certificate validation by setting the `Enable TLS but don't verify the certificate` option.

### Web UI

In the Web UI, open settings, and select Elastic APM.

![](https://res.cloudinary.com/djwdcmwdz/image/upload/v1674566041/Blogposts/Docs/screely-1674566018046_ci0st9.png)


### CLI

Or, if you prefer using the CLI, you can use this file config.

```yaml
type: DataStore
spec:
  name: Elastic Data Store
  type: elasticapm
  isDefault: true
    elasticapm:
      addresses:
        - https://es01:9200
      username: elastic
      password: changeme
      index: traces-apm-default
      insecureSkipVerify: true
```

Proceed to run this command in the terminal, and specify the file above.

```bash
tracetest datastore apply -f my/data-store/file/location.yaml
```
