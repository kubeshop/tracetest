import { CompositePropagator, W3CBaggagePropagator, W3CTraceContextPropagator } from '@opentelemetry/core';
import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { ZoneContextManager } from '@opentelemetry/context-zone';
import { UserInteractionInstrumentation } from '@tracetest/instrumentation-user-interaction';
import { getWebAutoInstrumentations } from '@opentelemetry/auto-instrumentations-web';

const createTracer = async () => {
  const { OTEL_EXPORTER_OTLP_TRACES_ENDPOINT = 'http://localhost:4318/v1/traces' } = process.env

  const provider = new WebTracerProvider({
    resource: new Resource({
      [SemanticResourceAttributes.SERVICE_NAME]: 'react-app',
    }),
  });

  provider.addSpanProcessor(
    new BatchSpanProcessor(
      new OTLPTraceExporter({ url: OTEL_EXPORTER_OTLP_TRACES_ENDPOINT }),
      { maxQueueSize: 10000, scheduledDelayMillis: 200 }
    )
  );

  provider.register({
    contextManager: new ZoneContextManager(),
    propagator: new CompositePropagator({
      propagators: [new W3CBaggagePropagator(), new W3CTraceContextPropagator()],
    }),
  });  

  registerInstrumentations({
    tracerProvider: provider,
    instrumentations: [
      new UserInteractionInstrumentation(),
      getWebAutoInstrumentations({
        '@opentelemetry/instrumentation-fetch': {
          propagateTraceHeaderCorsUrls: /.*/,
          clearTimingResources: true,
        },
        '@opentelemetry/instrumentation-user-interaction': {
          enabled: false,
        },
      }),
    ],
  });
};

createTracer();
