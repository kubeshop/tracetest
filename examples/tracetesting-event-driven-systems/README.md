# Tracetesting Event-driven systems

> [Read the detailed recipe for setting up OpenTelemetry Collector with Tractest in our documentation.](https://docs.tracetest.io/examples-tutorials/recipes/running-tracetest-without-a-trace-data-store)

This is a simple example from the article "Testing Event-driven Systems with OpenTelemetry" showing how to test an event-driven system using [Apache Kafka](https://kafka.apache.org/) as a event backbone.

If you want to run this example, just execute `docker compose up` on this folder.

Feel free to check out the [docs](https://docs.tracetest.io/), and join our [Discord Community](https://discord.gg/8MtcMrQNbX) for more info!

## Test Scenarios

We can test four different scenarios with this example:

### **Scenario 1**: Add order though Rest API

Here, we send a message as a regular user and check if all components are called as intended:

```sh
tracetest run test -f test-payment-order-submit-with-rest-api.yaml
```

### **Scenario 2**: Add order directly on Kafka

On this scenario, we want to validate only the event consumers and check if they are properly working:

```sh
tracetest run test -f test-payment-order-submit-with-message.yaml
```

### **Scenario 3**: Validate Risk Analysis for high value orders

This scenario checks if the Risk Analysis consumer flags a order sent into Kafka as risky and check if it was persisted correctly:

```sh
tracetest run test -f test-risk-analysis-using-order-with-high-value.yaml
```

### **Scenario 4**: Validate Risk Analysis for low value orders

This scenario checks if the Risk Analysis consumer flags a order sent into Kafka as normal and check if it was persisted correctly:

```sh
tracetest run test -f test-risk-analysis-using-order-with-low-value.yaml
```
