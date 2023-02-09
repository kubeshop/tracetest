# Running Tracetest With Lightstep

:::note
[Check out the source code on GitHub here.](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-lightstep)
:::

## OpenTelemetry Demo `v0.3.4-alpha` with Lightstep, OpenTelemetry and Tracetest

This is a simple sample app on how to configure the [OpenTelemetry Demo `v0.3.4-alpha`](https://github.com/open-telemetry/opentelemetry-demo) to use [Tracetest](https://tracetest.io/) for enhancing your E2E and integration tests with trace-based testing, and [Lightstep](https://lightstep.com/) as a trace data store.

## Prerequisites

You will need [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) installed on your machine to run this sample app! Additionally, you will need a Lightstep account and access token. Sign up to Lightstep [here](https://app.lightstep.com/signup/developer).

## Project structure

The project is built with Docker Compose. It contains two distinct `docker-compose.yaml` files.

### 1. OpenTelemetry Demo
The `docker-compose.yaml` file and `.env` file in the root directory are for the OpenTelemetry Demo.

### 2. Tracetest
The `docker-compose.yaml` file, `collector.config.yaml`, and `tracetest.config.yaml` in the `tracetest` directory are for the setting up Tracetest and the OpenTelemetry Collector.

The `tracetest` directory is self-contained and will run all the prerequisites for enabling OpenTelemetry traces and trace-based testing with Tracetest, as well as routing all traces the OpenTelemetry Demo generates to Lightstep.

### Docker Compose Network
All `services` in the `docker-compose.yaml` are on the same network and will be reachable by hostname from within other services. E.g. `tracetest:21321` in the `collector.config.yaml` will map to the `tracetest` service, where the port `21321` is the port where Tracetest accepts traces.

## OpenTelemetry Demo

The [OpenDelemetry Demo](https://github.com/open-telemetry/opentelemetry-demo) is a sample microservice-based app with the purpose to demo how to correctly set up OpenTelemetry distributed tracing.

Read more about the OpenTelemetry Demo [here](https://opentelemetry.io/blog/2022/announcing-opentelemetry-demo-release/).

The `docker-compose.yaml` contains 12 services.

To start the OpenTelemetry Demo by itself, run this command:

```bash
docker compose build # optional if you haven't already built the images
docker compose up
```

This will start the OpenTelemetry Demo. Open up `http://localhost:8084` to make sure it's working. But, you're not sending the traces anywhere.

Let's fix this by configuring Tracetest and OpenTelemetry Collector to forward trace data to both Lightstep and Tracetest.

## Tracetest

The `docker-compose.yaml` in the `tracetest` directory is configured with three services.

- **Postgres** - Postgres is a prerequisite for Tracetest to work. It stores trace data when running the trace-based tests.
- [**OpenTelemetry Collector**](https://opentelemetry.io/docs/collector/) - A vendor-agnostic implementation of how to receive, process and export telemetry data.
- [**Tracetest**](https://tracetest.io/) - Trace-based testing that generates end-to-end tests automatically from traces.

```yaml
version: "3.2"
services:
  tracetest:
    restart: unless-stopped
    image: kubeshop/tracetest:${TAG:-latest}
    platform: linux/amd64
    ports:
      - 11633:11633
    extra_hosts:
      - "host.docker.internal:host-gateway"
    volumes:
      - type: bind
        source: ./tracetest/tracetest.config.yaml
        target: /app/config.yaml
    healthcheck:
      test: ["CMD", "wget", "--spider", "localhost:11633"]
      interval: 1s
      timeout: 3s
      retries: 60
    depends_on:
      postgres:
        condition: service_healthy
      otel-collector:
        condition: service_started

  postgres:
    image: postgres
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_USER: postgres
    healthcheck:
      test: pg_isready -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"
      interval: 1s
      timeout: 5s
      retries: 60

  otel-collector:
    image: otel/opentelemetry-collector-contrib:0.68.0
    restart: unless-stopped
    command:
      - "--config"
      - "/otel-local-config.yaml"
    volumes:
      - ./tracetest/collector.config.yaml:/otel-local-config.yaml
```

Tracetest depends on both Postgres and the OpenTelemetry Collector. Both Tracetest and the OpenTelemetry Collector require config files to be loaded via a volume. The volumes are mapped from the root directory into the `tracetest` directory and the respective config files.

**Why?** To start both the OpenTelemetry Demo and Tracetest we will run this command:

```bash
docker-compose -f docker-compose.yaml -f tracetest/docker-compose.yaml up # add --build if the images are not built already
```

The `tracetest.config.yaml` file contains the basic setup of connecting Tracetest to the Postgres instance, and defining the trace data store and exporter. The data store is set to OTLP meaning the traces will be stored in Tracetest itself. The exporter is set to the OpenTelemetry Collector.

```yaml
# tracetest.config.yaml

---
postgresConnString: "host=postgres user=postgres password=postgres port=5432 sslmode=disable"

poolingConfig:
  maxWaitTimeForTrace: 30s
  retryDelay: 500ms

# This will populate sample tests in the Tracetest Web UI you can run to try out Tracetest.
demo:
  enabled: [otel]
  endpoints:
    otelFrontend: http://otel-frontend:8084
    otelProductCatalog: otel-productcatalogservice:3550
    otelCart: otel-cartservice:7070
    otelCheckout: otel-checkoutservice:5050

experimentalFeatures: []

googleAnalytics:
  enabled: false

telemetry:
  dataStores:
    otlp:
      type: otlp

  exporters:
    collector:
      serviceName: tracetest
      sampling: 100
      exporter:
        type: collector
        collector:
          endpoint: otel-collector:4317

server:
  telemetry:
    dataStore: otlp
    exporter: collector
    applicationExporter: collector
```

**How to send traces to Tracetest and Lightstep?**

The `collector.config.yaml` explains that. It receives traces via either `grpc` or `http`. Then, exports them to Tracetest's OTLP endpoint `tracetest:21321` in one pipeline, and to Lightstep in another.

Make sure to add your Lightstep access token in the headers of the `otlp/ls` exporter.

```yaml
receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:
    timeout: 100ms

exporters:
  logging:
    logLevel: debug
  # OTLP for Tracetest
  otlp/tt:
    endpoint: tracetest:21321 # Send traces to Tracetest. Read more in docs here: https://docs.tracetest.io/configuration/connecting-to-data-stores/opentelemetry-collector
    tls:
      insecure: true
  # OTLP for Lightstep
    endpoint: ingest.lightstep.com:443
    headers: 
      "lightstep-access-token": "<your-lightstep-access-token>" # Send traces to Lightstep. Read more in docs here: https://docs.lightstep.com/otel/otel-quick-start 

service:
  pipelines:
    traces/tt:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/tt]
    traces/ls:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, otlp/ls]
```

**Important!** Take a closer look at the sampling configs in both the `collector.config.yaml` and `tracetest.config.yaml`. They both set sampling to 100%. This is crucial when running trace-based e2e and integration tests.

## Run both the OpenTelemetry Demo app and Tracetest

To start both the OpenTelemetry and Tracetest we will run this command:

```bash
docker-compose -f docker-compose.yaml -f tracetest/docker-compose.yaml up # add --build if the images are not built already
```

This will start your Tracetest instance on `http://localhost:11633/`. Go ahead and open it up.

[Start creating tests in the Web UI](https://docs.tracetest.io/web-ui/creating-tests)! Make sure to use the endpoints within your Docker network like `http://otel-frontend:8084/` when creating tests.

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
tracetest configure --endpoint http://localhost:11633 --analytics
```

Once configure, you can run a test against the Tracetest instance via the terminal.

Check out the `http-test.yaml` file.

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
  - selector: span[tracetest.span.type="http" name="API HTTP GET" http.target="/api/products"
      http.method="GET"]
    assertions:
    - attr:http.status_code   =   200
    - attr:tracetest.span.duration  <  50ms
  - selector: span[tracetest.span.type="rpc" name="grpc.hipstershop.ProductCatalogService/ListProducts"]
    assertions:
    - attr:rpc.grpc.status_code = 0
  - selector: span[tracetest.span.type="rpc" name="hipstershop.ProductCatalogService/ListProducts"
      rpc.system="grpc" rpc.method="ListProducts" rpc.service="hipstershop.ProductCatalogService"]
    assertions:
    - attr:rpc.grpc.status_code = 0
```

This file defines the a test the same way you would through the Web UI.

To run the test, run this command in the terminal:

```bash
tracetest test run -d ./http-test.yaml -w
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

## View trace spans over time in Lightstep

To access a historical overview of all the trace spans the OpenTelemetry Demo generates, jump over to your Lightstep account.

![](https://res.cloudinary.com/djwdcmwdz/image/upload/v1672998664/Blogposts/tracetest-lightstep-partnership/screely-1672998658856_lae7ml.png)

You can also drill down into a particular trace as well.

![](https://res.cloudinary.com/djwdcmwdz/image/upload/v1672998974/Blogposts/tracetest-lightstep-partnership/screely-1672998969770_iwmjy5.png)

With Lightstep and Tracetest, you get the best of both worlds. You can run trace-based tests and automate running E2E and integration tests against real trace data. And, use Lightstep to get a historical overview of all traces your distributed application generates.

## Learn more

Feel free to check out our [examples in GitHub](https://github.com/kubeshop/tracetest/tree/main/examples), and join our [Discord Community](https://discord.gg/8MtcMrQNbX) for more info!
