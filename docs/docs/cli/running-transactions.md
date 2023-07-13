# Running Transactions From the Command Line Interface (CLI)

Once you have created a transaction, whether from the Tracetest UI or via a text editor, you will need the capabity to run it via the Command Line Interface (CLI) to integrate it into your CI/CD process or your local development workflow.

The command to run a transaction is the same as running a test from the CLI.

The documentation for running a test via the CLI can be found here:

- [tracetest run](./reference/tracetest_run.md): This page provides examples of using this command.

## Running Your First Transaction

To run a transaction, give the path to the transaction definition file with the `'-f'` option. This will launch a transaction, providing us with a link to the created transaction run.

```sh
tracetest run transaction -f path/to/transaction.yaml
```
```text title="Output:"
✔ Pokemon Transaction (http://localhost:11633/transaction/xcGqfHl4g/run/3)
	✔ Pokeshop - Import (http://localhost:11633/test/XRHjfH_4R/run/4/test)
	✔ Pokeshop - List (http://localhost:11633/test/QvPjBH_4g/run/4/test)
```

## Running a Transaction That Uses Environment Variables

There are two ways of referencing an environment when running a transaction.

You can reference an existing environment using its id. For example, given this defined environment with an id of `'testenv'`:

![testenv](../img/show-environment-definition.png)

We can run a transaction and specify that environment with this command:

```sh
tracetest run transaction -f path/to/transaction.yaml -e testenv
```

You can also reference an environment resource file which will be used to create a new environment or update an existing one. For example, if you have a file named `local.env` with this content:

```yaml
type: Environment
spec:
  id: local.env
  name: local.env
  values:
  - key: POKEID
    value: 45
  - key: POKENAME
    value: vileplume
```

```sh
tracetest run test -f path/to/transaction.yaml -e path/to/local.env
```

If you use the environment resource approach, a new environment will be created in Tracetest.

The second approach is very useful if you are running tests from a CI pipeline.
