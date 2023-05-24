# Creating Transactions

Just like other structures of Tracetest, you can also manager your transactions using the CLI and definition files.

A definition file for a transaction looks like the following:

```yaml
type: Transaction
spec:
    name: Test purchase flow
    description: Test a flow of purchasing an item
    steps:
      - ./tests/create-product.yaml
      - ./tests/add-product-to-cart.yaml
      - ./tests/complete-purschase.yaml
      - testID # you can also reference tests by their ids instead of referencing the definition file
```

In order to apply this transaction to your Tracetest instance, make sure to have your [CLI configured](./configuring-your-cli.md) and run:

```
tracetest apply transaction -f <transaction.yaml>
```

> If the file contains the property `spec.id`, the operation will be considered a transaction udpate.
