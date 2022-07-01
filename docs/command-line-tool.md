# Command Line Tool

Our web interface makes it easier to visualize your traces and add assertions, but sometimes we need a CLI to be able to automate our work. We developed the CLI thinking about users wanting to create tests and execute them everytime a new change is made in the system, so you could use Tracetest to detect regressions and check your service SLOs.

## Available commands

Here is a list of all available commands and how to use them:

### Configure
Configure your CLI to connect to your tracetest server

**How to use**: 

```sh
tracetest configure
```

or 
```sh
echo "your-server-url" | tracetest configure
```

### Test list

Allows you to list all tests

**How to use**:

```sh
tracetest test list
```

### Run a test
Allows you to run a test by referencing a [test definition file](/docs/test-definition-file.md).

> Note: if your definition file contain the field `id`, this command will not create a new test. Instead, it will update the test with that id. In case that test doesn't exist, a new one will be created with that id on your server.

Every time you run the test, we detect changes, and if any change is introduced, we use our [versioning](/docs/versioning.md) mechanism to ensure that it will not cause problems with your old test runs.

**How to use**:

```sh
tracetest test run --definition <file-path>
```

**Options**:

`--wait-for-result`: the CLI will only exit after the test run has completed (trace was retrieved and assertions were executed)