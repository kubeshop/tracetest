const { MeterProvider } = require('@opentelemetry/sdk-metrics');
const { PrometheusExporter } = require('@opentelemetry/exporter-prometheus');
const { Resource } = require('@opentelemetry/resources');
const { SemanticResourceAttributes } = require('@opentelemetry/semantic-conventions');

// Prometheus Exporter for metrics
const prometheusExporter = new PrometheusExporter({
  port: 9464,             // Port where metrics will be exposed
  endpoint: '/metrics',    // Endpoint for Prometheus to scrape
}, () => {
  console.log('Prometheus scrape endpoint: http://localhost:9464/metrics');
});

// MeterProvider for manual metrics instrumentation
const meterProvider = new MeterProvider({
  resource: new Resource({
    [SemanticResourceAttributes.SERVICE_NAME]: 'hello-world-app',  // Use semantic attributes for service name
  }),
});

// Bind the PrometheusExporter as a MetricReader to the MeterProvider
meterProvider.addMetricReader(prometheusExporter);

// Create a meter from the meterProvider
const meter = meterProvider.getMeter('hello-world-meter');

module.exports = meter;



