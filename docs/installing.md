# Installation

During the setup, we'll deploy Tracetest and Postgres with Helm.

For the architectural overview of the components, please check the [Architecture](architecture.md) page.

## **Prerequsities**

### **Installation Requirements**

Tools needed for the installation:

- [Helm v3](https://helm.sh/docs/intro/install/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/)

## **Installation**

### Install script

We provide a simple install script that can install all required components:

```
curl -L https://raw.githubusercontent.com/kubeshop/main/setup.sh | bash -s
```

This command will install Tracetest using the default settings. You can configure the following options:


| Option                   | description                                  | Default value      |
|--------------------------|----------------------------------------------|--------------------|
| --help                   | show help message                            | n/a                |
| --namespace              | target installation k8s namespace            | tracetest          |
| --trace-backend          | trace backend (jaeger or tempo)              | jaeger             |
| --trace-backend-endpoint | trace backend endpoint                       | jaeger-query:16685 |
| --skip-pma               | if set, don't install the sample application | n/a                |
| --skip-backend           | if set, don't install the jaeger backend     | n/a                |


Example with custom options:

```
curl -L https://raw.githubusercontent.com/kubeshop/main/setup.sh | bash -s -- --skip-pma --namespace my-custom-namespace
```

### Using Helm

Container images are hosted on the Docker Hub [Tracetest repository](https://hub.docker.com/r/kubeshop/tracetest).

Tracetest currently supports two traces backend: Jaeger and Grafana Tempo.

#### **Jaeger**

Tracetest uses [Jaeger Query Service `16685` port](https://www.jaegertracing.io/docs/1.32/deployment/#query-service--ui) to find Traces using gRPC protocol.

The commands below will install Tracetest connecting to the Jaeger tracing backend on `jaeger-query:16685`.

```sh
# Install Kubeshop Helm repo and update it
helm repo add kubeshop https://kubeshop.github.io/helm-charts
helm repo update

helm install tracetest kubeshop/tracetest \
  --set tracingBackend=jaeger \
  --set jaegerConnectionConfig.endpoint="jaeger-query:16685" # update this value to point to your jaeger install
```

#### **Grafana Tempo**

Tracetest uses [Grafana Tempo's Server's `9095` port](https://grafana.com/docs/tempo/latest/configuration/#server) to find Traces using gRPC protocol.


The commands below will install the Tracetest application connecting to Grafana Tempo tracing backend on `grafana-tempo:9095`:

```sh
# Install Kubeshop Helm repo and update it
helm repo add kubeshop https://kubeshop.github.io/helm-charts
helm repo update

helm install tracetest kubeshop/tracetest \
  --set tracingBackend=tempo \
  --set tempoConnectionConfig.endpoint="grafana-tempo:9095"  # update this value to point to your tempo install
```

## **Uninstallation**

The following command will uninstall Tracetest with Postgres:

```sh
helm delete tracetest
```
