# Testing Tracetest

On Tracetest, we work in two ways to test the entire system and guarantee that everything is working fine. One is through automatic tests, where we run unit, integration and end-to-end tests on Tracetest CLI, Web UI and Server.

Another source of tests is the manual tests that we execute on each release, simulating users, and checking if everything is ok on these tests.

## Manual Tests

On our manual tests, we aim to do some [sanity checks](https://en.wikipedia.org/wiki/Sanity_check) to see if the main features are working as expected. Usually, we run it on each Tracetest release.

### Testing Tracetest setup

This is a simple test to check if Tracetest is working correctly given it was provisioned with one Tracing Data Store.

The steps that we should follow are:

- [ ] Open WebUI and go to `/settings` page. The provisioned Data Store should be selected.
- [ ] Run `tracetest export datastore --id current` and check if the data was exported correctly.
- [ ] Create a test on WebUI that calls a demo API (like [Pokeshop](https://docs.tracetest.io/live-examples/pokeshop/overview) or [Open Telemetry Store](https://docs.tracetest.io/live-examples/opentelemetry-store/overview)). This test should fetch traces correctly and run without errors.

### Checklist on version release

This is the entire checklist on what we should do to assert that Tracetest is working fine on each version release. On each version release, we can copy the contents of this checklist and open a Github Discussion to start each test.

- [ ] Check if our release pipeline on [Release Tracetest](https://github.com/kubeshop/tracetest/actions/workflows/release-version.yml) workflow on Github Actions worked correctly.
- [ ] Double check [Detailed installation](https://docs.tracetest.io/getting-started/detailed-installation) doc and see if everything is documented correctly

### Tests to validate RC

- Test server installation via CLI

  - [ ] Docker Compose and no demo API
  - [ ] Docker Compose and demo API
  - [ ] Kubernetes and no demo API
  - [ ] Kubernetes and demo API

- Test Tracetest examples

  - [ ] [Amazon X-Ray example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-amazon-x-ray)
  - [ ] [Datadog example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-datadog)
  - [ ] [Elastic APM example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-elasticapm)
  - [ ] [Lightstep example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-lightstep)
  - [ ] [New Relic example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-new-relic)
  - [ ] [SignalFX example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-signalfx)

- Test specific features added/changed on this release:

  - [ ] Add features here

  ### Tests to validate final release

- Test CLI updates

  - [ ] MacOS via homebrew
  - [ ] MacOS via curl script
  - [ ] Windows via chocolatey
  - [ ] Windows via manual download

- Test specific features added/changed on this release:
  - [ ] Add features here

## Automatic Tests

Today Tracetest has 3 main components: a WebUI, a CLI and a Server.

### Web UI

- **Unit tests**: Run by executing `npm test` on `./web` folder
- **End-to-end tests**: Run using [cypress](https://www.cypress.io/) against a temporary Tracetest created on Kubernetes.

### CLI

- **Unit tests**: Run by executing `make test` on `./cli` folder

### Server

- **Unit tests**: Run by executing `make test` on `./server` folder
- **Integration tests**: Run along with some unit tests running `make test` on `./server` folder, but also done in a matrix test on Github actions, by executing the `test-examples` action.
- **End-to-end tests**: Run using Tracetest to test itself. We deploy two instances of Tracetest and use one to test API calls to another on Github actions, by executing the `trace-testing` action.
