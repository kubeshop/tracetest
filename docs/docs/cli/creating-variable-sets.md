# Defining Variable Sets as Text Files

This page showcases how to create and edit variable sets with the CLI.

:::tip
[To read more about variable sets check out variable sets concepts.](../concepts/variable-sets.md)
:::

Just like Data Stores, you can also manage your variable sets using the CLI and definition files.

A definition  file looks like the following:

```yaml
type: VariableSet
spec:
    name: Production
    description: Production env variables
    values:
    - key: URL
      value: https://app-production.company.com
    - key: API_KEY
      value: mysecret
```

In order to apply this configuration to your Tracetest instance, make sure to have your [CLI configured](./configuring-your-cli.md) and run:

```sh
tracetest apply variableset -f <variableset.yaml>
```

> If the file contains the property `spec.id`, the operation will be considered a variable set update. If you try to apply a variable set and you get an error: `could not apply variableset: 404 Not Found`, it means the provided id doesn't exist. Either update the id to reference an existing variable set, or just remove the property from the file, so Tracetest will create a new variable set and a new id.
