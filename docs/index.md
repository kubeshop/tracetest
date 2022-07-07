# Tracetest

Tracetest - Trace-based testing. End-to-end tests powered by your OpenTelemetry Traces.

<p align="center">
  <a target="_new" href="https://www.youtube.com/watch?v=WMRicNlaehc">
    <img src="https://img.youtube.com/vi/WMRicNlaehc/0.jpg" style="width:66%;height:auto">
    <p align="center">
      Click on the image or this link to watch the "Tracetest Intro Video" video (7 mins)
    </p>
  </a>
</p>

Tracetest allows you to quickly build integration and e2e tests, powered by your OpenTelementry traces.

- Point the system to your Jaeger or Tempo trace datastore.
- Define a triggering transaction, such as a GET against an API endpoint.
- The system runs this transaction, returning both the response data and a full trace.
- Define tests & assertions against this data, ensuring both your response and the underlying processes worked correctly, quickly, and without errors.
- Save your test.
- Run the tests either manually or via your CI build jobs.

## **Blog & Video Posts**

Check out the following blog-posts with Tracetest-related content:

- [Detect & Fix Performance Issues Using Tracetest](https://kubeshop.io/blog/detect-fix-performance-issues-using-tracetest) - July 6, 2022

- [Integrating Tracetest with GitHub Actions in a CI pipeline](https://kubeshop.io/blog/integrating-tracetest-with-github-actions-in-a-ci-pipeline) - June 24, 2022

And our recorded live stream: [Introduction to Tracetest - E2E Tests Powered by OpenTelemetry](https://youtu.be/mqwJRxqBNCg) - June 23, 2022


- [Tracetest — Assertions, Versioning & CI/CD - Release 0.5](https://kubeshop.io/blog/tracetest-assertions-versioning-ci-cd) - June 8, 2022

- [Tracing the History of Distributed Tracing & OpenTelemetry](https://kubeshop.io/blog/tracing-the-history-of-distributed-tracing-opentelemetryt) - May 26, 2022

- [Tracetest is released. What’s next?](https://kubeshop.io/blog/tracetest-is-released-whats-next) - May 6, 2022

- [Introducing Tracetest - Trace-based Testing with OpenTelemetry](https://kubeshop.io/blog/introducing-tracetest-trace-based-testing-with-opentelemetry) - April 26, 2022
