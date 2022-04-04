# Installation

During the setup we'll deploy Tracetest, and Postgres with Helm.

For the architectural overview of the components please check the [Architecture](architecture.md) page.

## Prerequsities

### Installation requirements

Tools needed for the installation:

- [Helm v3](https://helm.sh/docs/intro/install/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/)

## Installation

Container images are hosted on Docker Hub [Tracetest repository](https://hub.docker.com/r/kubeshop/tracetest).

There are two options to install Tracetest, if you use Jaeger to store your trace
### Jaeger

Tracetest uses [Jaeger Query Service `16685` port](https://www.jaegertracing.io/docs/1.32/deployment/#query-service--ui), which allows Tracetest to find Traces using GRPC protocol.

The commands below will install Tracetest application connecting to Jaeger tracing backend on `jaeger-query:16685`.

```sh
# Install Kubeshop Helm repo and update it
helm repo add kubeshop https://kubeshop.github.io/helm-charts
helm repo update

helm install tracetest kubeshop/tracetest \
  --set tracingBackend=jaeger \
  --set jaegerConnectionConfig.endpoint="jaeger-query:16685"
```

### Grafana Tempo

Tracetest uses [Grafana Tempo's Server's `9095` port](https://grafana.com/docs/tempo/latest/configuration/#server), which allows Tracetest to find Traces using GRPC protocol.


The commands below will install Tracetest application connecting to Grafana Tempo tracing backend on `grafana-tempo:9095`:

```sh
# Install Kubeshop Helm repo and update it
helm repo add kubeshop https://kubeshop.github.io/helm-charts
helm repo update

helm install tracetest kubeshop/tracetest \
  --set tracingBackend=tempo \ 
  --set tempoConnectionConfig.endpoint="grafana-tempo:9095"
```

## Uninstallation

The following command will uninstall Tracetest with Postgres:

```sh
# Delete releases

helm delete tracetest
```
