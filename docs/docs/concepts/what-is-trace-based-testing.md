# What is Trace-Based Testing


Trace-Based Testing is a means of conducting deep integration or system tests by utilizing the rich data contained in a distributed system trace.


## **What is a Distributed Trace?**

A Distributed Trace, more commonly known as a Trace, records the paths taken by requests (made by an application or end-user) take as they propagate through multi-service architectures, like microservice and serverless applications. [Source - OpenTelemetry.io](https://opentelemetry.io/docs/concepts/observability-primer/)

In Tracetest, after selecting a test from the first screen and clicking on the **Trace** tab, you will see the distributed trace for the selected test:

![Trace Example](../img/trace-example.png)


<!---![Trace & Spans Diagram](../img/trace-explainer.gif)-->

### **What is a Span?**

Traces are comprised of spans. A span represents a single operation in a trace. Spans are nested, typically with a parent child relationship to form a deeply nested tree.

![Span Example](../img/span-example.png)

### **What Data do Spans Contain?**


A span contains the data about the operation it represents. This data includes:

- The span name.

- Start and end timestamp.

- List of events (if instrumented).

- Attributes

### **What are Attributes?**

Attributes are a key-value pair, and they contain information about the operation. A developer can manually add additional attributes to a span, enriching the data. There are [Semantic Conventions](https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/) that provide recommended names for the attributes for common types of calls such as database, http, messaging, etc.

## **What is a Test Spec?**


In Tracetest, a Test Spec is comprised of two parts:


- Selectors
- Checks

<!--- ![Selectors and Checks](img/assertion-explainer.gif) -->

### **What is a Selector?**


A selector contains criteria to limit the scope of the spans from a trace that we wish to assert against. A selector can be very narrow, only selecting on one span, or very wide, selecting all spans or all spans of a certain type or other characteristics. Underlying this capability is a [selector language](./selectors).


### **What is a Check?**


A check is a logical verification that will be performed on all spans that match the selector. It is comprised of an attribute, a comparison operator and a value.

### **What is a Span Signature?**


A span signature is an automatically computed selector that has enough elements to specify a single span. It uses a combination of attributes in the selected span to automatically build the selector. If a trace has multiple spans that are almost identical, the span signature may still match more than one span. You can alter the selector in this case to be more specific by adding other attributes or specifying an ancestor span.
