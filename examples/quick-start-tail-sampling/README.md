# Quick Start - Trace-based tests with Tail Sampling configuration

> [Read the detailed recipe for setting up OpenTelemetry Collector with Tractest in our documentation.](https://docs.tracetest.io/examples-tutorials/recipes/running-tracetest-without-a-trace-data-store)

This is a simple quick start example on how to set up `tail_sampling` into OTel Collector, allowing Tracetest to run tests in environments where we have a probabilistic sampler enabled and a percentage of the traces are not sent to the final data store.

If you want to run this example, just execute `docker compose up` on this folder.

To execute a Trace-based test with Tracetest against this structure just run `tracetest run test -f test-api-working.yaml`.

Feel free to check out the [docs](https://docs.tracetest.io/), and join our [Discord Community](https://discord.gg/8MtcMrQNbX) for more info!
