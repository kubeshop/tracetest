# Quick Start

This page showcases getting started with Tracetest by using Docker or the Tracetest CLI.

More detailed guides can be found [here for Docker](./docker), and [here for the Tracetest CLI](./cli).

## Prerequisites

:::info
You need to add [OpenTelemetry instrumentation](https://opentelemetry.io/docs/instrumentation/) to your code and configure sending traces to a trace data store, or Tracetest directly, to benefit for Tracetest's trace-based testing.
:::

:::caution
Postgres is a prerequisite for Tracetest to work. It stores Tracetest's trace data. Make sure to have a Postgres service.
:::

## Docker

Download [this sample setup](https://github.com/kubeshop/tracetest/tree/main/examples/collector) from our GitHub examples.

Start Docker Compose.

```bash
docker compose up
```

```bash title="Condensed expected output from the Tracetest container:"
Starting tracetest ...
...
2022/11/28 18:24:09 HTTP Server started
...
```

Open your browser on [`http://localhost:11633`](http://localhost:11633).

Create a [test](../web-ui/creating-tests.md).

Read the detailed setup on the [Docker installation page](./docker).

:::info
Running a test against `localhost` will resolve as the 127.0.0.1 inside the Tracetest container. To run tests against apps running on your local machine, add them to the same network and use service name mapping instead. Example: Instead of running an app on `localhost:8080`, add it to your Docker Compose file, connect it to the same network as your Tracetest service, and use `service-name:8080` in the URL field when creating an app.

You can reach services running on your local machine using:

- Linux (docker version < 20.10.0): `172.17.0.1:8080`
- MacOS (docker version >= 18.03) and Linux (docker version >= 20.10.0): `host.docker.internal:8080`
:::

## CLI

Install the Tracetest CLI:

```bash
curl -L https://raw.githubusercontent.com/kubeshop/tracetest/main/install-cli.sh | bash
```

Install the Tracetest server:

```bash
tracetest server install
```

:::note
Follow the prompts and continue with all the default settings.
This will generate an empty `docker-compose.yaml` file and a `./tracetest/` directory that contains another `docker-compose.yaml` and two more config files. One for Tracetest and one for OpenTelemetry collector.
:::

Start Docker Compose from the directory where you ran `tracetest server install`.

```bash
docker compose -f docker-compose.yaml -f tracetest/docker-compose.yaml up -d
```

```bash title="Condensed expected output from the Tracetest container:"
Starting tracetest ...
...
2022/11/28 18:24:09 HTTP Server started
...
```

Open your browser on [`http://localhost:11633`](http://localhost:11633).

Create a [test](../web-ui/creating-tests.md).

Read the detailed setup on the [CLI installation page](./cli).

:::info
Running a test against `localhost` will resolve as the 127.0.0.1 inside the Tracetest container. To run tests against apps running on your local machine, add them to the same network and use service name mapping instead. Example: Instead of running an app on `localhost:8080`, add it to your Docker Compose file, connect it to the same network as your Tracetest service, and use `service-name:8080` in the URL field when creating an app.

You can reach services running on your local machine using:

- Linux (docker version < 20.10.0): `172.17.0.1:8080`
- MacOS (docker version >= 18.03) and Linux (docker version >= 20.10.0): `host.docker.internal:8080`
:::
