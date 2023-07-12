const opentelemetry = require("@opentelemetry/sdk-node")
const { getNodeAutoInstrumentations } = require("@opentelemetry/auto-instrumentations-node")
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-grpc')
const { Resource } = require("@opentelemetry/resources")
const { SemanticResourceAttributes } = require("@opentelemetry/semantic-conventions")
const { NodeTracerProvider } = require("@opentelemetry/sdk-trace-node")
const { BatchSpanProcessor } = require("@opentelemetry/sdk-trace-base")

const resource = Resource.default().merge(
  new Resource({
    [SemanticResourceAttributes.SERVICE_NAME]: "quick-start-nodejs-manual-instrumentation",
    [SemanticResourceAttributes.SERVICE_VERSION]: "0.0.1",
  })
)

const provider = new NodeTracerProvider({ resource: resource })
const exporter = new OTLPTraceExporter({ url: 'http://otel-collector:4317' })
const processor = new BatchSpanProcessor(exporter)
provider.addSpanProcessor(processor)
provider.register()

const sdk = new opentelemetry.NodeSDK({
  traceExporter: exporter,
  instrumentations: [getNodeAutoInstrumentations()],
  serviceName: 'quick-start-nodejs-manual-instrumentation'
})
sdk.start()
