const opentelemetry = require('@opentelemetry/sdk-node')
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node')
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-grpc')

const config = require("./config")

const sdk = new opentelemetry.NodeSDK({
  traceExporter: new OTLPTraceExporter({ url: config.otelExporterGrpcUrl }),
  instrumentations: [getNodeAutoInstrumentations({
    // load custom configuration for http instrumentation
    '@opentelemetry/instrumentation-fs': {
      enabled: false
    },
  })],
  serviceName: config.otelServiceName
})

sdk.start()
