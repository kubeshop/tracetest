# Tracetest + Grafana Tempo + Pokeshop

<!-- > [Read the detailed recipe for setting up Tracetest + Grafana Tempo + Pokeshop in our documentation.]() -->

This examples' objective is to show how you can:

1. Configure your Tracetest instance to connect to Grafana Tempo and use it as a trace data store.
2. Configure Grafana to query traces from Tempo

## Steps

1. [Install the tracetest CLI](https://docs.tracetest.io/installing/)
2. Run `tracetest configure --endpoint http://localhost:11633` on a terminal
3. Run the project by using docker-compose: `docker-compose up -d` (Linux) or `docker compose up -d` (Mac)
4. Test if it works by running: `tracetest run test -f tests/test.yaml`. This would trigger a test that will send and retrieve spans from the Grafana Tempo instance that is running on your machine. View the test on `http://localhost:11633`.
5. View traces in Grafana on `http://localhost:3000`. Use this TraceQL query:

    ```yaml
    { name="POST /pokemon/import" }
    ```
