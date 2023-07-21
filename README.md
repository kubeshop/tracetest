<a href="https://tracetest.io">
  <p align="center">
    <img style="width:66%" src="assets/tracetest-logo-color-w-white-text.svg#gh-dark-mode-only" alt="Tracetest Logo Light"/>
    <img style="width:66%" src="assets/tracetest-logo-color-w-black-text.svg#gh-light-mode-only" alt="Tracetest Logo Dark" />
  </p>
</a>

---

<p align="center">
  Build integration and end-to-end tests in minutes using OpenTelemetry and trace-based testing.
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

Tracetest allows you to use Observability-Driven Development to build integration and end-to-end tests 98% faster with disitrbuted traces. No plumbing, no mocks, just real data.

You can:

- Assert against your trace data at every point of a request transaction.
- Assert on the timing of trace spans. Eg. A database span executes within `100ms`.
- Wildcard assertions across common types of activities. Eg. All gRPC return codes should be `0`, all database calls should happen in less than `100ms`.
- Ensure both your response and the underlying processes worked correctly, quickly, and without errors.
- Test any side-effects in your distributed system. Eg. Message queue, async API calls, external APIs, etc.
- Work with your existing distributed tracing solution.
- Define multiple transaction triggers:
  - HTTP requests
  - GRPC requests
  - trace IDs
  - and many more...
- Save and run the tests manually or via CI build jobs.
- Write detailed trace-based tests as:
  - End-to-end tests
  - Integration tests
- Verify and analyze the quality of your OpenTelemetry instrumentation to enforce rules and standards.
- Testing long running processes

Build tests in minutes.

- Visually - in the Web UI
  ![Build tests visually](https://res.cloudinary.com/djwdcmwdz/image/upload/v1688476657/docs/screely-1688476653521_omxe4r.png)

- Programmatically - in YAML

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

- Works out of the box with your existing OTel instrumentation, supporting [numerous backend trace datastores](https://docs.tracetest.io/getting-started/supported-backends), including Jeager and Grafana Tempo. In addition, supports adding Tracetest as an [additional pipeline](https://docs.tracetest.io/getting-started/supported-backends#using-tracetest-without-a-backend) via your OpenTelemetry Collector config. Tell us others backend datastores you want supported!
- Supporting multiple ways of creating a test, including HTTP, GRPC and Postman Collections.
- Visualize the changes you are making to your trace as you develop, enabling Observability-Driven Development.
- [Add assertions](https://docs.tracetest.io/using-tracetest/adding-assertions) based on return data from trigger call and/or data contained in the spans in your distributed trace.
- Specify which spans to check in assertions via the [selector language](https://docs.tracetest.io/concepts/selectors).
- Define checks against the attributes in these spans, including properties, return status, or timing.
- Create tests via graphical UI or via [YAML-based test definition file](https://docs.tracetest.io/cli/test-definition-file).
- Use the test definition file to [enable Gitops flows](https://tracetest.io/blog/integrating-tracetest-with-github-actions-in-a-ci-pipeline).
- [Tracetest CLI](https://docs.tracetest.io/cli/command-line-tool) allows importing & exporting tests, running tests, and more.
- [Version tests](https://docs.tracetest.io/using-tracetest/versioning) as the definition of the test is altered.
- Install can include [an example microservice](https://kubeshop.github.io/tracetest/pokeshop/) that is instrumented with OpenTelemetry to use as an example application under test.

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

This command will launch an install wizard that automatically installs Tracetest and a sample Pokeshop Demo app into either Docker or Kubernetes based on your selection.

Or, install Tracetest with Helm. The Tracetest Helm charts are located [here](https://github.com/kubeshop/helm-charts/tree/main/charts/tracetest).

```bash
helm install tracetest kubeshop/tracetest --namespace=tracetest --create-namespace
```

> [Read the Server installation docs for more options and instructions.](https://docs.tracetest.io/getting-started/installation#install-the-tracetest-server)

## 3️⃣ Open Tracetest

Once you've installed Tracetest Server, access the Tracetest Web UI on `http://localhost:11633`.

Check out the [Opening Tracetest guide](https://docs.tracetest.io/getting-started/open) to access the Tracetest Dashboard and start creating and running tests!

# 🤔 How does Tracetest work?

1. Pick an endpoint to test.
2. Run a test, and get the trace.
3. The trace is the blueprint of your system under test. It shows all the steps the system has taken to execute the request.
4. Use this blueprint to define assertions through Tracetest UI.
5. Add assertions on different services, checking return statuses, data, or even execution times of a system.
6. Run the tests.

Once the test is built, it can be run automatically as part of a build process. Every test has a trace attached, allowing you to immediately see what worked, and what did not, reducing the need to reproduce the problem to see the underlying issue.

# 📂 What does the test definition file look like?

The Tracetest [test definition files](https://docs.tracetest.io/cli/test-definition-file) are written in a simple YAML format. You can write them directly or build them graphically via the UI. Here is an example of a test which:

- executes POST against the pokemon/import endpoint.
- verifies that the HTTP blocks return a 200 status code.
- verifies all database calls execute in less than 50ms.

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

Save a test definition file as `pokeshop_import.yaml`. Use the CLI to run a test.

```bash
tracetest run test --file pokeshop_import.yaml
```

# Tests

We strive to produce quality code and improve Tracetest rapidly and safely. Therefore, we have a full suite of both frontend and backend tests. We are using Cypress to test our frontend code and (surprise, surprise) Tracetest for our backend code. You can see the [test runs here](https://github.com/kubeshop/tracetest/actions/workflows/pull-request.yaml), and a blog post describing our [testing pipelines here](https://tracetest.io/blog/integrating-tracetest-with-github-actions-in-a-ci-pipeline).

# 🎤 Feedback

Have an idea to improve Tracetest?

You can:

- [Create an issue here](https://github.com/kubeshop/tracetest/issues/new/choose)!
- Visit our [Community Page](https://tracetest.io/community).
- Join our [Discord](https://discord.gg/eBvEQRVyKX), and ask us any questions there.
- Follow us on [Twitter at @tracetest_io](https://twitter.com/tracetest_io) for updates.
- Give us a ⭐️ on Github if you're interested in the project!

# 🌱 Contributing & Community

Whether it's big or small, we love contributions.

Not sure where to get started? You can:

- Visit our [Community Page](https://tracetest.io/community).
- See our contributing guidelines [here](./CONTRIBUTING.md).
