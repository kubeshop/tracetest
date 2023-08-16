# Defining Test Suites as Text Files

This page showcases how to create and edit Test Suites with the CLI.

:::tip
[To read more about Test Suites check out Test Suites concepts page.](../concepts/test-suites.md)
:::

Just like other structures of Tracetest, you can also manage your Test Suites using the CLI and definition files.

A definition file for a Test Suite looks like the following:

```yaml
type: TestSuite
spec:
  name: Test purchase flow
  description: Test a flow of purchasing an item
  steps:
    - ./tests/create-product.yaml
    - ./tests/add-product-to-cart.yaml
    - ./tests/complete-purchase.yaml
    - testID # you can also reference tests by their ids instead of referencing the definition file
```

In order to apply this Test Suite to your Tracetest instance, make sure to have your [CLI configured](./configuring-your-cli.md) and run:

```sh
tracetest apply testsuite -f <testsuite.yaml>
```

> If the file contains the property `spec.id`, the operation will be considered a Test Suite udpate.
