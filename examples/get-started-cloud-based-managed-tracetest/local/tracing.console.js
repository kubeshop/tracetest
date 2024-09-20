// Sample for exporting traces to the console.

const { NodeSDK } = require('@opentelemetry/sdk-node');
const { ConsoleSpanExporter } = require('@opentelemetry/sdk-trace-node');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');

// Configure Console Trace Exporter
const traceExporter = new ConsoleSpanExporter();

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
