# Command Line Tool

Our web interface makes it easier to visualize your traces and add assertions, but sometimes we need a CLI to automate our work. We developed the CLI thinking about users creating tests and executing them every time a new change is made in the system, so you can use Tracetest to detect regressions and check your service SLOs.


## Available Commands


Here is a list of all available commands and how to use them:

### Configure
Configure your CLI to connect to your Tracetest server.


**How to Use**: 


```sh
tracetest configure
```

or 
```sh
echo "your-server-url" | tracetest configure
```

### Test List


Allows you to list all tests.


**How to Use**:


```sh
tracetest test list
```

### Run a Test

Allows you to run a test by referencing a [test definition file](/docs/test-definition-file.md).

> Note: If your definition file contains the field `id`, this command will not create a new test. Instead, it will update the test with that ID. If that test doesn't exist, a new one will be created with that ID on your server.


Every time you run the test, we detect changes, and if any change is introduced, we use our [versioning](/docs/versioning.md) mechanism to ensure that it will not cause problems with your old test runs.

**How to Use**:


```sh
tracetest test run --definition <file-path>
```

**Options**:

`--wait-for-result`: The CLI will only exit after the test run has completed (the trace was retrieved and assertions were executed).