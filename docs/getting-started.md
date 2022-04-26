# Getting Started

Please follow the [install steps](/docs/installing.md) for the first installation of Tracetest.

Once installed, you can get started by launching the Tracetest Dashboard by following these instructions:

Run:

kubectl config set-context --current --namespace=tracetest

kubectl port-forward svc/tracetest 8080

Then launch a browser to [http://localhost:8080/](http://localhost:8080/).

To learn how to create your first test, see [create a test](/docs/create-test.md).
