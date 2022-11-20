# Tracetest

Generate end-to-end tests automatically from your traces. For QA, Dev, & Ops.

<!-- 
TODO: migrate video to youtube and use YT embed.

<p align="center">
 <script src="https://fast.wistia.com/embed/medias/dw06408oqz.jsonp" async></script><script src="https://fast.wistia.com/assets/external/E-v1.js" async></script><div class="wistia_responsive_padding" style="padding:56.25% 0 0 0;position:relative;"><div class="wistia_responsive_wrapper" style="height:100%;left:0;position:absolute;top:0;width:100%;"><div class="wistia_embed wistia_async_dw06408oqz videoFoam=true" style="height:100%;position:relative;width:100%"><div class="wistia_swatch" style="height:100%;left:0;opacity:0;overflow:hidden;position:absolute;top:0;transition:opacity 200ms;width:100%;"><img src="https://fast.wistia.com/embed/medias/dw06408oqz/swatch" style="filter:blur(5px);height:100%;object-fit:contain;width:100%;" alt="" aria-hidden="true" onload="this.parentNode.style.opacity=1;" /></div></div></div></div>
</p>

-->

Tracetest allows you to quickly build integration and e2e tests, powered by your OpenTelementry traces.

- Point the system to your Jaeger or Tempo trace datastore.
- Define a triggering transaction, such as a GET against an API endpoint.
- The system runs this transaction, returning both the response data and a full trace.
- Define tests & assertions against this data, ensuring both your response and the underlying processes worked correctly, quickly, and without errors.
- Save your test.
- Run the tests either manually or via your CI build jobs.

## Blog & Video Posts

Check out the following blog posts & videos with Tracetest-related content:

- [Frontend Overhaul of the OpenTelemetry Demo (Go to Next.js)](https://tracetest.io/blog/frontend-overhaul-opentelemetry-demo) - Oct 5, 2022

- [Enabling Tracetest to Work Directly with OpenSearch](https://tracetest.webflow.io/blog/tracetest-opensearch-integration) - Oct 3, 2022

- [Is it Observable? with Henrik Rexed](https://tracetest.webflow.io/blog/is-it-observable-with-henrik-rexed) - Sep 22, 2022

- [How Testability Drives Observability - Open Source Summit Dublin](https://www.youtube.com/watch?v=x5sQg4MNFxI) - Sep 14, 2022

- [Tracetest v0.7 Release Notes](https://tracetest.io/blog/tracetest-v0-7-release-notes) - Aug 23, 2022

- [Tracetest Roadmap Planning - In Person, In Cartagena!](https://tracetest.io/blog/tracetest-roadmap-planning-in-person-in-cartagena) - Aug 23, 2022

- [Common Cypress Testing Pitfalls & How to Avoid Them](https://tracetest.io/blog/common-cypress-testing-pitfalls-how-to-avoid-them) - Aug 18, 2022

- Recorded livestream: [Tracetest v0.6 Release - gRPC, Postman and More](https://www.youtube.com/watch?v=xpEKHK5VXB0) - July 27, 2022

- [Tracetest 0.6 Release - gRPC, Postman and More](https://tracetest.io/blog/tracetest-0-6-release-grpc-postman-and-more) - July 27, 2022

- [Creating a Custom Language Code Editor Using React](https://tracetest.io/blog/creating-a-custom-language-code-editor-using-react) - July 22, 2022

- [Integration Tests: Pros and Cons of Doubles vs. Trace-Based Testing](https://tracetest.io/blog/integration-tests-pros-and-cons-of-doubles-vs-trace-based-testing) - July 13, 2022

- [Detect & Fix Performance Issues Using Tracetest](https://tracetest.io/blog/detect-fix-performance-issues-using-tracetest) - July 6, 2022

- [Integrating Tracetest with GitHub Actions in a CI pipeline](https://tracetest.io/blog/integrating-tracetest-with-github-actions-in-a-ci-pipeline) - June 24, 2022

- Recorded livestream: [Introduction to Tracetest - E2E Tests Powered by OpenTelemetry](https://youtu.be/mqwJRxqBNCg) - June 23, 2022

- [Tracetest — Assertions, Versioning & CI/CD - Release 0.5](https://tracetest.io/blog/tracetest-assertions-versioning-ci-cd) - June 8, 2022

- [Tracing the History of Distributed Tracing & OpenTelemetry](https://kubeshop.io/blog/tracing-the-history-of-distributed-tracing-opentelemetryt) - May 26, 2022

- [Tracetest is released. What’s next?](https://kubeshop.io/blog/tracetest-is-released-whats-next) - May 6, 2022

- [Introducing Tracetest - Trace-based Testing with OpenTelemetry](https://kubeshop.io/blog/introducing-tracetest-trace-based-testing-with-opentelemetry) - April 26, 2022
