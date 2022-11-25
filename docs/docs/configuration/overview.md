# Configuration Overview

There is one way you can set configuration options in Tracetest. By using a configuration file, commonly known as the `tracetest.config.yaml` file.

When using Docker, ensure that the configuration file is mounted to `/app/config.yaml` within the Tracetest Docker container.

To view all the configuration options see the [config file reference page](./config-file-reference).

## Supported data sources

Tracetest is designed to work with different trace data stores. To enable Tracetest to run end-to-end tests against trace data, you need to configure Tracetest to access trace data.

Currently, Tracetest supports the following data stores. Click on the respective data store to view configuration examples:

- [Jaeger](./connecting-to-data-sources/jaeger)
- [Grafana Tempo](./connecting-to-data-sources/tempo)
- [OpenSearch](./connecting-to-data-sources/opensearch)
- [SignalFX](./connecting-to-data-sources/signalfx)

## Using Tracetest without a data source

Another option is to send traces directly to Tracetest using the OpenTelemetry Collector. And, you don't have to change your existing pipelines to do so.

View the [configuration for OpenTelemetry Collector](./connecting-to-data-sources/opentelemetry-collector) for more details.

## Data source configuration examples

Examples of configuring Tracetest to access different data stores can be found in the [`examples` folder of the Tracetest GitHub repo](https://github.com/kubeshop/tracetest/tree/main/examples). 

We will be adding new data stores over the next couple of months - [let us know](https://github.com/kubeshop/tracetest/issues/new/choose) which ones you would like to see us add support for.
