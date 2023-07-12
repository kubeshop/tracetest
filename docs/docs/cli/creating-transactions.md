# Defining Transactions as Text Files

This page showcases how to create and edit Transactions with the CLI.

:::tip
[To read more about transactions check out transactions concepts.](../concepts/transactions.md)
:::

Just like other structures of Tracetest, you can also manage your transactions using the CLI and definition files.

A definition file for a transaction looks like the following:

```yaml
type: Transaction
spec:
  name: Test purchase flow
  description: Test a flow of purchasing an item
  steps:
    - ./tests/create-product.yaml
    - ./tests/add-product-to-cart.yaml
    - ./tests/complete-purchase.yaml
    - testID # you can also reference tests by their ids instead of referencing the definition file
```

In order to apply this transaction to your Tracetest instance, make sure to have your [CLI configured](./configuring-your-cli.md) and run:

```sh
tracetest apply transaction -f <transaction.yaml>
```

> If the file contains the property `spec.id`, the operation will be considered a transaction udpate.
