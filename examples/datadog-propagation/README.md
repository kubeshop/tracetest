# Context Propagation between Rails API instrumented with ddtrace and OTel SDK

This repository shows an example of how to propagate context between `ddtrace` instrumented APIs and OTel SDK APIs.

## The problem

Datadog represents its TraceIds as an uint64, while OTel SDK represents it as a 128-bit hexadecimal string, causing some issues with Trace propagation when you use OTel Collector's `datadog` receiver along an external Tracing Backend that isn't Datadog.
