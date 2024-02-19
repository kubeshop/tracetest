# Tracetest Cloud + Artillery.io

This repo shows how to integrate tracetest tests with artillery.io tests


## Steps

1. Get your tracetest cloud api key and replace it on `atillery-test.yaml`
2. Install the `artillery-plugin-tracetest` // TODO: fix this
3. Run the project by using docker-compose: `docker-compose up -d` (Linux) or `docker compose up -d` (Mac)
4. Start the agent:
  > Since we need to use the otel collector from the docker compose, you need to run the agent
  > with changed otlp ports
  ```
  TRACETEST_OTLP_SERVER_GRPC_PORT=9999 TRACETEST_OTLP_SERVER_HTTP_PORT=9998 tracetest start
  ```
5. Run the tests `artillery run atillery-test.yaml`
