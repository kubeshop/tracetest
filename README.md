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

Tracetest is a trace-based testing tool that leverages the data captured by your existing Open Telemetry distributed traces to produce easy to create, yet super powerful integration tests. You can verify activity deep inside your system by asserting on data and flow information contained in the OpenTelemetry traces and span attributes. This can include:

- verify the quality of your OpenTelemetry instrumentation and enforce standards.
- Testing events that occur on 'the other side' of an async message queue, even though the original async call has returned earlier.
- Assertions based on the timing of different steps in your process.
- Wildcard assertions across common types of activities, ie all gRPC return codes should be 0, all database calls should happen in less than 100ms.
- Testing long running processes instrumented with OpenTelemetry tracing to assert proper operation deep in the process.

# Features

- Supporting multiple ways of creating a test, including HTTP, GRPC and Postman Collections.
- [Adding assertions](https://kubeshop.github.io/tracetest/adding-assertions/) based on return data from trigger call and/or data contained in the spans in your distributed trace.
- Specifying which spans to check in assertions via the [advanced selector language](https://kubeshop.github.io/tracetest/advanced-selectors/).
- Defining checks against the attributes in these spans, including properties, return status, or timing.
- Tests can be created via graphical UI or via [YAML-based test definition file](https://kubeshop.github.io/tracetest/test-definition-file/).
- Use the test definition file to [enable Gitops flows](https://kubeshop.io/blog/integrating-tracetest-with-github-actions-in-a-ci-pipeline).
- [Tracetest CLI](https://kubeshop.github.io/tracetest/command-line-tool/) allows importing & exporting tests, running tests, and more.
- Tests are [versioned](https://kubeshop.github.io/tracetest/versioning/) as the definition of the test is altered.
- Supports [numerous backend trace datastores](https://kubeshop.github.io/tracetest/architecture/), including Jeager and Grafana Tempo. Tell us which others you want!
- Install can include [an example microservice](https://kubeshop.github.io/tracetest/pokeshop/) that is instrumented with OpenTelemetry to use as an example application under test.

# Getting Started

You can install tracetest by running:

```sh
curl -L https://raw.githubusercontent.com/kubeshop/tracetest/main/setup.sh | bash -s
```

> :gear: To customize your Tracetest installation. Go to our [installation guide](https://kubeshop.github.io/tracetest/installing/) for more information.

Installation only takes a few minutes and is done with via a Helm command. After installing, take a look at the
[Accessing the Dashboard](https://kubeshop.github.io/tracetest/accessing-dashboard/) guide to access the Tracetest Dashboard and
create and run your first test.

# How does Tracetest work?

1. Pick an endpoint to test.
2. Run a test, and get the trace.
3. The trace is the blueprint of your system under test. It shows all the steps the system has taken to execute the request.
4. Use this blueprint to define assertions through Tracetest UI.
5. Add assertions on different services, checking return statuses, data, or even execution times of a system.
6. Run the tests.

Once the test is built, it can be run automatically as part of a build process. Every test has a trace attached, allowing you to immediately see what worked, and what did not, reducing the need to reproduce the problem to see the underlying issue.

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
# triggers can be http, grpc, etc
trigger:
    type: http
    httpRequest:
        method: POST
        url: http://demo-pokemon-api.demo.svc.cluster.local/pokemon/import
        headers:
            - key: Content-Type
              value: application/json
        body: '{"id":52}'

# Definition of the test specs which is a combination of a selector
# and an assertion
testDefinition:
    # the selector defines which spans will be targeted by the assertions
    selector: span[tracetest.span.type = "http"]
    # the assertions define the checks to be run. In this case, all
    # http spans will be checked for a status code = 200
    - assertions:
        - http.status_code = 200
    # this next test ensures all the database spans execute in less
    # than 50 ms
    selector: span[tracetest.span.type = "database"]
    - assertions:
        - tracetest.span.duration < "50ms"
```

# Feedback

Have an idea to improve Tracetest? Please [create an issue here](https://github.com/kubeshop/tracetest/issues/new/choose) or join our community on [Discord](https://discord.gg/eBvEQRVyKX).

Follow us on [Twitter at @tracetest_io](https://twitter.com/tracetest_io) for updates.

Give us a star on Github if you're interested in the project!

# Documentation

Is available at [https://kubeshop.github.io/tracetest](https://kubeshop.github.io/tracetest)

# Tests

We strive to produce quality code and improve Tracetest rapidly and safely. Therefore, we have a full suite of both frontend and backend tests. We are using Cypress to test our frontend code and (surprise, surprise) Tracetest for our backend code. You can see the [test runs here](https://github.com/kubeshop/tracetest/actions/workflows/pull-request.yaml), and a blog post describing our [testing pipelines here](https://kubeshop.io/blog/integrating-tracetest-with-github-actions-in-a-ci-pipeline).

# License

[MIT License](/LICENSE)
