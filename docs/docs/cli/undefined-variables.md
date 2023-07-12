# Undefined Variables

When a user runs a test or a transaction, any variables that will be needed but are not defined will be prompted for:

```sh
tracetest run test -f path/to/test.yaml
```

```text title="Output:"
WARNING  Some variables are required by one or more tests
INFO  Fill the values for each variable:
POKEID:
POKENAME:
```

Undefined variables are dependent on the environment selected and whether or not the variable is defined in the current environment. Select the environment to run the test or transaction by passing it into the test run command.

```sh
tracetest list environment
```

```text title="Output:"
 ID        NAME      DESCRIPTION
--------- --------- -------------
 testenv   testenv   testenv
```

```sh
tracetest get environment --id testenv
```

```text title="Output:"
---
type: Environment
spec:
  id: testenv
  name: testenv
  description: testenv
  values:
  - key: POKEID
    value: "42"
  - key: POKENAME
    value: oddish
```

```sh
tracetest run test -f path/to/test.yaml -e testenv
```

```text title="Output:"
✔ Pokeshop - Import (http://your-domain:11633/test/XRHjfH_4R/run/XYZ/test)
```

:::tip
[Check out use-cases for using undefined variables here.](../concepts/ad-hoc-testing.md)
:::
