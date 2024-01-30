# Observability Stack Example

This folder has the minimal code to run a local observability stack, plug an API into it, and do some trace based tests.

To run the Observability Stack and the local API, execute:

```sh
# run our Observability stack 
docker compose up -d

# install dependencies and run API
npm install
npm run with-telemetry

```

If you want to run it with Trace-based Tests, you can execute:

```sh
# run our Observability stack 
docker compose -f ./docker-compose.yaml -f docker-compose.tracetest.yaml up -d

# install dependencies and run API
npm install
npm run with-telemetry
```

And then run a test with:
```sh
tracetest run test -f ./test-api.yaml --server-url http://localhost:11633
```
