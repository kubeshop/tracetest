const opentelemetry = require("@opentelemetry/sdk-node");
const { getNodeAutoInstrumentations } = require("@opentelemetry/auto-instrumentations-node");
const { OTLPTraceExporter } = require("@opentelemetry/exporter-trace-otlp-http");
const { Resource } = require("@opentelemetry/resources");
const { NodeTracerProvider, ConsoleSpanExporter } = require("@opentelemetry/sdk-trace-node");
const { SemanticResourceAttributes } = require("@opentelemetry/semantic-conventions");

const provider = new NodeTracerProvider({
  resource: new Resource({
    [SemanticResourceAttributes.SERVICE_NAME]: "node-app",
  }),
});

const sdk = new opentelemetry.NodeSDK({
  traceExporter: new OTLPTraceExporter(),
  // traceExporter: new ConsoleSpanExporter(),
  tracerProvider: provider,
  instrumentations: [
    getNodeAutoInstrumentations({
      "@opentelemetry/instrumentation-fs": {
        enabled: false,
      },
      "@opentelemetry/instrumentation-net": {
        enabled: false,
      },
    }),
  ],
  serviceName: "node-app",
});
sdk.start();
