# Creating Transactions

This page showcases how to create and edit Transactions in the Web UI.

:::tip
[To read more about transactions check out transactions concepts.](../concepts/transactions.md)
:::

![Main Screen](../img/main-screen-0.11.png)

Click the **Create** button and select **Create New Transaction** in the drop down:

![Create a Test Button](../img/create-button-0.11.png)

Give your transaction a name, and click **Next**:

![Name the Transaction](https://res.cloudinary.com/djwdcmwdz/image/upload/v1685712802/docs/beta.tracetest.io__page_1_jynf6o.png)

Next, select which tests to run in the transaction and click **Create & Run**:

![Select Tests in Transaction](https://res.cloudinary.com/djwdcmwdz/image/upload/v1685712954/docs/beta.tracetest.io__page_1_1_agjvg0.png)

The transaction will start:

![Running Transaction](../img/running-transaction.png)

 On the automate tab, you find methods to automate the current transaction, including the YAML test file and the CLI command for Tracetest.

 ![Automate Tab](../img/automate-tab.png)

 With all of the toggles `Off`, each criteria is tested. Toggle individual criteria `On`, to select on certain checks for the selected test.

When the transaction is finished, you will get the following result:

![Finished Transaction](https://res.cloudinary.com/djwdcmwdz/image/upload/v1685713712/docs/demo.tracetest.io__x0o1gu.png)

You can now view individual [Test Results](test-results.md) executed by the transaction by clicking on any of the tests in the list.
