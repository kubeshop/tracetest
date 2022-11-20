# Command Line Tool

Our web interface makes it easier to visualize your traces and add assertions, but sometimes a CLI is needed for automation. The CLI was developed for users creating tests and executing them each time a change is made in the system, so Tracetest can detect regressions and check service SLOs.


## **Available Commands**


Here is a list of all available commands and how to use them:

### **Configure**
Configure your CLI to connect to your Tracetest server.


**How to Use**:


```sh
tracetest configure
```

If you want to set values without having to answer questions from a prompt, you can provide the flags `--endpoint` to define the server endpoint and `--analytics` to turn the analytics on and off.


```sh
# This will prompt a question to ask if you want to enable or not analytics
tracetest configure --endpoint http://my-tracetest-server:11633

# Analytics enabled
tracetest configure --endpoint http://my-tracetest-server:11633 --analytics

# Analytics disabled
tracetest configure --endpoint http://my-tracetest-server:11633 --analytics=false
```

### **Test List**


Allows you to list all tests.


**How to Use**:


```sh
tracetest test list
```

### **Run a Test**

Allows you to run a test by referencing a [test definition file](./test-definition-file).

> Note: If the definition file contains the field `id`, this command will not create a new test. Instead, it will update the test with that ID. If that test doesn't exist, a new one will be created with that ID on the server.


Every time the test is run, changes are detected and, if any change is introduced, we use Tractest's [versioning](../using-tracetest/versioning) mechanism to ensure that it will not cause problems with previous test runs.

**How to Use**:


```sh
tracetest test run --definition <file-path>
```

**Options**:

`--wait-for-result`: The CLI will only exit after the test run has completed (the trace was retrieved and assertions were executed).
