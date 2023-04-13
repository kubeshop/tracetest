# OpenTelemetry Bridge

This document includes design / developer / maintenance notes for the
Node.js APM Agent *OpenTelemetry Bridge*.

Spec: https://github.com/elastic/apm/blob/main/specs/agents/tracing-api-otel.md

## Maintenance

- We should release a new agent version with an updated "@opentelemetry/api"
  dependency relatively soon after any new *minor* release. Otherwise a user
  upgrading their "@opentelemetry/api" dep to "1.x+1", e.g. "1.2.0", will find
  that the OTel Bridge which uses version "1.x", e.g. "1.1.0" or lower, does
  not work.

  The reason is that the OTel Bridge registers global providers (e.g.
  `otel.trace.setGlobalTracerProvider`) with its version of the OTel API. When
  user code attempts to *get* a tracer with **its version** of the OTel API, the
  [OTel API compatibility logic](https://github.com/open-telemetry/opentelemetry-js-api/blob/v1.1.0/src/internal/semver.ts#L24-L33)
  decides that using a v1.1.x Tracer with a v1.2.0 Tracer API is not compatible
  and falls back to a noop implementation.


## Development / Debugging

When doing development on, or debugging the OTel Bridge, it might be helpful to
enable logging of (almost) every `@opentelemetry/api` call into the bridge.
This is done by setting this in `lib/opentelemetry-bridge/setup.js`.

    const LOG_OTEL_API_CALLS = true

It looks like this:

```
% cd test/opentelemetry-bridge/fixtures
% ELASTIC_APM_OPENTELEMETRY_BRIDGE_ENABLED=true node -r ../../../start.js start-span.js
otelapi: OTelTracerProvider.getTracer(...)
otelapi: OTelContextManager.active()
otelapi: OTelTracer.startSpan(name=mySpan, options={}, context=OTelBridgeRunContext<>)
otelapi: OTelContextManager.active()
otelapi: OTelBridgeRunContext.getValue(Symbol(OpenTelemetry Context Key SPAN))
otelapi: OTelSpan<Transaction<52260136515317aa, "mySpan">>.end(endTime=undefined)
```

Together with the agent's usual debug logging, this can help show how the bridge
is working.

```
% ELASTIC_APM_OPENTELEMETRY_BRIDGE_ENABLED=true \
    ELASTIC_APM_LOG_LEVEL=debug \
    node -r ../../../start.js start-span.js | ecslog
...
```


## Naming

In general, the following variable/class/file naming is used:

- A class that implements an OTel interface is prefixed with "OTel". For
  example `class OTelSpan` implements OTel `interface Span`.
- A class that bridges between an OTel interface and an object in the APM
  agent is prefixed with `OTelBridge`. For example `OTelBridgeRunContext`
  bridges between an OTel `interface Context` and the APM agent's `RunContext`,
  i.e. it implements both interfaces/APIs.
- A variable that holds an OpenTelemetry object is prefixed with `otel`, or
  `...OTel...` if it in the middle of the var name. Some examples:
  - `otelSpanOptions` holds an OTel `SpanOptions` instance
  - `parentOTelSpanContext`
  - `epochMsFromOTelTimeInput()` converts from an OTel `TimeInput` to a number
    of milliseconds since the Unix epoch


## Design Overview

The OpenTelemetry API is, currently, [these four interfaces](https://github.com/open-telemetry/opentelemetry-js-api/tree/main/src/api/):

- `otel.context.*` - API for managing Context, i.e. what the APM agent calls
  "run context". More below.
- `otel.trace.*` - API for manipulating spans, and getting a `Tracer` to
  create spans. More below.
- `otel.diag.*` - This is used to hook into internal OpenTelemetry diagnostics,
  i.e. internal logging. There is very little `otel.diag` usage in
  `@opentelemetry/api`, more in the SDK. The APM agent hooks up `otel.diag`
  logging to its own logger **if `logLevel=trace`**.
- `otel.propagation.*` - Used for abstracting trace-context propagation
  (reading/writing "traceparent" et al headers) and Baggage handling. This
  isn't touched by the OTel Bridge, and shouldn't be necessary until either
  the bridge supports Baggage or TextMapPropagator implementations like
  `W3CTraceContextPropagator`. The APM agent implements its own internally.

In `Agent#start()`, if the `opentelemetryBridgeEnabled` config is true, then
a global [`ContextManager`](./OTelContextManager.js) and a global [`TracerProvider`](./OTelTracerProvider.js) are registered, which "enables" the bridge.

From the OTel Bridge spec:

> In order to avoid potentially complex and tedious synchronization issues
> between OTel and our existing agent implementations, the bridge implementation
> SHOULD provide an abstraction to have a single "active context" storage.

For this bridge, the agent's `RunContext` class was extended to support the
small [`interface Context`](https://github.com/open-telemetry/opentelemetry-js-api/blob/v1.1.0/src/context/types.ts#L17-L41)
API and the agent's run context managers were updated to allow passing in a
subclass of `RunContext` to use. So the "single active context storage" is
instances of [`OTelBridgeRunContext`](./OTelBridgeRunContext.js) in the agent's
usual run context managers.

The way the "active span" is tracked by the OTel API is to call
`context.setValue(SPAN_KEY, span)`. The `OTelBridgeRunContext` class translates
calls using `SPAN_KEY` into the API that the agent's RunContext class uses.
Roughly this:

- `context.setValue(SPAN_KEY, span)` -> `this.enterSpan(span)`
- `context.getValue(SPAN_KEY)` -> `return new OTelSpan(this.currSpan())`

Otherwise the `*RunContextManager` classes in the agent map very well to the
OpenTelemetry `ContextManager` interface: the [`OTelContextManager`](./OTelContextManager.js)
implementation is very straightforward.

The `@opentelemetry/api` supports two ways to create objects that are internally
implemented and do not call the registered global providers.

1. `otel.trace.wrapSpanContext(...)` supports creating a `NonRecordingSpan` (a
   class that isn't exported) instance that implements `interface Span`. [This
   test fixture](../../test/opentelemetry-bridge/fixtures/nonrecordingspan-parent.js)
   shows a use case. The bridge wraps this in an `OTelBridgeNonRecordingSpan`
   that implements both OTel `interface Span` and the agent's Transaction API.
2. `otel.ROOT_CONTEXT` is a singleton object (an internal `BaseContext` class
   instance) that implements `interface Context` but is not created via any
   bridge API. That means bridge code cannot rely on a given `context` argument
   being an instance of its `OTelBridgeRunContext` class.
   [This test fixtures](../../test/opentelemetry-bridge/fixtures/using-root-context.js)
   shows an example.

The trickiest part of the bridge is handling these two cases, especially at
the top of `startSpan` in [`OTelTracer`](./OTelTracer.js)


## Limitations / Differences with OpenTelemetry SDK

- The OpenTelemetry SDK defines [SpanLimits](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/sdk.md#span-limits).
  This OpenTelemetry Bridge differs as follows:
  - Attribute count is not limited. The OTel SDK defaults to a limit of 128.
    (To implement this, start at `maybeSetOTelAttr` in "OTelSpan.js".)
  - Attribute value strings are truncated at 1024 bytes. The OpenTelemetry SDK
    uses `AttributeValueLengthLimit (Default=Infinity)`.
    (We could consider using the configurable `longFieldMaxLength` for the
    attribute value truncation limit, if there is a need.)
  - Span events are not currently supported by this bridge.
  - Span link *attributes* are not supported by the bridge (Elastic APM
    supports span links, but not span link attributes).

- The OpenTelemetry Bridge spec says APM agents
  ["MAY"](https://github.com/elastic/apm/blob/main/specs/agents/tracing-api-otel.md#attributes-mapping)
  report OTel span attributes as spad and transaction *labels* if the upstream
  APM Server is less than version 7.16. This implementation opts *not* to do
  that. The OTel spec allows a larger range of types for span attributes values
  than is allowed for "tags" (aka labels) in the APM Server intake API, so some
  further filtering of attributes would be required.

- There is a known issue with the `contextManager: "patch"` config option and
  `tracer.startActiveSpan(name, async function fn () { ... })` where run
  context is lost after the first `await ...` usage in that given `fn`.
  See https://github.com/elastic/apm-agent-nodejs/issues/2679.

- There is a semantic difference between this OTel Bridge and the OpenTelemetry
  SDK with `span.end()` that could impact parent/child relationships of spans.
  This demonstrates the different:

    ```js
    const otel = require('@opentelemetry/api')
    const tracer = otel.trace.getTracer()
    tracer.startActiveSpan('s1', s1 => {
      tracer.startActiveSpan('s2', s2 => {
        s2.end()
      })
      s1.end()
      tracer.startActiveSpan('s3', s3 => {
        s3.end()
      })
    })
    ```

  With the OTel SDK that will yield:

    ```
    span s1
    `- span s2
    `- span s3
    ```

  With the Elastic APM agent:

    ```
    transaction s1
    `- span s2
    transaction s3
    ```

  In current Elastic APM semantics, when a span is ended (e.g. `s1` above) it is
  *no longer the current/active span in that async context*. This is historical
  and allows a stack of current spans in sync code, e.g.:

    ```js
    const t1 = apm.startTransaction('t1')
    const s2 = apm.startSpan('s2')
    const s3 = apm.startSpan('s3') // s3 is a child of s2
    s3.end() // s3 is no longer active (popped off the stack)
    const s4 = apm.startSpan('s4') // s4 is a child of s2
    s4.end()
    s2.end()
    t1.end()
    ```

  This semantic difference is not expected to be common, because it is expected
  that typically OTel API user code will end a span only at the end of its
  function:

    ```js
    tracer.startActiveSpan('mySpan', mySpan => {
      // ...
      mySpan.end() // .end() only at end of function block
    })
    ```

  Note that active span context *is* properly maintained when a new async task
  is created (e.g. with `setTimeout`, etc.), so the following code produces
  the expected trace:

    ```js
    tracer.startActiveSpan('s1', s1 => {
    setImmediate(() => {
      tracer.startActiveSpan('s2', s2 => {
        s2.end()
      })
      setTimeout(() => {  // s1 is bound as the active span in this async task.
        tracer.startActiveSpan('s3', s3 => {
          s3.end()
        })
      }, 100)
      s1.end()
    })
    ```

  If this *does* turn out to be a common issue, the OTel semantics for span.end()
  can likely be accommodated.



