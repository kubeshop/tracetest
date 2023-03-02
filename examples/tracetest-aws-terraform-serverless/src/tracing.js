const { Resource } = require("@opentelemetry/resources");
const api = require("@opentelemetry/api");
const { BatchSpanProcessor } = require("@opentelemetry/sdk-trace-base");
const { OTLPTraceExporter } = require("@opentelemetry/exporter-trace-otlp-http");
const { NodeTracerProvider } = require("@opentelemetry/sdk-trace-node");
const { registerInstrumentations } = require("@opentelemetry/instrumentation");
const { getNodeAutoInstrumentations } = require("@opentelemetry/auto-instrumentations-node");
const { SemanticResourceAttributes } = require("@opentelemetry/semantic-conventions");

const provider = new NodeTracerProvider({
  resource: new Resource({
    [SemanticResourceAttributes.SERVICE_NAME]: "tracetest",
  }),
});
const spanProcessor = new BatchSpanProcessor(new OTLPTraceExporter());

provider.addSpanProcessor(spanProcessor);
provider.register();

registerInstrumentations({
  instrumentations: [
    getNodeAutoInstrumentations({
      "@opentelemetry/instrumentation-aws-lambda": {
        disableAwsContextPropagation: true,
        eventContextExtractor: (event) => {
          // Tracetest sends the traceparent header as Traceparent
          // passive aggressive 🥴 The propagators doesn't care about changing the casing before matching
          event.headers.traceparent = event.headers.Traceparent;
          const eventContext = api.propagation.extract(api.context.active(), event.headers);

          return eventContext;
        },
      },
    }),
  ],
});
