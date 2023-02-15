# Running Tracetest with Keptn
 
[Tracetest](https://tracetest.io/) is a testing tool based on OpenTelemetry that permits you to test your distributed application. It allows you to use your trace data generated on your OpenTelemetry tools to check and assert if your application has the desired behavior defined by your test definitions.

[Keptn](https://keptn.sh/) is a powerful tool to automate the lifecycle of your application running on Kubernetes. One of the tasks that we can do with `keptn` is to test an application and see if it is healthy and ready to be used by your users.

By using the Keptn [Job Executor Service](https://github.com/keptn-contrib/job-executor-service) plugin, we can upload a Tracetest test definition and a CLI configuration to a service and run a test using the following job:

```yaml
apiVersion: v2
actions:
- name: "Run tracetest on your service"
  events:
    - name: "sh.keptn.event.test.triggered"
  tasks:
    - name: "Run tracetest"
      files:
        - data/test-definition.yaml
        - data/tracetest-cli-config.yaml
      image: "kubeshop/tracetest:latest"
      cmd:
        - tracetest
      args:
        - --config
        - /keptn/data/tracetest-cli-config.yaml
        - test
        - run
        - --definition
        - /keptn/data/test-definition.yaml
        - --wait-for-result

```

## Quickstart

Here we will show how to use Tracetest to run these tests and help in your delivery and testing workflows using the [Pokeshop](https://docs.tracetest.io/pokeshop/) example, available at `http://demo-pokemon-api.demo` in a Kubernetes cluster.

### Prerequisites

In your Kubernetes cluster you should have:

1. `Keptn 1.0.x` [installed](https://keptn.sh/docs/1.0.x/install/)
2. `Job Executor Service 0.3.x` [installed](https://github.com/keptn-contrib/job-executor-service/blob/main/docs/INSTALL.md)
3. `Tracetest` server [installed](https://docs.tracetest.io/deployment/kubernetes) on `tracetest` namespace 

On your machine you should have:

1. `Keptn CLI` [installed](https://keptn.sh/docs/1.0.x/install/cli-install/)
2. and already [authenticated](https://keptn.sh/docs/1.0.x/install/authenticate-cli-bridge/) with Keptn API.

With everything set up, we will start configuring Keptn and Tracetest.

### 1. Setup a project and a service.
 
Keptn works with [concepts](https://keptn.sh/docs/concepts/glossary/) of a Project (element to maintain multiple services forming an application in stages) and a Service (smallest deployable unit and is deployed in all project stages according to the order).

Usually, these resources are managed by Keptn during Sequences (a set of tasks for realizing a delivery or operations process). In this example, we will create a sequence and run a Tracetest test in a system, once the sequence is triggered.
 
To do that, we will do the following steps:

1. Create a skeletal [`shipyard.yaml`](https://github.com/kubeshop/tracetest/tree/main/examples/keptn-integration/shipyard.yaml) file with the following content:
```yaml
apiVersion: "spec.keptn.sh/0.2.2"
kind: "Shipyard"
metadata:
 name: "shipyard-keptn-tracetest-integration"
spec:
 stages:
   - name: "production"
     sequences:
       - name: "validate-pokeshop"
         tasks:
           - name: "test"
```

2. Create a new `keptn-tracetest-integration` project using that `shipyard` file:
```sh
keptn create project keptn-tracetest-integration -y -s shipyard.yaml
```
 
**Note:** Keptn may ask you to have a Git repository for this project to enable GitOps. If so, you need to create an empty Git repository and a Git token and pass it through the flags `--git-remote-url`, `--git-user`, and `--git-token`. More details about this setup can be seen on [Keptn docs/Git-based upstream](https://keptn.sh/docs/1.0.x/manage/git_upstream).
 
3. Create a `pokeshop` service:
```sh
keptn create service pokeshop --project keptn-tracetest-integration -y
```
 
### 2. Add Tracetest files and job files as resources of a service.
 
Now, we will set up a job associated with the `pokeshop` service, listening to the task event `test-services`:
 
1. Create the [`tracetest-cli-config.yaml`](https://github.com/kubeshop/tracetest/tree/main/examples/keptn-integration/tracetest-cli-config.yaml) configuration file for the Tracetest CLI in your current directory, identifying the Tracetest instance that should run the tests:
```yaml
scheme: http
endpoint: tracetest.tracetest.svc.cluster.local:11633
```

2. Create the [`test-definition.yaml`](https://github.com/kubeshop/tracetest/tree/main/examples/keptn-integration/test-definition.yaml) Tracetest test definition in your current directory:
```yaml
type: Test
spec:
  id: apdCx-h4g
  name: Pokeshop - List Pokemons
  description: Get a Pokemon
  trigger:
    type: http
    httpRequest:
      url: http://demo-pokemon-api.demo/pokemon?take=20&skip=0
      method: GET
      headers:
      - key: Content-Type
        value: application/json
  specs:
  - selector: span[name = “Tracetest trigger”]
    assertions:
    - attr:tracetest.span.duration < 500ms
```
 
3. Add these files as resources for the `pokeshop` service:
```sh
keptn add-resource --project keptn-tracetest-integration --service pokeshop --stage production --resource test-definition.yaml --resourceUri data/test-definition.yaml
keptn add-resource --project keptn-tracetest-integration --service pokeshop --stage production --resource tracetest-cli-config.yaml --resourceUri data/tracetest-cli-config.yaml
```

These files will be located in the folder `data` and will be injected into our Keptn job that we will set up in the next step.

4. Create [`job-config.yaml`](https://github.com/kubeshop/tracetest/tree/main/examples/keptn-integration/job-config.yaml):
```yaml
apiVersion: v2
actions:
- name: "Run tracetest on your service"
  events:
    - name: "sh.keptn.event.test.triggered"
  tasks:
    - name: "Run tracetest"
      files:
        - data/test-definition.yaml
        - data/tracetest-cli-config.yaml
      image: "kubeshop/tracetest:latest"
      cmd:
        - tracetest
      args:
        - --config
        - /keptn/data/tracetest-cli-config.yaml
        - test
        - run
        - --definition
        - /keptn/data/test-definition.yaml
        - --wait-for-result
```

This job will run Tracetest every time a `test` event happens, listening to the event `sh.keptn.event.test.triggered` (event emitted by the `test` task on the `validate-pokeshop` sequence).

5. Add this job as a resource on Keptn:
```sh
keptn add-resource --project keptn-tracetest-integration --service pokeshop --stage production --resource job-config.yaml --resourceUri job/config.yaml
```
 
### 3. Setup Job Executor Service to see events emitted by the test step.
 
To guarantee that our job will be called by Keptn when we execute the `deployment` sequence, we need to configure the integration `Job Executor Service` on `keptn-tracetest-integration` project to listen to `sh.keptn.event.test.triggered` events if it is not configured. We can do that only through the Keptn Bridge (their Web UI), by going to our project, choosing the `Settings` option, and later `Integrations`.
 
Choose the `job-executor-service` integration, and add a subscription to the event `sh.keptn.event.test.triggered` and the project `keptn-tracetest-integration`.
 
### 4. Run sequence when needed.
 
Finally, to see the integration running, we only need to execute the following command:
```sh
keptn trigger sequence test-sequence --project keptn-tracetest-integration --service pokeshop --stage production
```

Now you should be able to see the sequence running for `keptn-tracetest-integration` project on Keptn Bridge.
