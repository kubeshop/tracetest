# Introduction to Trace Based Testing

Trace Based Testing is a means of conducting deep integration or system tests by utilizing the rich data contained in a distributed system trace.

## What is a Distributed Trace?

A Distributed Trace, more commonly known as a Trace, records the paths taken by requests (made by an application or end-user) take as they propagate through multi-service architectures, like microservice and serverless applications. [Source - OpenTelemetry.io](https://opentelemetry.io/docs/concepts/observability-primer/)

## What is a Span?

Traces are comprised of spans. A span represents a single operation in a trace. Spans are nested, typically with a parent child relationship to form a deeply nested tree.

## What data do Spans contain?

A span contains the data about the operation it represents. This data includes:

- The span name
- Start and end timestamp
- List of events (if instrumented)
- Attributes

## What are Attributes?

Attributes are a key-value pair, and they contain information about the operation. A developer can manually add additional attributes to a span, enriching the data. There are [Semantic Conventions](https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/) that provide recommended names for the attributes for common types of calls such as database, http, messaging, etc.

## What is an assertion?

In Tracetest an assertion is comprised of two parts:

- Selector / filters
- Checks

## What is a selector / filter?

A selector / filter contains criteria to limit the scope of the spans from a trace that we wish to assert against. A selector / filter can be very narrow, only selecting on span, or very wide, selecting all spans or all spans of a certain type or other characteristic.

## What is a check?

A check is a logical verification that will be performed on all spans that match the selector / filter. It is comprised of an attribute, a comparison operator, and a value.

## What is a span signature?

A span signature is an automatically computed selector / filter that has enough elements to specify a single span. It uses a combination of attributes in the selected span to automatically build the filter. If a trace has multiple spans that are almost identical, the span signature may still match more than one span. You can alter the selector / filter in this case to be more specific by adding other attributes or specifying an ancestor span.
