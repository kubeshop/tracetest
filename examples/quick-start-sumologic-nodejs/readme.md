# Quick Start - Node.js app with Sumo Logic and Tracetest

> [Read the detailed recipe for setting up Sumo Logic with Tractest in our documentation.](https://docs.tracetest.io/examples-tutorials/recipes/running-tracetest-with-sumologic)

This is a simple quick start on how to configure a Node.js app to use OpenTelemetry instrumentation with traces and Tracetest for enhancing your e2e and integration tests with trace-based testing. The infrastructure will use Sumo Logic as the trace data store, and the Sumo Logic distribution of the OpenTelemetry Collector to receive traces from the Node.js app and send them to Sumo Logic.

## Steps to run Tracetest

### 1. Configure Tracetest Agent and Sumo Logic

- Sign up on [`app.tracetest.io`](https://app.tracetest.io).
- Create a new environment.
- Copy the start command under `Settings > Agent`.

Configure Sumo Logic as a Tracing Backend:

- Configure the Tracetest CLI by creating an environment token under `Settings > Tokens`.

```bash
tracetest configure -t <YOUR_API_TOKEN>
```

- Add your Sumo Logic AccessID and AccessKey to `tracetest.datastore.yaml`
- Apply the tracing backend configuration

```bash
tracetest apply datastore -f ./tracetest.datastore.yaml
```

- Configure the OpenTelemetry Collector (`collector.config.yaml`) with a Sumo Logic installation token.

> Note: Here's a guide which Sumo Logic API endpoint to use: https://help.sumologic.com/docs/api/getting-started/#which-endpoint-should-i-should-use

### 2. Start Node.js App

Set the env vars in the `.env` file.

Run the example with Docker.

```bash
docker compose up -d
```

### 3. Run tests

- Create and run a test against `http://localhost:8080` on [`https://app.tracetest.io/`](https://app.tracetest.io/).
- View the `./test-api.yaml` for reference.

```bash
tracetest run test -f ./test-api.yaml
```

## Steps to run Tracetest Core

### 1. Start Node.js App and Tracetest Core with Docker Compose

```bash
docker compose -f ./docker-compose.yaml -f ./tracetest/docker-compose.yaml up --build
```

### 2. Run tests

Once started, you will need to make sure to trigger tests with correct service names since both the Node.js app and Tracetest Core are running in the same Docker Network. In this example the Node.js app would be at `http://app:8080`. View the `./test-api.yaml` for reference.

---

Feel free to check out the [docs](https://docs.tracetest.io/), and join our [Slack Community](https://dub.sh/tracetest-community) for more info!
