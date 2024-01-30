# Observability Stack Example

This folder has minimal code to run a local observability stack, plug an API into it, and do some trace-based tests.

To run the observability stack and the local API, execute:

```sh
# run the observability stack 
docker compose up -d

# install dependencies and run API
npm install
npm run with-telemetry
```

To run it with trace-based tests, execute:

```sh
# run the observability stack 
export TRACETEST_API_KEY="{{Your Tracetest.io Token}}"
docker compose -f ./docker-compose.yaml -f docker-compose.tracetest.yaml up -d

# install dependencies and run API
npm install
npm run with-telemetry
```

And then run a test with:
```sh
tracetest run test -f ./test-api.yaml
```
