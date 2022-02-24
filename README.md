# projectx


## API

### Create Test

```
curl -X POST http://localhost:8080/api/tests \
   -H "Content-Type: application/json" \
   -d '{"name":"first-test","serviceUnderTest":{"url":"http://http-app:3030/hello-instrumented"}}'

```

Response:

```
{"id":"3ba6b38a-99b5-4bb5-83eb-e8fa78f377fb","name":"first-test","serviceUnderTest":{"url":"http://localhost:3030/hello-instrumented"}}
```

### Create Assertions


### Run Test

```
curl -X POST http://localhost:8080/api/test/3ba6b38a-99b5-4bb5-83eb-e8fa78f377fb/assertions \
   -H "Content-Type: application/json" \
   -d '{"selector":"resourceSpans[?resource.attributes[?key=='service.name' && value.stringValue=='productcatalog']]","comparable":"Comparable","operator":"Operator"}'
```

```
curl -X POST http://localhost:8080/api/tests/3ba6b38a-99b5-4bb5-83eb-e8fa78f377fb/run
```
Response:

```
{"id":"833aa474-a8db-494f-ba13-9a7fa92ea03d"}
```

### Get Results

```
curl -X GET http://localhost:8080/api/tests/fbd008c4-ec08-4312-9bb3-5c70ae55b255/results/833aa474-a8db-494f-ba13-9a7fa92ea03d
```

### Get Trace

```
curl -X GET http://localhost:8080/api/tests/fbd008c4-ec08-4312-9bb3-5c70ae55b255/results/833aa474-a8db-494f-ba13-9a7fa92ea03d/trace
```
