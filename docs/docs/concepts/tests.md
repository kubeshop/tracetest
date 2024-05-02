# Tests

Tracetest uses the concept of Tests to define how to trigger a test against your application, define assertions against its trace data, and automate its execution. Every time a Test is triggered it will create a Run.
- [Runs](../concepts/runs.md) (link)

![Tests](../img/tests.png)

A test allows you to:

- Execute a trigger, such as an HTTP request, a gRPC call, a TraceID, a Kafka queue, etc. to generate a trace.
- Fetch the resulting trace and analyze it.
- Add assertions against the trace data to verify the behavior of the system at every step of the request transaction.
- The assertions can check things like HTTP status codes, database call durations, gRPC return codes, and other aspects of the distributed system's behavior.
- Tests can be saved and run manually or as part of a CI/CD pipeline to ensure the quality and reliability of the distributed application.

A test in Tracetest is a way to define trace-based assertions that validate the end-to-end behavior of a distributed system, going beyond just UI or API-level testing.


