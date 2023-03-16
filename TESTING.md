# Testing Tracetest

On Tracetest, we work in two ways to test the entire system and guarantee that everything is working fine. One is through automatic tests, where we run unit, integration and end-to-end tests on Tracetest CLI, Web UI and Server.

Another source of tests is the manual tests that we execute on each release, simulating users, and checking if everything is ok on these tests.

## Manual Tests

On our manual tests, we aim to do some [sanity checks](https://en.wikipedia.org/wiki/Sanity_check) to see if the main features are working as expected. Usually, we run it on each Tracetest release.

### Testing Tracetest setup

This is a simple test to check if Tracetest is working correctly given it was provisioned with one Tracing Data Store.

The steps that we should follow are:

- [ ] Open WebUI and go to `/settings` page. The provisioned Data Store should be selected.
- [ ] Run `tracetest datastore export -d {provisioned_datastore}` and check if the data was exported correctly.
- [ ] Create a test on WebUI that calls a demo API (like [Pokeshop](https://docs.tracetest.io/live-examples/pokeshop/overview) or [Open Telemetry Store](https://docs.tracetest.io/live-examples/opentelemetry-store/overview)). This test should fetch traces correctly and run without errors.

### Checklist on version release

This is the entire checklist on what we should do to assert that Tracetest is working fine on each version release. On each version release, we can copy the contents of this checklist and open a Github Discussion to start each test.

- [ ] Check if our release pipeline on [Release Tracetest](https://github.com/kubeshop/tracetest/actions/workflows/release-version.yml) workflow on Github Actions worked correctly.
- [ ] Test CLI update on MacOS via homebrew
- [ ] Test CLI update on MacOS via curl script
- [ ] Test CLI update on Linux via APT
- [ ] Test CLI update on Linux via YUM
- [ ] Test CLI update on Linux via curl script
- [ ] Test CLI update on Windows via chocolatey
- [ ] Test CLI update on Windows via manual download
- [ ] Test server installation via CLI with Docker Compose and no demo API
- [ ] Test server installation via CLI with Docker Compose and demo API
- [ ] Test server installation via CLI with Kubernetes and no demo API
- [ ] Test server installation via CLI with Kubernetes and demo API
- [ ] Double check [Detailed installation](https://docs.tracetest.io/getting-started/detailed-installation) doc and see if everything is documented correctly

- [ ] Test Tracetest setup with [Amazon X-Ray example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-amazon-x-ray)
- [ ] Test Tracetest setup with [Datadog example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-datadog)
- [ ] Test Tracetest setup with [Elastic APM example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-elasticapm)
- [ ] Test Tracetest setup with [Lightstep example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-lightstep)
- [ ] Test Tracetest setup with [New Relic example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-new-relic)
- [ ] Test Tracetest setup with [SignalFX example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-signalfx)

- [ ] Test specific features added/changed on this release on a local installation or in a Kubernetes installation:
  - [ ] Merge https://github.com/kubeshop/helm-charts/pull/436

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
