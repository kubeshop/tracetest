# Context Propagation between Rails API instrumented with ddtrace and OTel SDK

This repository shows an example of how to propagate context between `ddtrace` instrumented APIs and OTel SDK APIs.

## The problem

Datadog represents its TraceIds as an uint64, while OTel SDK represents it as a 128-bit hexadecimal string, causing some issues with Trace propagation when you use OTel Collector's `datadog` receiver along an external Tracing Backend that isn't Datadog.

## Solution

Datadog can build an entire TraceID using a tag called `_dd.p.tid` as the upper part of a TraceId and its own TraceId in hexadecimal format as the lower part. To use this data to generate the OTel TraceID we changed how our `ddtrace` sends Traces to send the `_dd.p.tid` tag in all spans (normally it sends this attribute only for the first trace span) and changed our OTel Collector to have a `transform` processor, that will grab the upper and lower part of the TraceID and will reconstruct it.

### How to test it

There are 4 `docker-compose` files, showing each stage of the solution. Each file has two Rails APIs defined, one Datadog-instrumented and another OTelSDK-instrumented, one OTel Collector configured to receive traces from OTLP and Datadog and a Jaeger to ingest the data and show it.

After starting the APIs with one of the docker compose files, you can test the APIs locally by executing one of the following `curl` commands and checking on Jaeger on `http://localhost:16686/`:
```sh
# Test API call from OTelSDK-instrumented API to Datadog-instrumented API
curl http://localhost:8080/remotehello

# Test API call from OTelSDK-instrumented API to Datadog-instrumented API
curl http://localhost:8081/remotehello
```

The first file, `docker-compose.step1.yaml`, shows the propagation problem happening. When you call any of the `curl` requests, you should see two split traces on Jaeger, one for each API. You can run it with:
```sh
docker compose -f ./docker-compose.step1.yaml up
```

The second file, `docker-compose.step2.yaml`, runs the OTel Collector with a `transform` processor that fixes the `TraceID` for spans that have `_dd.p.tid` attribute. After calling a `curl` request, you should see one trace that has the first span of the Datadog-instrumented API call along with all spans emitted by the OTelSDK-instrumented API, and another trace with Datadog-instrumented API that doesn't have `_dd.p.tid` attribute. You can run it with:
```sh
docker compose -f ./docker-compose.step2.yaml up
```

The third file, `docker-compose.step3.yaml`, does a patch on `ddtrace` to inject `_dd.p.tid` attribute on each span instrumented by `ddtrace`, also renaming `_dd.p.tid` attribute to `propagation.upper_trace_id`. This version runs the OTel Collector with a `transform` processor that fixes the `TraceID` for spans that have `propagation.upper_trace_id` attribute. After calling a `curl` request, you should see just one trace with all spans emitted by the OTelSDK-instrumented API and Datadog-instrumented API. You can run it with:
```sh
docker compose -f ./docker-compose.step3.yaml up
```

The last file, `docker-compose.step4.yaml`, adds Tracetest to the stack to allow us to test the traces and guarantee that everything is working right. You can run it with:
```sh
TRACETEST_API_KEY=your-api-key docker compose -f ./docker-compose.step4.yaml up
```
