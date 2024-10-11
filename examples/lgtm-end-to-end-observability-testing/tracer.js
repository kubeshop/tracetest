const opentelemetry = require('@opentelemetry/sdk-node')
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node')
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-http')
const { ConsoleSpanExporter } = require('@opentelemetry/sdk-trace-node');
const dotenv = require("dotenv")
dotenv.config()

const sdk = new opentelemetry.NodeSDK({
  // traceExporter: new ConsoleSpanExporter(),
  traceExporter: new OTLPTraceExporter({ url: process.env.OTEL_EXPORTER_OTLP_TRACES_ENDPOINT }),
  instrumentations: [getNodeAutoInstrumentations()],
})
sdk.start()
