# Defining Test Specifications in Text Files

Test Specifications may be added to a trace to set a value for a step in the trace to determine success or failure.

## Assertions and Selectors

Assertions are as important as how you trigger your test. Without them, your test is just a fancy way of executing a request using a CLI command. In this section, we will discuss how you can declare your assertions in your definition file.

Before we start, there are two concepts that you must understand to write your tests:

- [Selectors](../concepts/selectors.md)
- [Assertions](../concepts/assertions.md)

### Selectors

**Selectors** are queries that are executed against your trace tree and select a set of spans based on some attributes. They are responsible for defining which spans will be tested against your assertions.

### Assertions

**Assertions** are tests against a specific span based on its attributes. A practical example might be useful:

Imagine you have to ensure that all your `database select statements` take `less than 500ms`. To write a test for that you must:

1. Select all spans in your trace related to `select statements`.
2. Check if all those spans lasted `less than 500ms`.

For the first task, we use a selector: `span[db.statement contains "SELECT"]`. While the second one is achieved by using an assertion: `attr:tracetest.span.duration < 500ms`.

> **Note:** When asserting time fields, you can use the following time units: `ns` (nanoseconds), `us` (microseconds), `ms` (milliseconds), `s` (seconds), `m` (minutes), and `h` (hours). Instead of defining `attr:tracetest.span.duration <= 3600s`, you can set it as `attr:tracetest.span.duration <= 1h`.

To write that in your test definition, you can define the following YAML definition:

```yaml
specs:
- selector: span[db.statement contains "SELECT"]
  assertions:
    - attr:tracetest.span.duration < 500ms
```

As you probably noticed in the test definition structure, you can have multiple assertions for the same selector. This is useful to group related validations. For example, ensuring that all your HTTP calls are successful and take less than 1000ms:

```yaml
specs:
- selector: span[tracetest.span.type="http"]
  assertions:
    - attr:http.status_code >= 200
    - attr:http.status_code < 300
    - attr:tracetest.span.duration < 1000ms
```

#### Referencing Other Fields from the Same Span

You also can reference fields from the same span in your assertions. For example, you can define an assertion to ensure the output number is greater than the input number.

```yaml
specs:
- selector: span[name = "my operation"]
  assertions:
    - attr:myapp.output > attr:myapp.input
```

You also can use basic arithmetic expressions in your assertions:

```yaml
assertions:
  - attr:myapp.output = myapp.input + 1
```

:::note
This does not take into account the order of operators yet. So an expression `1 + 2 * 3` will be resolved as `9` instead of `7`. This will be fixed in future releases.
:::

Available operations in an expression are: `+`, `-`, `*`, and `/`.

For more information about selectors or assertions, take a look at the documentation for those topics.

## Available Operations

| Operator              | Description |
| :------------- | ------------------------------------------------------------------------------------------------------------------------ |
| `=`            | Check if two values are equal.                                                                                            |
| `!=`           | Check if two values have different values.                                                                                |
| `<`            | Check if value from left side is smaller than the one on the right side of the operation.                                 |
| `<=`           | Check if value from left side is smaller or equal to the one on the right side of the operation.                          |
| `>`            | Check if value from left side is larger than the one on the right side of the operation.                                  |
| `>=`           | Check if value from left side is larger or equal to the one on the right side of the operation.                           |
| `contains`     | Check if value on the right side of the operation is contained inside of the value of the left side of the operation.     |
| `not-contains` | Check if value on the right side of the operation is not contained inside of the value of the left side of the operation. |

## Testing Span Events

As an MVP of how to test span events, we are injecting `all spans` events into the span attributes as a JSON array. To assert your span events, use the `json_path` filter to select and test the **write events**.

```yaml
specs:
  - selector: span[name = "my span"]
    assertions:
      - attr:span.events | json_path '$[?(@.name = "event name")].attributes.key' = "expected_value"
```

### Query breakdown

* **@.name = "event name"**: select the event with name "**event name**"
* $[?(@.name = "event name")]**.attributes.key**: select the attribute "**key**" from the event with name "**event name**"
