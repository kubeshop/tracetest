<a href="https://tracetest.io">
  <p align="center">
    <picture>
      <source media="(prefers-color-scheme: light)" srcset="assets/tracetest-logo-color-w-black-text.svg">
      <source media="(prefers-color-scheme: dark)" srcset="assets/tracetest-logo-color-w-white-text.svg">
      <img alt="Tracetest Logo" src="assets/tracetest-logo-color-w-black-text.svg" style="max-width:450px">
    </picture>
  </p>
</a>

---

<p align="center">
  Build integration and end-to-end tests in minutes, instead of days, using OpenTelemetry and trace-based testing.
</p>

<b>
  <p align="center">
    <a href="https://docs.tracetest.io/getting-started/installation">
      Get Started! &nbsp;👉&nbsp;
    </a>
  </p>
</b>

<b>
  <p align="center">
    <a href="https://docs.tracetest.io/">Docs</a>&nbsp;|&nbsp;
    <a href="https://docs.tracetest.io/examples-tutorials/overview">Tutorials</a>&nbsp;|&nbsp;
    <a href="https://docs.tracetest.io/examples-tutorials/recipes">Recipes</a>&nbsp;|&nbsp;
    <a href="https://github.com/kubeshop/tracetest/tree/main/examples">Examples</a>&nbsp;|&nbsp;
    <a href="https://discord.gg/eBvEQRVyKX">Discord</a>&nbsp;|&nbsp;
    <a href="https://tracetest.io/blog">Blog</a>&nbsp;|&nbsp;
    <a href="https://tracetest.io">Website</a>
  </p>
</b>

<h4 align="center">
   <a href="https://github.com/kubeshop/tracetest/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/Licence-MIT-blue" alt="Tracetest is released under the MIT License">
  </a>
  <a href="https://kubeshop.io/"><img src="https://img.shields.io/website?url=https://kubeshop.io/&up_message=Kubeshop&up_color=%232635F1&label=Accelerator&down_color=%232635F1&down_message=Kubeshop"></a>
  <a href="https://github.com/kubeshop/tracetest/releases">
    <img title="Release" src="https://img.shields.io/github/v/release/kubeshop/tracetest"/>
  </a>
  <a href="https://github.com/kubeshop/tracetest/releases">
    <img title="Release date" src="https://img.shields.io/github/release-date/kubeshop/tracetest"/>
  </a>
  <a href="https://github.com/kubeshop/tracetest/blob/main/CONTRIBUTING.md">
    <img src="https://img.shields.io/badge/PRs-Welcome-brightgreen?logo=github" alt="PRs welcome!" />
  </a>
  <a href="https://github.com/kubeshop/tracetest/issues">
    <img src="https://img.shields.io/github/stars/kubeshop/tracetest?color=%23EAC54F&logo=github&label=Help us reach 1k stars! Now at:" alt="Help us reach 1k stars!" />
  </a>
  <a href="https://tracetest.io/community">
    <img src="https://img.shields.io/badge/Join-Community!-purple" alt="Join our Community!" />
  </a>
  <a href="https://discord.gg/8MtcMrQNbX">
    <img src="https://img.shields.io/badge/Chat-on Discord!-red" alt="Talk to us on Discord!" />
  </a>
  <a href="https://hub.docker.com/r/kubeshop/tracetest"><img title="Docker Pulls" src="https://img.shields.io/docker/pulls/kubeshop/tracetest?logo=docker"/></a>
  <a href="https://hub.docker.com/r/kubeshop/tracetest/tags?page=1&name=latest"><img title="Docker Image Size" src="https://img.shields.io/docker/image-size/kubeshop/tracetest/latest?logo=docker&label='kubeshop/tracetest:latest' image size"/></a>
  <a href="https://twitter.com/tracetest_io">
    <img src="https://img.shields.io/badge/follow-%40tracetest__io-1DA1F2?logo=twitter&style=social" alt="Tracetest Twitter" />
  </a>
</h4>

<p align="center">
  Tracetest lets you build integration and end-to-end tests 98% faster with disitrbuted traces.
</p>

<p align="center">
  No plumbing, no mocks, no fakes. Test against real data.
</p>

You can:

- **Assert** against both the **response and trace data** at every point of a request transaction.
- **Assert** on the **timing of trace spans**.
  - Eg. A database span executes within `100ms`.
- **Wildcard assertions** across common types of activities.
  - Eg. All gRPC return codes should be `0`.
  - Eg. All database calls should happen in less than `100ms`.
- **Assert** against **side-effects** in your distributed system.
  - Eg. Message queues, async API calls, external APIs, etc.
- **Integrate** with your **existing distributed tracing solution**.
- Define multiple test triggers:
  - HTTP requests
  - gRPC requests
  - Trace IDs
  - and many more...
- Save and run tests manually and via CI build jobs.
- Verify and analyze the quality of your OpenTelemetry instrumentation to enforce rules and standards.
- Test long-running processes.

**Build tests in minutes**.

- **Visually** - in the Web UI
  <p align="center">
    <img src="https://res.cloudinary.com/djwdcmwdz/image/upload/v1688476657/docs/screely-1688476653521_omxe4r.png" style="width:66%;height:auto">
  </p>

- **Programmatically** - in YAML

  ```yaml
  type: Test
  spec:
    id: Yg9sN-94g
    name: Pokeshop - Import
    description: Import a Pokemon
    trigger:
      type: http
      httpRequest:
        url: http://demo-api:8081/pokemon/import
        method: POST
        headers:
        - key: Content-Type
          value: application/json
        body: '{"id":52}'
    specs:
    - name: 'All Database Spans: Processing time is less than 100ms'
      selector: span[tracetest.span.type="database"]
      assertions:
      - attr:tracetest.span.duration < 100ms
  ```

# 🔥 Features

- Works out of the box with your existing OpenTelemetry instrumentation, supporting [numerous trace data stores](https://docs.tracetest.io/configuration/overview/#supported-trace-data-stores), including:
  - Jeager
  - Grafana Tempo
  - OpenSearch
  - Elastic
  - And, many more...
  - Tell us which other trace data stores you want supported!
- Works out of the box by adding Tracetest as an [additional pipeline](https://docs.tracetest.io/getting-started/supported-backends#using-tracetest-without-a-backend) via your OpenTelemetry Collector config.
- Supporting multiple ways of creating a test, including HTTP, GRPC and Postman Collections.
- Visualize the changes you are making to your trace as you develop, enabling Observability-Driven Development.
- [Add assertions](https://docs.tracetest.io/using-tracetest/adding-assertions) based on response data from the trigger request and all trace data contained in the spans of your distributed trace.
- Specify which spans to check in assertions via the [selector language](https://docs.tracetest.io/concepts/selectors).
- Define checks against the attributes in these spans, including properties, return status, or timing.
- Create tests visually in the Tracetest Web UI or programatically via [YAML-based test definition files](https://docs.tracetest.io/cli/test-definition-file).
- Use test definition files and the Tracetest CLI to [enable Gitops flows and CI/CD automation](https://docs.tracetest.io/ci-cd-automation/overview).
- [Tracetest CLI](https://docs.tracetest.io/cli/cli-installation-reference) allows importing & exporting tests, running tests, and more.
- [Version tests](https://docs.tracetest.io/concepts/versioning/) as the definition of the test is altered.
- The [guided install](https://docs.tracetest.io/getting-started/installation) can include [an example "Pokeshop" microservice](https://docs.tracetest.io/live-examples/pokeshop/overview) that is instrumented with OpenTelemetry to use as an example application under test.
- Create [environment variables](https://docs.tracetest.io/concepts/environments) to assert the same behavior across multiple environments (dev, staging, and production, for example)
- Create [test outputs](https://docs.tracetest.io/web-ui/creating-test-outputs/) by defining a variable based on the information contained in a particular span's attributes.
- Run [ad-hoc tests](https://docs.tracetest.io/concepts/ad-hoc-testing) by using undefined variables to enable supplying variables at runtime.
- Define [test suites/transactions](https://docs.tracetest.io/concepts/transactions) to chain tests together and use variables obtained from a test in a subsequent test. These variables can also be loaded from the environment.
- Run comprehensive [trace analysis and validation](https://docs.tracetest.io/analyzer/concepts) to adhere to OpenTelemetry rules and standards.
- Configure [test runner](https://docs.tracetest.io/configuration/test-runner) behavior with required gates used when executing your tests to determine whether to mark the test as passed or failed.

# 🚀 Getting Started

<p align="center">
  <a target="_new" href="https://kubeshop.wistia.com/medias/dw06408oqz">
    <img src="/assets/introvideo.png" style="width:66%;height:auto">
    <p align="center">
      Click on the image or this link to watch the "Tracetest Intro Video" video (< 2 minutes)
    </p>
  </a>
</p>

## 1️⃣ Install the Tracetest CLI

```bash
curl -L https://raw.githubusercontent.com/kubeshop/tracetest/main/install-cli.sh | bash -s
```

> [:gear: Read the CLI installation docs for more options and instructions.](https://docs.tracetest.io/getting-started/installation#install-the-tracetest-cli)

## 2️⃣ Install the Tracetest Server

```bash
tracetest server install
```

This command will launch an install wizard that automatically installs Tracetest and a [sample Pokeshop Demo app](https://docs.tracetest.io/live-examples/pokeshop/overview) into either Docker or Kubernetes based on your selection.

Or, install Tracetest with Helm. The Tracetest Helm charts are located [here](https://github.com/kubeshop/helm-charts/tree/main/charts/tracetest).

```bash
helm install tracetest kubeshop/tracetest --namespace=tracetest --create-namespace
```

> [:gear: Read the Server installation docs for more options and instructions.](https://docs.tracetest.io/getting-started/installation#install-the-tracetest-server)

## 3️⃣ Open Tracetest

Once you've installed Tracetest Server, access the Tracetest Web UI on `http://localhost:11633`.

Check out the [Opening Tracetest guide](https://docs.tracetest.io/getting-started/open) to start creating and running tests!

# 🤔 How does Tracetest work?

1. Pick an endpoint to test.
2. Run a test, and get the trace.
3. The trace is the blueprint of your system under test. It shows all the steps the system has taken to execute the request.
4. Use this blueprint to define assertions in the Tracetest Web UI.
5. Add assertions on different services, checking return statuses, data, or even execution times of a system.
6. Run the tests.

Once the test is built, it can be run automatically as part of a build process. Every test has a trace attached, allowing you to immediately see what worked, and what did not, reducing the need to reproduce the problem to see the underlying issue.

# 📂 What does the test definition file look like?

The Tracetest [test definition files](https://docs.tracetest.io/cli/test-definition-file) are written in a simple YAML format. You can write them directly or build them graphically via the UI. Here is an example of a test which:

- Executes `POST` against the `pokemon/import` endpoint.
- Verifies that the HTTP blocks return a `200` status code.
- Verifies all database calls execute in less than `50ms`.

```yaml
type: Test
spec:
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

# 🤖 How to run an automated test?

Save a test definition file above as `pokeshop_import.yaml`. Use the CLI to run a test.

```bash
tracetest run test --file /path/to/pokeshop_import.yaml
```

Check out the [CI/CD docs](https://docs.tracetest.io/ci-cd-automation/overview) to learn more about test automation.

# Tests

We strive to produce quality code and improve Tracetest rapidly and safely. Therefore, we have a full suite of both front-end and back-end tests. Cypress tests are running against the front-end code and (surprise, surprise) Tracetest against the back-end code. You can see the [test runs here](https://github.com/kubeshop/tracetest/actions/workflows/pull-request.yaml), and a blog post describing our [testing pipelines here](https://tracetest.io/blog/50-faster-ci-pipelines-with-github-actions).

# 🎤 Feedback

Have an idea to improve Tracetest?

You can:

- [Create an issue here](https://github.com/kubeshop/tracetest/issues/new/choose)!
- Join our [Discord](https://discord.gg/eBvEQRVyKX), and ask us any questions there.
- Follow us on [Twitter at @tracetest_io](https://twitter.com/tracetest_io) for updates.
- Give us a ⭐️ on Github if you like what we're doing!

# 🌱 Contributing & Community

Whether it's big or small, we love contributions.

Not sure where to get started? You can:

- Visit our [Community Page](https://tracetest.io/community).
- See our contributing guidelines [here](./CONTRIBUTING.md).
