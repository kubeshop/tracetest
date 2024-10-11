# End-to-End Observability

This project demonstrates an end-to-end observability setup using the Grafana stack: **Grafana**, **Prometheus**, **Loki**, and **Tempo**. The application uses **OpenTelemetry (OTel)** to generate traces and metrics, and **Winston-Loki** to produce logs.

## Stack Overview

- **Grafana**: Visualization of metrics, traces, and logs.
- **Prometheus**: Metrics collection and storage.
- **Loki**: Logs aggregation.
- **Tempo**: Distributed tracing backend.
- **OpenTelemetry (OTel)**: Traces and metrics generator.
- **Winston-Loki**: Logging integration for Loki.

## Prerequisites

Make sure you have the following installed:

- [Docker](https://docs.docker.com/get-docker/)
- [Node.js](https://nodejs.org/)

## Getting Started

### 1. Copy the `.env.sample` into a `.env`

- Get your Tracetest token and env id [here](https://app.tracetest.io/retrieve-token).

### 2. Start the Observability Stack and the App

Use Docker Compose:

```bash
docker-compose up -d
```

This will spin up the following services:

- Grafana on [localhost:3000](localhost:3000)
- Prometheus on [localhost:9090](localhost:9090)
- Loki on [localhost:3100](localhost:3100)
- Tempo on [localhost:3200](localhost:3200)
- Node app on [localhost:8081](localhost:8081)

### 3. Configure Trace Ingestion

```bash
tracetest apply datastore -f ./tracetest-trace-ingestion.yaml
```

### 3. Trigger a Test

```bash
tracetest run test -f ./tracetest-test.yaml
```
