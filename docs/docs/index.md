# Tracetest Documentation

<!-- 
TODO: migrate video to youtube and use YT embed.

<p align="center">
 <script src="https://fast.wistia.com/embed/medias/dw06408oqz.jsonp" async></script><script src="https://fast.wistia.com/assets/external/E-v1.js" async></script><div class="wistia_responsive_padding" style="padding:56.25% 0 0 0;position:relative;"><div class="wistia_responsive_wrapper" style="height:100%;left:0;position:absolute;top:0;width:100%;"><div class="wistia_embed wistia_async_dw06408oqz videoFoam=true" style="height:100%;position:relative;width:100%"><div class="wistia_swatch" style="height:100%;left:0;opacity:0;overflow:hidden;position:absolute;top:0;transition:opacity 200ms;width:100%;"><img src="https://fast.wistia.com/embed/medias/dw06408oqz/swatch" style="filter:blur(5px);height:100%;object-fit:contain;width:100%;" alt="" aria-hidden="true" onload="this.parentNode.style.opacity=1;" /></div></div></div></div>
</p>

-->

Tracetest allows you to quickly build integration and end-to-end tests, powered by your OpenTelemetry traces.

## In a Nutshell

Tracetest uses your existing [OpenTelemetry](https://opentelemetry.io/docs/getting-started/) traces to power trace-based testing with assertions against your trace data at every point of the request transaction. You only need to point Tracetest to your existing trace data source, or send traces to Tracetest directly!

We make it possible to:

- Define tests and assertions against every single microservice that a request goes through.
- Use your preferred trace back-end, like Jaeger or Tempo, or OpenTelemetry Collector.
- Define multiple transaction triggers, such as a GET against an API endpoint, a GRPC request, etc.
- Return both the response data and a full trace.
- Define assertions against both the response and trace data, ensuring both your response and the underlying processes worked correctly, quickly, and without errors.
- Save tests.
- Run the tests manually or via CI build jobs with the Tracetest CLI.

New to trace-based testing? Read more about the concepts, [here](./concepts/what-is-trace-based-testing).

## Prerequisites

You need to add [OpenTelemetry instrumentation](https://opentelemetry.io/docs/instrumentation/) to your code and configure sending traces to a trace data store, or Tracetest directly, to benefit from Tracetest's trace-based testing.

## Who Uses Tracetest?

Our users are typically developers or QA engineers building distributed systems with microservices using back-end languages like Go, Rust, Node.js, and Python.

Tracetest enables you to write detailed trace-based tests, primarily:

- End-to-end tests
- Integration tests

## What Makes Tracetest Special?

Tracetest can be compared with Cypress or Selenium; however Tracetest is fundamentally different.

Cypress and Selenium are constrained by using the browser for testing. Tracetest bypasses this entirely by using your existing OpenTelemetry instrumentation and trace data to run tests and assertions against traces in every step of a request transaction.

Move on to the [Quick Start](./getting-started/installation) to hit the ground running!
