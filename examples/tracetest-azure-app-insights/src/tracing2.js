const { registerInstrumentations } = require("@opentelemetry/instrumentation");
const { NodeTracerProvider } = require("@opentelemetry/sdk-trace-node");
const { BatchSpanProcessor } = require("@opentelemetry/tracing");
const { Resource } = require("@opentelemetry/resources");
const {
  SemanticResourceAttributes,
} = require("@opentelemetry/semantic-conventions");
const {
  getNodeAutoInstrumentations,
} = require("@opentelemetry/auto-instrumentations-node");
const {
  AzureMonitorTraceExporter,
} = require("@azure/monitor-opentelemetry-exporter");
const { diag, DiagConsoleLogger, DiagLogLevel } = require("@opentelemetry/api");
diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.DEBUG);

const azExporter = new AzureMonitorTraceExporter({
  connectionString: process.env.CONNECTION_STRING,
});

// Create and configure NodeTracerProvider
const provider = new NodeTracerProvider({
  resource: new Resource({
    [SemanticResourceAttributes.SERVICE_NAME]: "tracetest",
  }),
});

// Initialize the provider
provider.register();

provider.addSpanProcessor(
  new BatchSpanProcessor(azExporter, {
    bufferSize: 1000, // 1000 spans
    bufferTimeout: 5000, // 5 seconds
  })
);

// register and load instrumentation and old plugins - old plugins will be loaded automatically as previously
// but instrumentations needs to be added
registerInstrumentations({
  instrumentations: getNodeAutoInstrumentations(),
});
