const { NodeSDK } = require('@opentelemetry/sdk-node');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-grpc');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');

// Configure OTLP gRPC Trace Exporter
const traceExporter = new OTLPTraceExporter({
  // Default endpoint for OTLP gRPC is localhost:4317
  // You can change this to your OTLP collector or backend URL.
  url: 'http://tracetest-agent.default.svc.cluster.local:4317',
});

// Initialize the OpenTelemetry Node SDK
const sdk = new NodeSDK({
  traceExporter,
  instrumentations: [getNodeAutoInstrumentations()],
});

// Start the SDK (this enables tracing)
sdk.start();

// Graceful shutdown on exit
process.on('SIGTERM', () => {
  sdk.shutdown().then(() => {
    console.log('OpenTelemetry tracing terminated');
    process.exit(0);
  });
});
