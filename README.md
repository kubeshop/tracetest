# Tracetest

## Overview

Testing and debugging software built on Micro-Services architectures is hard.

As many as 30 to 100 services may be involved in a single flow. Written in multiple languages. With several backend data stores, message busses, and technologies. Understanding the flow is hard - having enough experience & wide ranging knowledge to create tests to verify it is working properly is even harder.

Tracetest makes this easy. Pick an api to test. Tracetest uses your tracing infrastructure to trace this api call. This trace is the blueprint of your entire system, showing all the activity. Use this blueprint to graphically define assertions on different services throughout the trace, checking return statuses, data, or even execution times of systems.

Examples:
assert that all database calls return in less than 250 ms
assert that one particular micro service returns a 200 code when called
assert that a Kafka queue successful delivers a payload to a dependent micro service.

Once the test is built, it can be run automatically as part of a build process or manually. Every test has a trace attached, allowing you to immediately see what worked, and what did not, reducing the need to reproduce the problem to see the underlying issue.

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
