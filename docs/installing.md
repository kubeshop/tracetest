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
curl -L https://raw.githubusercontent.com/kubeshop/tracetest/main/setup.sh | bash -s
```

This command will install Tracetest using the default settings. You can configure the following options:

| Option                   | description                                  | Default value      |
| ------------------------ | -------------------------------------------- | ------------------ |
| --help                   | show help message                            | n/a                |
| --namespace              | target installation k8s namespace            | tracetest          |
| --trace-backend          | trace backend (jaeger or tempo)              | jaeger             |
| --trace-backend-endpoint | trace backend endpoint                       | jaeger-query:16685 |
| --skip-collector         | if set, don't install the otel-collector     | n/a                |
| --skip-pma               | if set, don't install the sample application | n/a                |
| --skip-backend           | if set, don't install the jaeger backend     | n/a                |

Example with custom options:

```
curl -L https://raw.githubusercontent.com/kubeshop/tracetest/main/setup.sh | bash -s -- --skip-pma --namespace my-custom-namespace
```

### **Using Helm**

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
  --set telemetry.dataStores.jaeger.jaeger.endpoint="jaeger-query:16685" \ # update this value to point to your jaeger install
  --set telemetry.exporters.collector.exporter.collector.endpoint="otel-collector:4317" \ # update this value to point to your collector install
  --set server.telemetry.dataStore="jaeger"
```

#### **Grafana Tempo**

Tracetest uses [Grafana Tempo's Server's `9095` port](https://grafana.com/docs/tempo/latest/configuration/#server) to find Traces using gRPC protocol.

The commands below will install the Tracetest application connecting to the Grafana Tempo tracing backend on `grafana-tempo:9095`:

```sh
# Install Kubeshop Helm repo and update it
helm repo add kubeshop https://kubeshop.github.io/helm-charts
helm repo update

helm install tracetest kubeshop/tracetest \
  --set telemetry.dataStores.tempo.tempo.endpoint="grafana-tempo:9095" \ # update this value to point to your tempo install
  --set telemetry.exporters.collector.exporter.collector.endpoint="otel-collector:4317" \ # update this value to point to your collector install
  --set server.telemetry.dataStore="tempo"
```

#### **Opensearch**

The commands below will install the Tracetest application connecting to the Opensearchtracing backend on `opensearch:9200`:

```sh
# Install Kubeshop Helm repo and update it
helm repo add kubeshop https://kubeshop.github.io/helm-charts
helm repo update

helm install tracetest kubeshop/tracetest \
  --set telemetry.dataStores.opensearch.opensearch.addresses={"http://opensearch:9200"} \ # update this value to point to your opensearch install
  --set telemetry.dataStores.opensearch.opensearch.index="traces" \ # update this value to use the index where your traces are being stored
  --set telemetry.dataStores.opensearch.opensearch.username="admin" \ # update this value with your opensearch username
  --set telemetry.dataStores.opensearch.opensearch.password="admin" \ # update this value with your opensearch password
  --set server.telemetry.dataStore="opensearch"
```

#### **SignalFX**

The commands below will install the Tracetest application connecting to the Opensearchtracing backend on `opensearch:9200`:

```sh
# Install Kubeshop Helm repo and update it
helm repo add kubeshop https://kubeshop.github.io/helm-charts
helm repo update

helm install tracetest kubeshop/tracetest \
  --set telemetry.dataStores.signalfx.signalfx.token="your signalfx token" \ # update this value to point to your signalfx account
  --set telemetry.dataStores.signalfx.signalfx.realm="your realm (us1?)" \ # update this value to point to your signalfx account
  --set telemetry.dataStores.signalfx.signalfx.url="your custom url" \ # update this value to point to your signalfx custom url. This is optional.
  --set server.telemetry.dataStore="signalfx"
```


### **Have a different backend trace data store?**

[Tell us](https://github.com/kubeshop/tracetest/issues/new?assignees=&labels=&template=feature_request.md&title=) which one you have and we will see if we can add support for it!

## **Uninstallation**

The following command will uninstall Tracetest with Postgres:

```sh
helm delete tracetest
```

## CLI Installation
Every time we release a new version of Tracetest, we generate binaries for Linux, MacOS, and Windows. Supporting both amd64, and ARM64 architectures, in `tar.gz`, `deb`, `rpm` and `exe` formats
You can find the latest version [here](https://github.com/kubeshop/tracetest/releases/latest).

### Linux/MacOS

Tracetest CLI can be installed automatically using the following script:
```sh
curl -L https://raw.githubusercontent.com/kubeshop/tracetest/main/install-cli.sh | sh
```

It works for systems with Hombrew, `apt-get`, `dpkg`, `yum`, `rpm` installed, and if no package manager is available, it will try to download the build and install it manually.

You can also manually install with one of the following methods

#### Homebrew

```sh
brew install kubeshop/tracetest/tracetest
```

#### apt

```sh
# requirements for our deb repo to work
sudo apt-get update && sudo apt-get install -y apt-transport-https ca-certificates

# add repo
echo "deb [trusted=yes] https://apt.fury.io/tracetest/ /" | sudo tee /etc/apt/sources.list.d/fury.list

# update and install
sudo apt-get update
sudo apt-get install tracetest
```

#### yum

```sh
# add repository
cat <<EOF | $SUDO tee /etc/yum.repos.d/tracetest.repo
[tracetest]
name=Tracetest
baseurl=https://yum.fury.io/tracetest/
enabled=1
gpgcheck=0
EOF

# install
sudo yum install tracetest --refresh
```

### Windows
Download one of the files from the latest tag, extract to your machine, and then [add the tracetest binary to your PATH variable](https://stackoverflow.com/a/41895179)
