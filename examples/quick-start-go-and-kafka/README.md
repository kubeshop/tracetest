# Quick Start - Go API with Kafka messaging

> [Read the detailed recipe for setting up OpenTelemetry Collector with Tractest in our documentation.](https://docs.tracetest.io/examples-tutorials/recipes/running-tracetest-without-a-trace-data-store)

This is a simple quick start on how to configure two Go lang apps that interacts with [Apache Kafka](https://kafka.apache.org/): a `producer-api` and a `consumer-api` , and how to test if the messaging is properly working with a trace-based test.

If you want to run this example, just execute `docker compose up` on this folder.

To execute a Trace-based test with Tracetest against this structure just run `tracetest run test -f test.yaml`.

Feel free to check out the [docs](https://docs.tracetest.io/), and join our [Discord Community](https://discord.gg/8MtcMrQNbX) for more info!
