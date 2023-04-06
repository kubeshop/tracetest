# Configuration Overview

There are several configuration options with Tracetest:
- [Server configuration](./server) to set database connection information needed to connect to required PostgreSQL instance.
- [Provisioning configuration](./provisioning) to 'preload' the Tracetest server with resources when first running the Tracetest server.

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
- [AWS X-Ray](./connecting-to-data-stores/awsxray)
- [Datadog](./connecting-to-data-stores/datadog)

## Using Tracetest without a Trace Data Store

Another option is to send traces directly to Tracetest using the OpenTelemetry Collector. And, you don't have to change your existing pipelines to do so.

View the [configuration for OpenTelemetry Collector](./connecting-to-data-stores/opentelemetry-collector) for more details.

## Trace Data Store Configuration Examples

Examples of configuring Tracetest to access different data stores can be found in the [`examples` folder of the Tracetest GitHub repo](https://github.com/kubeshop/tracetest/tree/main/examples). Check out the [**Recipes**](../examples-tutorials/recipes.md) for guided walkthroughs of sample use cases.

We will be adding new data stores over the next couple of months - [let us know](https://github.com/kubeshop/tracetest/issues/new/choose) any additional data stores you would like to see us support.
