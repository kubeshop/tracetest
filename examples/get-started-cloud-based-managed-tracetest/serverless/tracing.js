const api = require('@opentelemetry/api');
const { BatchSpanProcessor } = require('@opentelemetry/sdk-trace-base');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-grpc');
const { NodeTracerProvider } = require('@opentelemetry/sdk-trace-node');
const { registerInstrumentations } = require('@opentelemetry/instrumentation');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');

api.diag.setLogger(new api.DiagConsoleLogger(), api.DiagLogLevel.ALL);

const provider = new NodeTracerProvider();
const spanProcessor = new BatchSpanProcessor(
  new OTLPTraceExporter({
    // Use the HTTP endpoint for the agent
    url: 'https://agent-<ID>.tracetest.io:443',
  }),
);
provider.addSpanProcessor(spanProcessor);
provider.register();

registerInstrumentations({
  instrumentations: [
    getNodeAutoInstrumentations({
      '@opentelemetry/instrumentation-aws-lambda': {
        disableAwsContextPropagation: true,
      },
    }),
  ],
});
