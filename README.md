<p align="center">
  <img style="width:66%" src="assets/tracetest-logo-color-w-white-text.svg#gh-dark-mode-only" alt="Tracetest Logo Light"/>
  <img style="width:66%" src="assets/tracetest-logo-color-w-black-text.svg#gh-light-mode-only" alt="Tracetest Logo Dark" />
</p>

<p align="center">
  Tracetest - the best way to develop and test your distributed system with OpenTelemetry. For QA, Dev, & Ops.
</p>

<p align="center">
  <!--<a href="https://tracetest.io">Website</a>&nbsp;|&nbsp; -->
  <!--<a href="https://github.com/kubeshop/tracetest#try-the-demo--give-us-feedback">Live Demo</a>&nbsp;|&nbsp;-->
  <a href="https://kubeshop.github.io/tracetest/installing/">Install</a>&nbsp;|&nbsp;
  <a href="https://kubeshop.github.io/tracetest">Documentation</a>&nbsp;|&nbsp;
  <a href="https://twitter.com/tracetest_io">Twitter</a>&nbsp;|&nbsp;
  <a href="https://discord.gg/eBvEQRVyKX">Discord</a>&nbsp;|&nbsp;
  <a href="https://tracetest.io/blog">Blog</a>
</p>

<p align="center">
  <a href="https://github.com/kubeshop/tracetest/releases"><img title="Release" src="https://img.shields.io/github/v/release/kubeshop/tracetest"/></a>
  <a href=""><img title="Docker builds" src="https://img.shields.io/docker/automated/kubeshop/tracetest"/></a>
  <a href="https://github.com/kubeshop/tracetest/releases"><img title="Release date" src="https://img.shields.io/github/release-date/kubeshop/tracetest"/></a>
</p>

<p align="center">
  <a target="_new" href="https://kubeshop.wistia.com/medias/dw06408oqz">
    <img src="/assets/introvideo.png" style="width:66%;height:auto">
    <p align="center">
      Click on the image or this link to watch the "Tracetest Intro Video" video (< 2 minutes)
    </p>
  </a>
</p>

# Tracetest

Tracetest is a OpenTelemetry based tool that helps you develop and test your distributed applications. It assists you in the development process by enabling you to trigger your code and see the trace as you add OTel instrumentation. It also empowers you to create trace-based tests based on the data contained in your OpenTelemetry trace. You can verify against both the triggering transactions response AND any of the information contained deep in a span in your trace. This can include:

- verify the quality of your OpenTelemetry instrumentation and enforce standards.
- Testing events that occur on 'the other side' of an async message queue, even though the original async call has returned earlier.
- Assertions based on the timing of different steps in your process.
- Wildcard assertions across common types of activities, ie all gRPC return codes should be 0, all database calls should happen in less than 100ms.
- Testing long running processes instrumented with OpenTelemetry tracing to assert proper operation deep in the process.

# Features

- Works out of the box with your existing OTel instrumentation, supporting [numerous backend trace datastores](https://docs.tracetest.io/supported-backends/), including Jeager and Grafana Tempo. In addition, supports adding Tracetest as an [additional pipeline](https://docs.tracetest.io/supported-backends/#using-tracetest-without-a-backend) via your OpenTelemetry Collector config. Tell us others backend datastores you want supported!
- Supporting multiple ways of creating a test, including HTTP, GRPC and Postman Collections.
- Visualize the changes you are making to your trace as you develop, enabling Observability-Driven Development.
- [Add assertions](https://docs.tracetest.io/adding-assertions/) based on return data from trigger call and/or data contained in the spans in your distributed trace.
- Specifying which spans to check in assertions via the [advanced selector language](https://docs.tracetest.io/advanced-selectors/).
- Defining checks against the attributes in these spans, including properties, return status, or timing.
- Tests can be created via graphical UI or via [YAML-based test definition file](https://docs.tracetest.io/test-definition-file/).
- Use the test definition file to [enable Gitops flows](https://tracetest.io/blog/integrating-tracetest-with-github-actions-in-a-ci-pipeline).
- [Tracetest CLI](https://docs.tracetest.io/command-line-tool/) allows importing & exporting tests, running tests, and more.
- Tests are [versioned](https://docs.tracetest.io/versioning/) as the definition of the test is altered.
- Install can include [an example microservice](https://kubeshop.github.io/tracetest/pokeshop/) that is instrumented with OpenTelemetry to use as an example application under test.

# Getting Started

You can install tracetest by running:

```sh
curl -L https://raw.githubusercontent.com/kubeshop/tracetest/main/install-cli.sh | bash -s
tracetest server install
```

> :gear: To customize your Tracetest installation. Go to our [installation guide](https://docs.tracetest.io/installing/) for more information.

Installation only takes a few minutes and is done with via a Helm command. After installing, take a look at the
[Accessing the Dashboard](https://docs.tracetest.io/accessing-dashboard/) guide to access the Tracetest Dashboard and
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

The Tracetest [test definition files](https://docs.tracetest.io/test-definition-file/) are written in a simple YAML format. You can write them directly or build them graphically via the UI. Here is an example of a test which:

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
specs:
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

Is available at [https://docs.tracetest.io/](https://docs.tracetest.io/)

# Tests

We strive to produce quality code and improve Tracetest rapidly and safely. Therefore, we have a full suite of both frontend and backend tests. We are using Cypress to test our frontend code and (surprise, surprise) Tracetest for our backend code. You can see the [test runs here](https://github.com/kubeshop/tracetest/actions/workflows/pull-request.yaml), and a blog post describing our [testing pipelines here](https://tracetest.io/blog/integrating-tracetest-with-github-actions-in-a-ci-pipeline).

# License

[MIT License](/LICENSE)
