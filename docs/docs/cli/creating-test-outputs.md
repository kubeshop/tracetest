# Defining Test Outputs in Text Files

Outputs are really useful when running [Test Suites](../concepts/test-suites). They allow for exporting values from a test so they become available in the [Variable Sets](../concepts/variable-sets.md) of the current Test Suite.

## Outputs are Expression Results

An output exports the result of an [Expression](../concepts/expressions) and assigns it to a name, so it can be injected into the variable set of a running Test Suite.
A `selector` is needed only if the provided expression refers to a/some span/s attribute or meta attributes.

It can be defined using the following YAML definition:

```yaml
outputs:
  - name: USER_ID
    selector: span[name = "user creation"]
    value: attr:myapp.users.created_id
```

The `value` attribute is an `expression` and is a very powerful tool.

## Examples

### Basic Expression

You can output basic expressions:

```yaml
outputs:

- name: ARITHMETIC_RESULT
  value: 1 + 1
  # results in ARITHMETIC_RESULT = 2

- name: INTERPOLATE_STRING
  # assume PRE_EXISTING_VALUE=someValue from env vars
  value: "the value ${env:PRE_EXISTING_VALUE} comes from the env var PRE_EXISTING_VALUE"
  # results in INTERPOLATE_STRING = "the value someValue comes from the env var PRE_EXISTING_VALUE
```

### Extract a Value from a JSON

Imagine a hypotetical `/users/create` endpoint that returns the full `user` object, including the new ID, when the operation is successful.

```yaml
outputs:
- name: USER_ID
  selector: span[name = "POST /user/create"]
  value: attr:http.response.body | json_path '.id'
```

### Multiple Values

Using the same hypotethical user creation endpoint, a user creation might result on multiple sql queries, for example:

- `INSERT INTO users ...`
- `INSERT INTO permissions...`
- `SELECT remaining_users FROM accounts`
- `UPDATE accounts SET remaining_users ...`

In this case, the service is instrumented so that each query generates a span of type `database`.
You can get a list of SQL operations:

```yaml
outputs:
- name: SQL_OPS
  selector: span[tracetest.span.type = "database"]
  value: attr:sql.operation
  # result: SQL_OPS = ["INSERT", "INSERT", "SELECT", "UPDATE"]
```

Since the value is an array, you can also apply filters to it:

```yaml
outputs:
- name: LAST_SQL_OP
  selector: span[tracetest.span.type = "database"]
  value: attr:sql.operation | get_index 'last'
  # result: LAST_SQL_OP = "INSERT"
```
