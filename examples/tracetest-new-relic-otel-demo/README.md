# OpenTelemetry Demo with Tracetest and New Relic

This examples uses OpenTelemetry Demo `v1.2.1`.

This is a production-ready sample app on how to configure the [OpenTelemetry Demo](https://github.com/open-telemetry/opentelemetry-demo) to use [Tracetest](https://tracetest.io/) for enhancing your E2E and integration tests with trace-based testing, and [New Relic](https://newrelic.com/) as a trace data store.

## Prerequisites

You will need [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) installed on your machine to run this sample app! Additionally, you will need a New Relic account and ingest licence key. Sign up to New Relic [here](https://newrelic.com/signup).

## Project structure

The project is built with Docker Compose. It contains a `docker-compose.yaml` file with 25 services.

### 1. OpenTelemetry Demo

The `docker-compose.yaml` file and `.env` file in the root directory are for the OpenTelemetry Demo.

### 2. Tracetest & New Relic

At the bottom of the `docker-compose.yaml` file you'll see the Tracetest service. In the `./otelcollector/otelcol-config-extras.yml` you'll see the config for forwarding traces to both Tracetest and New Relic. The `./tracetest/tracetest.config.yaml` is for the setting up Tracetest and the OpenTelemetry Collector.

The `tracetest` directory also contains an `e2e` directory with a `http-test.yaml` file which is a Tracetest test definition for running a test via the Tracetest CLI.

### Docker Compose Network

All `services` in the `docker-compose.yaml` are on the same network and will be reachable by hostname from within other services. E.g. `tracetest:4317` in the `otelcol-config-extras.yml` will map to the `tracetest` service, where the port `4317` is the port where Tracetest accepts traces.

## OpenTelemetry Demo

The [OpenDelemetry Demo](https://github.com/open-telemetry/opentelemetry-demo) is a sample microservice-based app with the purpose to demo how to correctly set up OpenTelemetry distributed tracing.

Read more about the OpenTelemetry Demo [here](https://opentelemetry.io/blog/2022/announcing-opentelemetry-demo-release/).

The `docker-compose.yaml` contains 25 services. View the file in its entirety [here](./docker-compose.yaml).

To start the OpenTelemetry Demo by itself, run this command:

```bash
docker compose up
```

> Note: Building the images locally is currently not supported in this example app.

This will start the OpenTelemetry Demo. Open up `http://localhost:8080` to make sure it's working. But, you're not sending the traces anywhere.

Let's fix this by configuring Tracetest and OpenTelemetry Collector to forward trace data to both New Relic and Tracetest.

## Tracetest

At the bottom of the `docker-compose.yaml` you'll see a `# Tracetest` comment. There you'll see two configured services.

- **Postgres** - Postgres is a prerequisite for Tracetest to work. It stores trace data when running the trace-based tests.
- [**Tracetest**](https://tracetest.io/) - Trace-based testing that generates end-to-end tests automatically from traces.

The `TRACETEST_SERVICE_PORT` is configured in the `.env` file

```yaml
# ...

  tracetest:
    image: kubeshop/tracetest:${TAG:-latest}
    ports:
      - "${TRACETEST_SERVICE_PORT}:${TRACETEST_SERVICE_PORT}"
    volumes:
      - type: bind
        source: ./tracetest/tracetest.config.yaml
        target: /app/config.yaml
      - type: bind
        source: ./tracetest/tracetest-provision.yaml
        target: /app/provision.yaml
    command: --provisioning-file /app/provision.yaml
    healthcheck:
      test: ["CMD", "wget", "--spider", "localhost:11633"]
      interval: 1s
      timeout: 3s
      retries: 60
    depends_on:
      tt_postgres:
        condition: service_healthy
      otelcol:
        condition: service_started
    logging: *logging

  # Postgres used by the Tracetest instance
  tt_postgres:
    image: postgres:14
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_USER: postgres
    healthcheck:
      test: pg_isready -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"
      interval: 1s
      timeout: 5s
      retries: 60
    logging: *logging
```

Tracetest depends on both Postgres and the OpenTelemetry Collector. Tracetest requires a config file to be loaded via a volume. The volume is mapped from the root directory into the `tracetest` directory and the respective config file.

The `tracetest.config.yaml` file contains the basic setup of connecting Tracetest to the Postgres instance, and defining the exporter. The exporter is set to the OpenTelemetry Collector.

```yaml
# tracetest.config.yaml
---
postgres:
  host: tt_postgres
  user: postgres
  password: postgres
  port: 5432
  dbname: postgres
  params: sslmode=disable

poolingConfig:
  maxWaitTimeForTrace: 30s
  retryDelay: 500ms

# This section will populate the Tracetest Web UI with sample tests for you to try out
demo:
  enabled: [otel]
  endpoints:
    otelFrontend: http://frontend:8080
    otelProductCatalog: productcatalogservice:3550
    otelCart: cartservice:7070
    otelCheckout: checkoutservice:5050

experimentalFeatures: []

googleAnalytics:
  enabled: true

telemetry:
  exporters:
    collector:
      serviceName: tracetest
      sampling: 100
      exporter:
        type: collector
        collector:
          endpoint: otelcol:4317

server:
  telemetry:
    exporter: collector
    applicationExporter: collector
```

The `tracetest-provision.yaml` file contains the data store setup. The data store is set to OTLP meaning the traces will be stored in Tracetest itself.

```yaml
# tracetest-provision.yaml
---
dataStore:
  type: otlp
```

**How to send traces to Tracetest and New Relic?**

The `otelcol-config-extras.yml` explains that. But first, check the `otelcol-config.yml`. It receives traces via either `grpc` or `http`. Then, in the `otelcol-config-extras.yml` you see a `exporters` that exports traces to Tracetest's OTLP endpoint `tracetest:4317` in one pipeline, and to New Relic in another.

Make sure to add your New Relic access token in the headers of the `otlp/newrelic` exporter.

```yaml
# otelcol-config-extras.yml

# extra settings to be merged into OpenTelemetry Collector configuration
# do not delete this file

processors:
  batch:
    timeout: 100ms

exporters:
  # OTLP for Tracetest
  otlp/tracetest:
    endpoint: tracetest:4317 # Send traces to Tracetest. Read more in docs here:  https://docs.tracetest.io/configuration/connecting-to-data-stores/opentelemetry-collector
    tls:
      insecure: true
  otlp/newrelic:
    endpoint: otlp.nr-data.net:443
    headers:
      api-key: <new_relic_ingest_licence_key> # Send traces to New Relic.
      # Read more in docs here: https://docs.newrelic.com/docs/more-integrations/open-source-telemetry-integrations/opentelemetry/opentelemetry-setup/#collector-export
      # And here: https://docs.newrelic.com/docs/more-integrations/open-source-telemetry-integrations/opentelemetry/collector/opentelemetry-collector-basic/

service:
  pipelines:
    traces/tracetest:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/tracetest]
    traces/newrelic:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, otlp/newrelic]
```

## Run the OpenTelemetry Demo with Tracetest and New Relic

To start both the OpenTelemetry Demo and Tracetest we will run this command:

```bash
docker-compose up
```

This will start your Tracetest instance on `http://localhost:11633/`. Go ahead and open it up.

[Start creating tests in the Web UI](https://docs.tracetest.io/web-ui/creating-tests)! Make sure to use the endpoints within your Docker network like `http://otel-frontend:8080/` when creating tests.

This is because your OpenTelemetry Demo and Tracetest are in the same network.

> Note: View the `demo` section in the `tracetest.config.yaml` to see which endpoints from the OpenTelemetry Demo are available for running tests.

Here's a sample of a failed test run, which happens if you add this assertion:

```
attr:tracetest.span.duration  < 50ms
```

![](https://res.cloudinary.com/djwdcmwdz/image/upload/v1672998179/Blogposts/tracetest-lightstep-partnership/screely-1672998159326_depw45.png)

Increasing the duration to a more reasonable `500ms` will make the test pass.

![](https://res.cloudinary.com/djwdcmwdz/image/upload/v1672998252/Blogposts/tracetest-lightstep-partnership/screely-1672998249450_mngghb.png)

## Run Tracetest tests with the Tracetest CLI

First, [install the CLI](https://docs.tracetest.io/getting-started/installation#install-the-tracetest-cli).
Then, configure the CLI:

```bash
tracetest configure --endpoint http://localhost:11633
```

Once configure, you can run a test against the Tracetest instance via the terminal.

Check out the `./tracetest/e2e/http-test.yaml` file.

```yaml
# http-test.yaml

type: Test
spec:
  id: YJmFC7hVg
  name: Otel - List Products
  description: Otel - List Products
  trigger:
    type: http
    httpRequest:
      url: http://otel-frontend:8084/api/products
      method: GET
      headers:
        - key: Content-Type
          value: application/json
  specs:
    - selector:
        span[tracetest.span.type="http" name="API HTTP GET" http.target="/api/products"
        http.method="GET"]
      assertions:
        - attr:http.status_code   =   200
        - attr:tracetest.span.duration  <  50ms
    - selector: span[tracetest.span.type="rpc" name="grpc.hipstershop.ProductCatalogService/ListProducts"]
      assertions:
        - attr:rpc.grpc.status_code = 0
    - selector:
        span[tracetest.span.type="rpc" name="hipstershop.ProductCatalogService/ListProducts"
        rpc.system="grpc" rpc.method="ListProducts" rpc.service="hipstershop.ProductCatalogService"]
      assertions:
        - attr:rpc.grpc.status_code = 0
```

This file defines the a test the same way you would through the Web UI.

To run the test, run this command in the terminal:

```bash
tracetest run test -f ./tracetest/e2e/http-test.yaml
```

This test will fail just like the sample above due to the `attr:tracetest.span.duration  <  50ms` assertion.

```bash
✘ Otel - List Products (http://localhost:11633/test/YJmFC7hVg/run/9/test)
	✘ span[tracetest.span.type="http" name="API HTTP GET" http.target="/api/products" http.method="GET"]
		✘ #cb68ccf586956db7
			✔ attr:http.status_code   =   200 (200)
			✘ attr:tracetest.span.duration  <  50ms (72ms) (http://localhost:11633/test/YJmFC7hVg/run/9/test?selectedAssertion=0&selectedSpan=cb68ccf586956db7)
	✔ span[tracetest.span.type="rpc" name="grpc.hipstershop.ProductCatalogService/ListProducts"]
		✔ #634f965d1b34c1fd
			✔ attr:rpc.grpc.status_code = 0 (0)
	✔ span[tracetest.span.type="rpc" name="hipstershop.ProductCatalogService/ListProducts" rpc.system="grpc" rpc.method="ListProducts" rpc.service="hipstershop.ProductCatalogService"]
		✔ #33a58e95448d8b22
			✔ attr:rpc.grpc.status_code = 0 (0)
```

If you edit the duration as in the Web UI example above, the test will pass!

## View trace spans over time in New Relic

To access a historical overview of all the trace spans the OpenTelemetry Demo generates, jump over to your New Relic account.

![](https://res.cloudinary.com/djwdcmwdz/image/upload/v1673009546/Blogposts/tracetest-new-relic-partnerships/screely-1673009541979_glib3l.png)

You can also drill down into a partucular trace as well.

![](https://res.cloudinary.com/djwdcmwdz/image/upload/v1673010042/Blogposts/tracetest-new-relic-partnerships/screely-1673010038074_aodp7d.png)

With New Relic and Tracetest, you get the best of both worlds. You can run trace-based tests and automate running E2E and integration tests against real trace data. And, use New Relic to get a historical overview of all traces your distributed application generates.

## Learn more

Feel free to check out our [docs](https://docs.tracetest.io/), and join our [Discord Community](https://discord.gg/8MtcMrQNbX) for more info!
