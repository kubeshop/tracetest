# Testing Tracetest

On Tracetest, we work in two ways to test the entire system and guarantee that everything is working fine. One is through automatic tests, where we run unit, integration and end-to-end tests on Tracetest CLI, Web UI and API server.

Another source of tests are the manual tests that we execute on each release, simulating users, and checking if everything is ok on these tests.

## Automatic Tests

todo

## Manual Tests

### Checklist on version release

- [] Check if our release pipeline on [Release Tracetest](https://github.com/kubeshop/tracetest/actions/workflows/release-version.yml) workflow on Github Actions worked correctly.
- [] Test CLI update on MacOS via homebrew
- [] Test CLI update on MacOS via curl script
- [] Test CLI update on Linux via APT
- [] Test CLI update on Linux via YUM
- [] Test CLI update on Linux via curl script
- [] Test CLI update on Windows via chocolatey
- [] Test CLI update on Windows via manual download
- [] Test server installation via CLI with Docker Compose and no demo API
- [] Test server installation via CLI with Docker Compose and demo API
- [] Test server installation via CLI with Kubernetes and no demo API
- [] Test server installation via CLI with Kubernetes and demo API
- [] Double check [Detailed installation](https://docs.tracetest.io/getting-started/detailed-installation) doc and see if everything is documented correctly

- [] Test Tracetest setup with [No tracing example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-no-tracing)
- [] Test Tracetest setup with [Amazon X-Ray example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-amazon-x-ray)
- [] Test Tracetest setup with [Datadog example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-datadog)
- [] Test Tracetest setup with [Elastic APM example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-elasticapm)
- [] Test Tracetest setup with [Jaeger example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-jaeger)
- [] Test Tracetest setup with [Lightstep example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-lightstep)
- [] Test Tracetest setup with [New Relic example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-new-relic)
- [] Test Tracetest setup with [OpenSearch example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-opensearch)
- [] Test Tracetest setup with [SignalFX example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-signalfx)
- [] Test Tracetest setup with [Tempo example](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-tempo)
