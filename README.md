# projectx


## API

### Create Test
```
curl -X POST http://localhost:8080/tests \
   -H "Content-Type: application/json" \
   -d '{"name":"first-test","serviceUnderTest":{"url":"http://localhost:3030/hello-instrumented"}}'
```

Response:

```
{"id":"3ba6b38a-99b5-4bb5-83eb-e8fa78f377fb","name":"first-test","serviceUnderTest":{"url":"http://localhost:3030/hello-instrumented"}}
```

### Run Test


curl -X POST http://localhost:8080/tests/
```
curl -X POST http://localhost:8080/tests/3ba6b38a-99b5-4bb5-83eb-e8fa78f377fb/run
```

### See the trace in jaeger

```
http://localhost:8081/trace/0194fdc2fa2ffcc041d3ff12045b73c8
```
