# Tracetest Run Groups with CLI

> [Read the detailed recipe for using Run Groups in Tractest in our documentation.](https://docs.tracetest.io/concepts/run-groups)

This is a simple quick start on how to run 3 tests in a Run Group.

## Run Tracetest with Docker Compose

```bash
docker compose up --build
```

## Run 3 tests in a Group

```bash
tracetest run test -f ./test-api-1.yaml -f test-api-2.yaml -f test-api-3.yaml --group nodejs-group-X
```

> Edit the `X` every time you run a group of tests.

---

Feel free to check out the [docs](https://docs.tracetest.io/), and join our [Slack Community](https://dub.sh/tracetest-community) for more info!
