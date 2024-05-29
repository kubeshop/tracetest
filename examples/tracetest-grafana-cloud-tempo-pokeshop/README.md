# Tracetest + Grafana Cloud Tempo + Pokeshop

> [Read the detailed recipe for setting up Grafana Cloud Tempo with Tractest in our blog.](https://tracetest.io/blog/monitoring-and-testing-cloud-native-apis-with-grafana)

This examples' objective is to show how you can configure your Tracetest account to connect to Grafana Cloud Tempo and use it as a trace data store.

Feel free to check out the [docs](https://docs.tracetest.io/), and join our [Slack Community](https://dub.sh/tracetest-community) for more info!

## Steps

1. [Install the tracetest CLI](https://docs.tracetest.io/installing/)
2. Run `tracetest configure -t <TRACETEST_API_TOKEN>`
3. Run the project by using docker-compose: `docker-compose up -d` (Linux) or `docker compose up -d` (Mac)
4. Set your Grafana Cloud tokens by base64 encoding your `username:password` using headers like this: `authorization: Basic <base64 encoded username:password>`.
5. Test if it works by running: `tracetest run test -f tests/test.yaml`. This would trigger a test that will send and retrieve spans from Tempo in your Grafana Cloud account.
6. View traces in your Grafana Cloud. Use this TraceQL query:

    ```text
    { name="Tracetest trigger" } || { name="POST /pokemon/import?" } || { name="POST /pokemon/import" }
    ```
