# Step-by-step

1. Install dependencies

      ```bash
      npm i @opentelemetry/core \
        @opentelemetry/sdk-trace-web \
        @opentelemetry/sdk-trace-base \
        @opentelemetry/instrumentation \
        @opentelemetry/exporter-trace-otlp-http \
        @opentelemetry/context-zone \
        @tracetest/instrumentation-user-interaction \
        @opentelemetry/auto-instrumentations-web \
        @opentelemetry/resources
      ```

2. Initialize Tracing

`instrumentation.js`

```js
import { CompositePropagator, W3CBaggagePropagator, W3CTraceContextPropagator } from "@opentelemetry/core";
import { WebTracerProvider } from "@opentelemetry/sdk-trace-web";
import { BatchSpanProcessor } from "@opentelemetry/sdk-trace-base";
import { registerInstrumentations } from "@opentelemetry/instrumentation";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-http";
import { ZoneContextManager } from "@opentelemetry/context-zone";
import { UserInteractionInstrumentation } from "@tracetest/instrumentation-user-interaction";
import { getWebAutoInstrumentations } from "@opentelemetry/auto-instrumentations-web";
import { Resource } from "@opentelemetry/resources";

const provider = new WebTracerProvider({
  resource: new Resource({
    "service.name": "browser-app",
  }),
});

provider.addSpanProcessor(
  new BatchSpanProcessor(new OTLPTraceExporter({ url: "http://localhost:4318/v1/traces" }), {
    maxQueueSize: 10000,
    scheduledDelayMillis: 200,
  })
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
      "@opentelemetry/instrumentation-fetch": {
        propagateTraceHeaderCorsUrls: /.*/,
        clearTimingResources: true,
      },
      "@opentelemetry/instrumentation-user-interaction": {
        enabled: false,
      },
    }),
  ],
});
```

`index.js` or entrypoint file

```js
import "./instrumentation";
// rest of the app's entrypoint code
```
