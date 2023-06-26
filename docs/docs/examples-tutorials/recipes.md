# Recipes

These recipes will show you the best practices for using Tracetest.

## Trace Data Stores

These recipes show integrations with trace data stores and tracing vendors/providers.

### OpenTelemetry Collector

This integration point uses the OpenTelemetry Collector as a router to send trace data to both Tracetest and tracing vendors/providers.

- [Sending traces directly to Tracetest from a Node.js app using OpenTelemetry Collector](./recipes/running-tracetest-without-a-trace-data-store.md)
- [Sending traces with manual instrumentation directly to Tracetest from a Node.js app using OpenTelemetry Collector](./recipes/running-tracetest-without-a-trace-data-store-with-manual-instrumentation.md)
- [Sending traces with manual instrumentation directly to Tracetest from a Python app using OpenTelemetry Collector](./recipes/running-python-app-with-opentelemetry-collector-and-tracetest.md)
- [Sending traces to Lightstep and Tracetest from the OpenTelemetry Demo with OpenTelemetry Collector](./recipes/running-tracetest-with-lightstep.md)
- [Sending traces to New Relic and Tracetest from the OpenTelemetry Demo with OpenTelemetry Collector](./recipes/running-tracetest-with-new-relic.md)
- [Sending traces to Datadog and Tracetest from the OpenTelemetry Demo with OpenTelemetry Collector](./recipes/running-tracetest-with-datadog.md)
- [Sending traces to Honeycomb and Tracetest from a Node.js app using the OpenTelemetry Collector](./recipes/running-tracetest-with-honeycomb.md)

### Jaeger

- [Sending traces to Jaeger from a Node.js app and fetching them from Jaeger with Tracetest](./recipes/running-tracetest-with-jaeger.md)
- [Running Tracetest on AWS Fargate with Terraform](./recipes/running-tracetest-with-aws-terraform.md)

### OpenSearch

- [Sending traces to OpenSearch from a Node.js app and fetching them from OpenSearch with Tracetest](./recipes/running-tracetest-with-opensearch.md)

### Elastic

- [Sending traces to Elastic APM from a Node.js app and fetching them from Elasticsearch with Tracetest](./recipes/running-tracetest-with-elasticapm.md)

### Grafana Tempo

- [Sending traces to Tempo from a Node.js app and fetching them from Tempo with Tracetest](./recipes/running-tracetest-with-tempo.md)

### AWS X-Ray

- [Running Tracetest with AWS X-Ray (AWS X-Ray Node.js SDK)](./recipes/running-tracetest-with-aws-x-ray.md)
- [Running Tracetest with AWS X-Ray (AWS X-Ray Node.js SDK & AWS Distro for OpenTelemetry)](./recipes/running-tracetest-with-aws-x-ray-adot.md)
- [Running Tracetest with AWS X-Ray (AWS Distro for OpenTelemetry & Pokeshop API)](./recipes/running-tracetest-with-aws-x-ray-pokeshop.md)
- [Running Tracetest with AWS Step Functions, AWS X-Ray and Terraform](./recipes/running-tracetest-with-step-functions-terraform.md)

### Azure App Insights

- [Running Tracetest with Azure App Insights (AppInsights Otel Node.js SDK)](./recipes/running-tracetest-with-azure-app-insights.md)
- [Running Tracetest with Azure App Insights (Otel Node.js SDK & OpenTelemetryCollector)](./recipes/running-tracetest-with-azure-app-insights-collector.md)
- [Running Tracetest with Azure App Insights (OpenTelemetry Collector & Pokeshop API)](./recipes/running-tracetest-with-azure-app-insights-pokeshop.md)

## CI/CD Automation

These guides show integrations with CI/CD tools.

- [GitHub Actions](../ci-cd-automation/github-actions-pipeline.md)
- [Testkube](../ci-cd-automation/testkube-pipeline.md)
- [Tekton](../ci-cd-automation/tekton-pipeline.md)

## Tools

These guides show integrations with other tools and vendors.

- [Running Tracetest with Testkube](../tools-and-integrations/testkube.md)
- [Running Tracetest with k6](../tools-and-integrations/k6.md)
- [Running Tracetest with Keptn](../tools-and-integrations/keptn.md)

Stay tuned! More recipes are coming soon. 🚀
