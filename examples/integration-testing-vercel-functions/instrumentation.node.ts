import { NodeSDK } from '@opentelemetry/sdk-node'
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http'
import { Resource } from '@opentelemetry/resources'
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions'
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node'
import { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch'

const sdk = new NodeSDK({
  // The OTEL_EXPORTER_OTLP_ENDPOINT env var is passed into "new OTLPTraceExporter" automatically.
  // If the OTEL_EXPORTER_OTLP_ENDPOINT env var is not set the "new OTLPTraceExporter" will
  // default to use "http://localhost:4317" for gRPC and "http://localhost:4318" for HTTP.
  // This sample is using HTTP.
  traceExporter: new OTLPTraceExporter(),
  instrumentations: [
    getNodeAutoInstrumentations(),
    new FetchInstrumentation(),
  ],
  resource: new Resource({
    [SemanticResourceAttributes.SERVICE_NAME]: 'next-app',
  }),
})
sdk.start()
