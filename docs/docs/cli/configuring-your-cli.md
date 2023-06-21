# Configuring your CLI

Our web interface makes it easier to visualize your traces and add assertions, but sometimes a CLI is needed for automation. The CLI was developed for users creating tests and executing them each time a change is made in the system, so Tracetest can detect regressions and check service Service Level Objectives (SLOs).


## **Available Commands**

Here is a list of all available commands and how to use them:

### **Configure**
Configure your CLI to connect to your Tracetest server.


**How to Use**:


```sh
tracetest configure
```

If you want to set values without having to answer questions from a prompt, you can provide the flag `--endpoint` to define the server endpoint.


```sh
tracetest configure --endpoint http://my-tracetest-server:11633
```

### **Test List**


Allows you to list all tests.


**How to Use**:


```sh
tracetest test list
```

### **Run a Test**

Allows you to run a test by referencing a [test definition file](./creating-tests).

> Note: If the definition file contains the field `id`, this command will not create a new test. Instead, it will update the test with that ID. If that test doesn't exist, a new one will be created with that ID on the server.


Every time the test is run, changes are detected and, if any change is introduced, we use Tractest's [versioning](../concepts/versioning) mechanism to ensure that it will not cause problems with previous test runs.

**How to Use**:


```sh
tracetest test run --definition <file-path>
```

**Options**:

`--wait-for-result`: The CLI will only exit after the test run has completed (the trace was retrieved and assertions were executed).

### Running Tracetest CLI From Docker

There are times when it is easier to directly execute the Tracetest CLI from a Docker image rather than installing the CLI on your local machine. This can be convenient when you wish to execute the CLI in a CI/CD environment.


**How to Use**:

Use the command below, substituting the following placeholders:
- <your-tracetest-server-url> - The URL to the running Tracetest server you wish to execute the test on. Example: http://localhost:11633/

- <file-path> - The path to the saved Tracetest test. Example: ./mytest.yaml


```sh
docker run --rm -it -v$(pwd):$(pwd) -w $(pwd) --network host --entrypoint tracetest kubeshop/tracetest:latest -s <your-tracetest-server-url> test run  --definition <file-path> --wait-for-result
```


