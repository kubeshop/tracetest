# Testkube Tracetest Executor

Testkube Tracetest Executor is a test executor to run Tracetest tests with Testkube.

**Why?**

Tracetest leverages existing OpenTelemetry instrumentation to run assertions against every part of an HTTP transaction. But, by integrating with Testkube you can now add Tracetest to the native CI/CD pipeline in your Kubernetes cluster. This allows you to run scheduled test runs and synthetic tests. All while following the trace-based testing principle and enabling full in-depth assertions against trace data, not just the response.

## Prerequisites

1. Make sure you have a running Kubernetes cluster.
2. Kubectl
3. Helm

## Quickstart

Follow these steps to get started.

## 1. Install Testkube CLI

Install Testkube CLI by following [these instructions](https://docs.testkube.io/getting-started/step1-installing-cli) for your OS.

```bash
# MacOS example
brew install testkube
```

## 2. Install Testkube in your Kubernetes cluster

```bash
testkube init
```

Confirm that Testkube is running:

```bash
kubectl get all -n testkube

[Output]
NAME                                                        READY   STATUS    RESTARTS      AGE
pod/testkube-api-server-8f5cf8b8f-vrpk6                     1/1     Running   3 (14m ago)   18m
pod/testkube-dashboard-584846b754-4sxq8                     1/1     Running   0             18m
pod/testkube-minio-testkube-bd549c85d-lpwbm                 1/1     Running   0             18m
pod/testkube-mongodb-d78699775-cjcd7                        1/1     Running   0             18m
pod/testkube-nats-0                                         3/3     Running   0             18m
pod/testkube-nats-box-5b555bc9c4-5mpz5                      1/1     Running   0             18m
pod/testkube-operator-controller-manager-5f84bb4fd4-mjr2l   2/2     Running   0             18m

NAME                                                           TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                                                 AGE
service/testkube-api-server                                    ClusterIP   10.96.78.171    <none>        8088/TCP                                                18m
service/testkube-dashboard                                     ClusterIP   10.96.204.214   <none>        8080/TCP                                                18m
service/testkube-minio-service-testkube                        ClusterIP   10.96.85.219    <none>        9000/TCP,9090/TCP,9443/TCP                              18m
service/testkube-mongodb                                       ClusterIP   10.96.213.36    <none>        27017/TCP                                               18m
service/testkube-nats                                          ClusterIP   None            <none>        4222/TCP,6222/TCP,8222/TCP,7777/TCP,7422/TCP,7522/TCP   18m
service/testkube-operator-controller-manager-metrics-service   ClusterIP   10.96.243.219   <none>        8443/TCP                                                18m
service/testkube-operator-webhook-service                      ClusterIP   10.96.106.138   <none>        443/TCP                                                 18m

NAME                                                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/testkube-api-server                    1/1     1            1           18m
deployment.apps/testkube-dashboard                     1/1     1            1           18m
deployment.apps/testkube-minio-testkube                1/1     1            1           18m
deployment.apps/testkube-mongodb                       1/1     1            1           18m
deployment.apps/testkube-nats-box                      1/1     1            1           18m
deployment.apps/testkube-operator-controller-manager   1/1     1            1           18m

NAME                                                              DESIRED   CURRENT   READY   AGE
replicaset.apps/testkube-api-server-8f5cf8b8f                     1         1         1       18m
replicaset.apps/testkube-dashboard-584846b754                     1         1         1       18m
replicaset.apps/testkube-minio-testkube-bd549c85d                 1         1         1       18m
replicaset.apps/testkube-mongodb-d78699775                        1         1         1       18m
replicaset.apps/testkube-nats-box-5b555bc9c4                      1         1         1       18m
replicaset.apps/testkube-operator-controller-manager-5f84bb4fd4   1         1         1       18m

NAME                             READY   AGE
statefulset.apps/testkube-nats   1/1     18m
```

By default, Testkube is installed in the `testkube` namespace.

To explore the Testkube dashboard, run the command:

```bash
testkube dashboard
```

![Testkube Dashboard](https://res.cloudinary.com/djwdcmwdz/image/upload/v1679064043/Blogposts/Docs/screely-1679064032115_bnl4pc.png)

## 3. Install Tracetest CLI

Install Tracetest CLI by following [these instructions](https://docs.tracetest.io/getting-started/installation) for your OS.

```bash
# MacOS example
brew install kubeshop/tracetest/tracetest
```

## 4. Install Tracetest in your Kubernetes cluster

```bash
tracetest server install

[Output]
How do you want to run TraceTest? [type to search]:
  Using Docker Compose
> Using Kubernetes
```

Confirm that Tracetest is running:

```bash
kubectl get all -n tracetest

[Output]
NAME                                  READY   STATUS    RESTARTS        AGE
pod/otel-collector-7f4d87489f-vp6zn   1/1     Running   0               5m41s
pod/tracetest-78b9c84c57-t4prx        1/1     Running   3 (4m15s ago)   5m29s
pod/tracetest-postgresql-0            1/1     Running   0               5m42s

NAME                              TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)               AGE
service/otel-collector            ClusterIP   10.96.173.226   <none>        4317/TCP              5m46s
service/tracetest                 ClusterIP   10.96.248.146   <none>        11633/TCP,21321/TCP   5m42s
service/tracetest-postgresql      ClusterIP   10.96.155.147   <none>        5432/TCP              5m42s
service/tracetest-postgresql-hl   ClusterIP   None            <none>        5432/TCP              5m42s

NAME                             READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/otel-collector   1/1     1            1           5m46s
deployment.apps/tracetest        1/1     1            1           5m42s

NAME                                        DESIRED   CURRENT   READY   AGE
replicaset.apps/otel-collector-7f4d87489f   1         1         1       5m46s
replicaset.apps/tracetest-78b9c84c57        1         1         1       5m42s

NAME                                    READY   AGE
statefulset.apps/tracetest-postgresql   1/1     5m42s
```

By default, Testkube is installed in the `tracetest` namespace.

To explore the Tracetest Web UI, run the command:

```bash
kubectl --kubeconfig ${HOME}/.kube/config --context kind-kind --namespace tracetest port-forward svc/tracetest 11633
```

![Tracetest Web UI](https://res.cloudinary.com/djwdcmwdz/image/upload/v1679064296/Blogposts/Docs/screely-1679064291876_jxlhmn.png)

## 5. Create a test in Tracetest

Start by clicking `Create` > `Create New Test` > `HTTP Request` > `Next` > `Choose Example` (dropdown) > `Pokeshop - List` (generates a sample test from the Tracetest demo) > `Next` > `URL` is prefilled with `http://demo-pokemon-api.demo/pokemon?take=20&skip=0` > `Create & Run`

This will trigger the test and display a distributed trace in the `Trace` tab to run assertions against.

![Tracetest distributed trace test run view](https://res.cloudinary.com/djwdcmwdz/image/upload/v1679064990/Blogposts/Docs/screely-1679064984975_s0psbr.png)

Proceed to add a test spec to assert all database queries return within 500 ms. Click the `Test` tab and proceed to click the `Add Test Spec` button.

In the span selector make sure to add this selector:

```css
span[tracetest.span.type="database"]
```

In the assertion field add:

```css
attr:tracetest.span.duration < 500ms
```

Save the test spec and publish the test.

![Assertion for database queries](https://res.cloudinary.com/djwdcmwdz/image/upload/v1679071121/Blogposts/Docs/screely-1679071115690_hqhzh2.png)

The database spans that are returning in more than 500ms are labeled in red.

![Assertions failing](https://res.cloudinary.com/djwdcmwdz/image/upload/v1679071183/Blogposts/Docs/screely-1679071177655_cjqwlk.png)

This is an example of a trace-based test that asserts against every single part of an HTTP transaction, including all interactions with the database.

However, Tracetest cannot run this test as part of your CI/CD without integrating with another tool.

Let's introduce how Testkube makes it possible.

## 6. Deploy the Tracetest Executor

Testkube works with the concept of Executors. An Executor is a wrapper around a testing framework, Tracetest in this case, in the form of a Docker container and runs as a Kubernetes job. To start you need to register and deploy the Tracetest executor in your cluster using the Testkube CLI.

```bash
kubectl testkube create executor --image kubeshop/testkube-executor-tracetest:latest --types "tracetest/test" --name tracetest-executor

[Output]
Executor created tracetest-executor 🥇
```

## 7. Create a trace-based test with Tracetest in Testkube

Click the ⚙️ button in the top right. Then click `Test Definition`.

![test settings](https://res.cloudinary.com/djwdcmwdz/image/upload/v1679065450/Blogposts/Docs/screely-1679065444972_zzsila.png)

This will open a YAML definition for the test run.

![tracetest test yaml file](https://res.cloudinary.com/djwdcmwdz/image/upload/v1679071430/Blogposts/Docs/screely-1679071422136_ygbo8q.png)

Save this into a file called `test.yaml`:

```yaml
type: Test
spec:
  id: RUkKQ_aVR
  name: Pokeshop - List
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
  - name: Database queries less than 500 ms
    selector: span[tracetest.span.type="database"]
    assertions:
    - attr:tracetest.span.duration  <  500ms
```

Execute the following command to create the test executor object in Testkube. Do not forget to provide the path to your Tracetest definition file using the `--file` argument, and also the Tracetest Server endpoint using the `TRACETEST_ENDPOINT` `--variable`.

```bash
kubectl testkube create test --file location/test.yaml --type "tracetest/test" --name pokeshop-tracetest-test --variable TRACETEST_ENDPOINT=http://tracetest

[Output]
Test created testkube / pokeshop-tracetest-test 🥇
```

Opening the Testkube Dashboard will show the test is created successfully.

![Tracetest test created in Testkube](https://res.cloudinary.com/djwdcmwdz/image/upload/v1679071918/Blogposts/Docs/screely-1679071913649_yrgucd.png)

## 8. Run the Tracetest trace-based test in Testkube

Finally, to run the test, execute the following command, or run the test from the Testkube Dashboard.

```bash
kubectl testkube run test --watch pokeshop-tracetest-test
```

Here's what the Testkube Dashboard will look like if the test fails.
![testkube failing tests](https://res.cloudinary.com/djwdcmwdz/image/upload/v1679322512/Blogposts/Docs/screely-1679322504316_v7qfbf.png)

And, if they pass, it'll look like this.
![testkube passing tests]()

## Infrastructure Overview

The following is high level sequence diagram on how Testkube and Tracetest interact with the different pieces of the system:

```mermaid
sequenceDiagram
    testkube client->>+testkube: Trigger Testkube test run
    testkube->>+executor CRDs: Get executor details
    executor CRDs-->>-testkube: Send details
    testkube->>+tracetest executor job: Schedules execution
    tracetest executor job->>+tracetest executor job: Configure Tracetest CLI
    tracetest executor job->>+tracetest server: Executes the Tracetest test run
    tracetest server->>+instrumented service: Trigger request
    instrumented service-->>-tracetest server: Get response
    instrumented service->>+data store: Send telemetry data
    tracetest server->>+data store: Fetch trace
    data store-->>-tracetest server: Get trace
    tracetest server->>+tracetest server: Run assertions
    tracetest server-->>-tracetest executor job: Return test run results
    tracetest executor job-->>-testkube: Return test run results
    testkube-->>-testkube client: Send details
```
