// tracing.js
'use strict';

const opentelemetry = require('@opentelemetry/sdk-node');
const { OTLPTraceExporter } =  require('@opentelemetry/exporter-trace-otlp-http');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');

const sdk = new opentelemetry.NodeSDK({
  // environment vars are loaded in the start step
  traceExporter: new OTLPTraceExporter(),
  instrumentations: [
    getNodeAutoInstrumentations({
      // we recommend disabling fs autoinstrumentation since it can be noisy and expensive during startup
      '@opentelemetry/instrumentation-fs': {
          enabled: false,
      },
    }),
  ],
});

sdk.start();
