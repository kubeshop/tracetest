# Running Tracetest with Keptn
 
[Keptn](https://keptn.sh/) is a powerful tool to automate the lifecycle of your application running on Kubernetes. One of the tasks that we can do on `keptn` is to test an application and see if it is healthy and ready to be used by your users.

By using Keptn [Job Executor Service](https://github.com/keptn-contrib/job-executor-service) plugin, we can upload a Tracetest test definition and a CLI configuration to a service, and run a test using the following job:

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
         - sh
       args:
         - -c
         - "./tracetest --config /keptn/data/tracetest-cli-config.yaml test run --definition /keptn/data/test-definition.yaml --wait-for-result"
```

## Quickstart

Here we will show how to use Tracetest to do these tests and help in your delivery and testing workflows. 

### Assumptions

1. We are using the [Pokeshop](https://docs.tracetest.io/pokeshop/) example, exposed on `http://demo-pokemon-api.demo`;
2. Tracetest is installed in `tracetest` namespace in that cluster;
3. We are considering that Keptn is [already installed](https://keptn.sh/docs/1.0.x/install/) with the [Job Executor Service](https://github.com/keptn-contrib/job-executor-service) plugin on your Kubernetes cluster and that you have the [CLI](https://keptn.sh/docs/1.0.x/install/cli-install/) installed on your machine and already [authenticated](https://keptn.sh/docs/1.0.x/install/authenticate-cli-bridge/) with Keptn API.
 
### 1. Setup a project and a service
 
Keptn works with [concepts](https://keptn.sh/docs/concepts/glossary/) of a Project (element to maintain multiple services forming an application in stages) and a Service (smallest deployable unit and is deployed in all project stages according to the order).
Usually, these resources are managed by Keptn during Sequences (a set of tasks for realizing a delivery or operations process). An example of a sequence of **Deploying**** a system, that could do 3 tasks:
1. Update a service / a set of services in a namespace
2. Run tests to see if everything is working fine
3. Change traffic to this new set of services
 
To integrate Tracetest with Keptn at first glance, we recommend setting up a task that will invoke Tracetest to run tests on these services.
 
In our example, we are setting up a project through the CLI by first, defining a `shipyard` file called [`shipyard.yaml`](./shipyard.yaml) for our project:
```yaml
apiVersion: "spec.keptn.sh/0.2.2"
kind: "Shipyard"
metadata:
 name: "shipyard-keptn-tracetest-integration"
spec:
 stages:
   - name: "production"
     sequences:
       - name: "deployment"
         tasks:
           - name: "update-services"
           - name: "test-services"
```
 
And later creating a project `keptn-tracetest-integration` with it:
```sh
keptn create project keptn-tracetest-integration -y -s shipyard.yaml
```
 
**Note:** Keptn may ask you to have a git repository for this project, to enable GitOps. If so, you need to create an empty git repository, and a git token and pass it through the flags `--git-remote-url`, `-`-git-user`, and `--`git-token`. More details about this setup can be seen on [Keptn docs/Git-based upstream](https://keptn.sh/docs/1.0.x/manage/git_upstream).
 
After that, we can create a service `pokeshop` for this project by executing:
```sh
keptn create service pokeshop --project keptn-tracetest-integration -y
```
 
### 2. Add Tracetest files and job files as resources of a service
 
After creating a service, we will set up a job associated with the `pokeshop` service and the task event `test-services`.
 
To do that, first, we will create a config file for Tracetest CLI, telling which instance of Tracetest should run the tests. The file should be saved in your current directory with the name [`tracetest-cli-config.yaml`](./tracetest-cli-config.yaml):
```yaml
scheme: http
endpoint: tracetest.tracetest.svc.cluster.local:11633
analyticsEnabled: false
```
 
Then, we will create a Tracetest test definition file in your current directory called [`test-definition.yaml`](./test-definition.yaml):
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
 
Now we will add these files as resources for the service `pokeshop` on Keptn with the following command:
```sh
keptn add-resource --project keptn-tracetest-integration --service pokeshop --stage production --resource test-definition.yaml --resourceUri data/test-definition.yaml
keptn add-resource --project keptn-tracetest-integration --service pokeshop --stage production --resource tracetest-cli-config.yaml --resourceUri data/tracetest-cli-config.yaml
```
 
These files will be located on the folder `data` and will be injected into our Keptn job that we will set up in the next step. To that, we will set up a job definition YAML, telling to run Tracetest every time that an event of `test-services` happens. This file will be named [`job-config.yaml`](./job-config.yaml), and will listen for the event `sh.keptn.event.test-services.triggered` (event emitted by the `test-services` task on the `deployment` sequence defined in our project):
```yaml
apiVersion: v2
actions:
 - name: "Run tracetest on your service"
   events:
     - name: "sh.keptn.event.test-services.triggered"
   tasks:
     - name: "Run tracetest"
       files:
         - data/test-definition.yaml
         - data/tracetest-cli-config.yaml
       image: "kubeshop/tracetest:latest"
       cmd:
         - sh
       args:
         - -c
         - "./tracetest --config /keptn/data/tracetest-cli-config.yaml test run --definition /keptn/data/test-definition.yaml --wait-for-result"
```
 
Finally, we will add this job as a resource on Keptn:
```sh
keptn add-resource --project keptn-tracetest-integration --service pokeshop --stage production --resource job-config.yaml --resourceUri job/config.yaml
```
 
### 3. Setup Job Executor Service to see events emitted by the test step
 
To guarantee that our job will be called by Keptn when we execute the `deployment` sequence, we need to configure the integration `Job Executor Service` on `keptn-tracetest-integration` project to listen to `sh.keptn.event.test-services.triggered` events. We can do that only through the Keptn Bridge (their Web UI), by going to our project, choosing the `Settings` option, and later `Integrations``.
 
Choose the `job-executor-service` integration, and add a subscription to the event `sh.keptn.event.test-services.triggered` and the project `keptn-tracetest-integration`.
 
### 4. Run sequence when needed
 
At last, to see the integration running, we only need to execute the following command:
```sh
keptn trigger sequence deployment --project keptn-tracetest-integration --service pokeshop --stage production
```

Now you should be able to see the sequence running for `keptn-tracetest-integration` project on Keptn Bridge.
