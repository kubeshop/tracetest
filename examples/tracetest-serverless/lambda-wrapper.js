const api = require('@opentelemetry/api');
const { BatchSpanProcessor } = require('@opentelemetry/sdk-trace-base');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-grpc');
const { NodeTracerProvider } = require('@opentelemetry/sdk-trace-node');
const { registerInstrumentations } = require('@opentelemetry/instrumentation');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');

api.diag.setLogger(new api.DiagConsoleLogger(), api.DiagLogLevel.ALL);

const { COLLECTOR_ENDPOINT = '' } = process.env;

const provider = new NodeTracerProvider();

const spanProcessor = new BatchSpanProcessor(
  new OTLPTraceExporter({
    url: COLLECTOR_ENDPOINT,
  })
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
