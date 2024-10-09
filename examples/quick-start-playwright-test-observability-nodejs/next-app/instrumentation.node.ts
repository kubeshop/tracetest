import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http'
import { Resource } from '@opentelemetry/resources'
import { NodeSDK } from '@opentelemetry/sdk-node'
import { ATTR_SERVICE_NAME } from '@opentelemetry/semantic-conventions'
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node'
import { NodeTracerProvider } from '@opentelemetry/sdk-trace-node'
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base'

const resource = new Resource({
  [ATTR_SERVICE_NAME]: 'next-app',
})
 
const provider = new NodeTracerProvider({ resource: resource })
const exporter = new OTLPTraceExporter()
const processor = new BatchSpanProcessor(exporter)
provider.addSpanProcessor(processor)
provider.register()

const sdk = new NodeSDK({
  traceExporter: exporter,
  instrumentations: [
    getNodeAutoInstrumentations({
      '@opentelemetry/instrumentation-fs': {
        enabled: false
      },
      '@opentelemetry/instrumentation-net': {
        enabled: false
      },
    })
  ],
  serviceName: 'next-app'
})
sdk.start()
