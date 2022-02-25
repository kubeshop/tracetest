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


```
curl -X POST http://localhost:8080/tests/3ba6b38a-99b5-4bb5-83eb-e8fa78f377fb/run
```

```

```

### Get Results

```
curl -X GET http://localhost:8080/tests/fbd008c4-ec08-4312-9bb3-5c70ae55b255/results/833aa474-a8db-494f-ba13-9a7fa92ea03d
```

### See the trace in jaeger

```
http://localhost:8081/trace/0194fdc2fa2ffcc041d3ff12045b73c8
```
