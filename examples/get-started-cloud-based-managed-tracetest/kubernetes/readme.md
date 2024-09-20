# Running Node.js with OpenTelemetry and Cloud-based Managed Tracetest in Kubernetes

## Build and Push the Docker Image to Docker Hub

1. Edit the `docker-compose.yaml` to use your Docker Hub username.
2. Build the image.

    ```bash
    docker compose build
    ```

3. Push the image.

    ```bash
    docker push YOUR_DOCKERHUB_USERNAME/get-started-cloud-based-managed-tracetest-app
    ```

## Run Kubernetes Manifests (both `app` and `tracetest-agent`)

1. Update the `app.yaml` to use the `YOUR_DOCKERHUB_USERNAME/get-started-cloud-based-managed-tracetest-app` image you pushed to DockerHub.
2. [Sign in to Tracetest](https://app.tracetest.io/).
3. [Create an Organization](https://docs.tracetest.io/concepts/organizations).
4. Retrieve your [Tracetest Organization API Key/Token and Environment ID](https://app.tracetest.io/retrieve-token).
5. Update the [Tracetest Agent](https://docs.tracetest.io/concepts/agent) env vars `<TRACETEST_API_KEY>` and `<TRACETEST_ENVIRONMENT_ID>` from step 2 in the `./manifests/tracetest-agent.yaml` file.

```bash
kubectl apply -f ./manifests
```

Tracetest Agent will run on gRPC and HTTP ports and use the Kubernetes DNS where you can access it via its service name.

- `http://tracetest-agent.default.svc.cluster.local:4317` — gRPC
- `http://tracetest-agent.default.svc.cluster.local:4318/v1/traces` — HTTP

## Configure Trace Ingestion for Localhost

Go to the Trace Ingestion tab in the settings, select OpenTelemetry, and enable the toggle.

## Run Trace-based Tests

You can now run tests against your apps on `http://app.default.svc.cluster.local:8080` by going to Tracetest and creating a new HTTP test.

1. Click create a test and select HTTP.
2. Add `http://app.default.svc.cluster.local:8080` in the URL field.
3. Click Run. You’ll see the response and trace data right away.

## (Optional) Use Helm to Deploy Tracetest Agent

Install Helm chart

```bash
helm repo add tracetestcloud https://kubeshop.github.io/tracetest-cloud-charts --force-update  && \

helm install demo tracetestcloud/tracetest-agent \
--set agent.apiKey=<TRACETEST_API_KEY> --set agent.environmentId=<TRACETEST_ENVIRONMENT_ID>
```

Tracetest Agent will run on gRPC and HTTP ports and use the Kubernetes DNS where you can access it via its service name. Where `demo` will be prefixed as the release name.

- `http://demo-tracetest-agent.default.svc.cluster.local:4317` — gRPC
- `http://demo-tracetest-agent.default.svc.cluster.local:4318/v1/traces` — HTTP

Edit trace exporter

```js
// Configure OTLP gRPC Trace Exporter
const traceExporter = new OTLPTraceExporter({
  // Default endpoint for OTLP gRPC is localhost:4317
  // You can change this to your OTLP collector or backend URL.
  url: 'http://demo-tracetest-agent.default.svc.cluster.local:4317',
});
```
