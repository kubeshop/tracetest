# Defining Environments as Text Files

This page showcases how to create and edit environments with the CLI.

:::tip
[To read more about environments check out environment concepts.](../concepts/variable-sets.md)
:::

Just like Data Stores, you can also manage your environments using the CLI and definition files.

A definition  file looks like the following:

```yaml
type: Environment
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
tracetest apply environment -f <environment.yaml>
```

> If the file contains the property `spec.id`, the operation will be considered an environment udpate. If you try to apply an environment and you get an error: `could not apply environment: 404 Not Found`, it means the provided id doesn't exist. Either update the id to reference an existing environment, or just remove the property from the file, so Tracetest will create a new environmenta and a new id.
