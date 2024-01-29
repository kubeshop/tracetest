const opentelemetry = require('@opentelemetry/sdk-node')
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node')
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-grpc')
const { OTLPMetricExporter } = require('@opentelemetry/exporter-metrics-otlp-grpc')
const { PeriodicExportingMetricReader } = require('@opentelemetry/sdk-metrics')
const grpc = require('@grpc/grpc-js')

const exporterConfig = {
  url: 'localhost:4317',
  credentials: grpc.ChannelCredentials.createInsecure()
}

const sdk = new opentelemetry.NodeSDK({
  metricReader: new PeriodicExportingMetricReader({
    exporter: new OTLPMetricExporter(exporterConfig)
  }),
  traceExporter: new OTLPTraceExporter(exporterConfig),
  instrumentations: [getNodeAutoInstrumentations()],
  serviceName: 'test-api',
})
sdk.start()
