<p align="center">  
  <img style="width:66%" src="assets/tracetest-logo-color-w-white-text.svg#gh-dark-mode-only" alt="Tracetest Logo Light"/>
  <img style="width:66%" src="assets/tracetest-logo-color-w-black-text.svg#gh-light-mode-only" alt="Tracetest Logo Dark" />
</p>

<p align="center">
  End-to-end tests powered by OpenTelemetry. For QA, Dev, & Ops.
</p>

<p align="center">
  <!--<a href="https://tracetest.io">Website</a>&nbsp;|&nbsp; -->
  <!--<a href="https://github.com/kubeshop/tracetest#try-the-demo--give-us-feedback">Live Demo</a>&nbsp;|&nbsp;-->
  <a href="https://kubeshop.github.io/tracetest/installing/">Install</a>&nbsp;|&nbsp;
  <a href="https://kubeshop.github.io/tracetest">Documentation</a>&nbsp;|&nbsp; 
  <a href="https://twitter.com/tracetest_io">Twitter</a>&nbsp;|&nbsp; 
  <a href="https://discord.gg/eBvEQRVyKX">Discord</a>&nbsp;|&nbsp; 
  <a href="https://kubeshop.io/blog-projects/tracetest">Blog</a>
</p>

<p align="center">
  <a href="https://github.com/kubeshop/tracetest/releases"><img title="Release" src="https://img.shields.io/github/v/release/kubeshop/tracetest"/></a>
  <a href=""><img title="Docker builds" src="https://img.shields.io/docker/automated/kubeshop/tracetest"/></a>
  <a href="https://github.com/kubeshop/tracetest/releases"><img title="Release date" src="https://img.shields.io/github/release-date/kubeshop/tracetest"/></a>
</p>

<p align="center">
  <a target="_new" href="https://www.youtube.com/watch?v=WMRicNlaehc">
    <img src="https://img.youtube.com/vi/WMRicNlaehc/0.jpg" style="width:66%;height:auto">
    <p align="center">
      Click on the image or this link to watch the "Tracetest Intro Video" video (7 mins)
    </p>
  </a>
</p>

# Tracetest

Tracetest is a trace-based testing tool that leverages the data contained in your distributed traces to produce easy to create, yet super powerful integration tests. You can verify activity deep inside your system, even events that occur on 'the other side' of an async message queue. Verifying the timing of different steps in your process via an automated test is also a valuable use case.

# Features

- Test by executing a REST or gRPC call to trigger the test. Can also import your Postman Collections.
- [Add assertions](https://kubeshop.github.io/tracetest/adding-assertions/) based on return data from trigger call and/or data contained in the spans in your distributed trace.
- Enables white box testing in which internal structure, design and coding of software are tested to verify flow of input-output and to improve design, usability and security.
- Specify which spans to check in assertions via the [advanced selector language](https://kubeshop.github.io/tracetest/advanced-selectors/).
- Define checks against the attributes in these spans, including properties, return status, or timing.
- Tests can be created via graphical UI or via [YAML-based test definition file](https://kubeshop.github.io/tracetest/test-definition-file/).
- Use the test definition file to [enable Gitops flows](https://kubeshop.io/blog/integrating-tracetest-with-github-actions-in-a-ci-pipeline).
- [Tracetest CLI](https://kubeshop.github.io/tracetest/command-line-tool/) allows importing & exporting tests, running tests, and more.
- Tests are [versioned](https://kubeshop.github.io/tracetest/versioning/) as the definition of the test is altered.
- Supports [numerous backend trace datastores](https://kubeshop.github.io/tracetest/architecture/), including Jeager and Grafana Tempo. Tell us which others you want!
- Easy [install via Helm command](https://kubeshop.github.io/tracetest/installing/).
- Install can include [an example microservice](https://kubeshop.github.io/tracetest/pokeshop/) that is instrumented with OpenTelemetry to use as an example application under test.

# How does Tracetest work?

1. Pick an endpoint to test.
2. Run a test, and get the trace.
3. The trace is the blueprint of your system under test. It shows all the steps the system has taken to execute the request.
4. Use this blueprint to define assertions through Tracetest UI.
5. Add assertions on different services, checking return statuses, data, or even execution times of a system.
6. Run the tests.

Once the test is built, it can be run automatically as part of a build process. Every test has a trace attached, allowing you to immediately see what worked, and what did not, reducing the need to reproduce the problem to see the underlying issue.

# Getting Started

The [install](https://kubeshop.github.io/tracetest/installing/) only takes a few minutes, and is done with via a Helm command. After installing, take a look at the
[Accessing the Dashboard](https://kubeshop.github.io/tracetest/accessing-dashboard/) guide to access the Tracetest Dashboard and
create and run your first test.

# What does the test definition file look like?

The Tracetest [test definition files](https://kubeshop.github.io/tracetest/test-definition-file/) are written in a simple YAML format. You can write them directly or build them graphically via the UI. Here is an example of a test which:

- executes POST against the pokemon/import endpoint.
- verifies that the HTTP blocks return a 200 status code.
- verifies all database calls execute in less than 200ms.

```yaml
id: 5dd03dda-fad2-49f0-b9d9-5143b746c1d0
name: DEMO Pokemon - Import - Import a Pokemon
description: "Import a pokemon"

# Configure how tracetest triggers the operation on your application
trigger:
    type: http
    httpRequest:
        method: POST
        url: http://demo-pokemon-api.demo.svc.cluster.local/pokemon/import
        headers:
            - key: Content-Type
              value: application/json
        body: '{"id":52}'

# Setup what will be asserted on the resulting trace
testDefinition:
    selector: span[tracetest.span.type = "http"]
    - assertions:
        - http.status_code = 200
    selector: span[tracetest.span.type = "database"]
    - assertions:
        - tracetest.span.duration < "50ms"
```

# Feedback

We are in the our early days with the project and need your help. Have an idea to improve it? Please [Create an issue here](https://github.com/kubeshop/tracetest/issues/new/choose) or join our community on [Discord](https://discord.gg/eBvEQRVyKX).

Follow us on [Twitter at @tracetest_io](https://twitter.com/tracetest_io) for updates

Give us a star on Github if you're interested in the project!

# Documentation

Is available at [https://kubeshop.github.io/tracetest](https://kubeshop.github.io/tracetest)

# Tests

We strive to produce quality code and improve Tracetest rapidly and safely. Therefore, we have a full suite of both frontend and backend code. We are using Cypress for our frontend code and (surprise surprise) Tracetest for our backend code. You can see the [test runs here](https://github.com/kubeshop/tracetest/actions/workflows/pull-request.yaml), and a blog post describing our [testing pipelines here](https://kubeshop.io/blog/integrating-tracetest-with-github-actions-in-a-ci-pipeline).

# License

[MIT License](/LICENSE)
