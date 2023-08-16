# Undefined Variables

When a user runs a test or a Test Suite, any variables that will be needed but are not defined will be prompted for:

```sh
tracetest run test -f path/to/test.yaml
```

```text title="Output:"
WARNING  Some variables are required by one or more tests
INFO  Fill the values for each variable:
POKEID:
POKENAME:
```

Undefined variables are dependent on the variable set selected and whether or not the variable is defined in the current variable set. Select the variable set to run the Test or Test Suite by passing it into the test run command.

```sh
tracetest list variableset
```

```text title="Output:"
 ID        NAME      DESCRIPTION
--------- --------- -------------
 testvars   testvars   testvars
```

```sh
tracetest get variableset --id testvars
```

```text title="Output:"
---
type: VariableSet
spec:
  id: testvars
  name: testvars
  description: testvars
  values:
  - key: POKEID
    value: "42"
  - key: POKENAME
    value: oddish
```

```sh
tracetest run test -f path/to/test.yaml --vars testvars
```

```text title="Output:"
✔ Pokeshop - Import (http://your-domain:11633/test/XRHjfH_4R/run/XYZ/test)
```

:::tip
[Check out use-cases for using undefined variables here.](../concepts/ad-hoc-testing.md)
:::
