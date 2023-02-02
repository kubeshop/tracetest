# Run Tracetest Locally

Tracetest depends on a postgres database and a trace store backend (Jaeger or Tempo). The frontend requires node and npm and the backend requires the go tooling.

## **Run on Local Kubernetes**

Tracetest and its dependencies can be installed in a local Kubernetes cluster (microk8s, minikube, Kubernetes for Docker Desktop, etc).
Following the [install steps](./getting-started/installation) will install a running instance of Tracetest and Postgres. Installing Jaeger is the easiest way to get a trace store backend.

The Tracetest install can be exposed with a `LoadBalancer`, `NodePort` or any similar mechanism. It can also be kept internally, only expose the Jaeger and postgres port,
and use them to run local development builds. This is useful to quickly test changes on both the front and back end.

### **Installing Jaeger**

Before installing Tracetest, we need to setup the [Jaeger operator](https://www.jaegertracing.io/docs/1.32/operator/), which in turn has a dependency on [cert-mnanager](https://cert-manager.io/).
cert-manager has different [install options](https://cert-manager.io/docs/installation/). The simplest way to do this is to use a static install:

```
$ kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.8.0/cert-manager.yaml
```

Once the pods in `cert-manager` namespace are running, we can install the Jaeger operator:

```
kubectl create namespace observability
kubectl create -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.32.0/jaeger-operator.yaml -n observability
```

Now, create an `AllInOne` jaeger instance:

```
cat <<EOF | kubectl create -f -
apiVersion: jaegertracing.io/v1
kind: Jaeger
metadata:
  name: jaeger
EOF
```

### **Install Tracetest**

Follow the [install steps](./getting-started/installation):

```sh
helm repo add kubeshop https://kubeshop.github.io/helm-charts
helm repo update

helm install tracetest kubeshop/tracetest \
  --set telemetry.dataStores.jaeger.jaeger.endpoint="jaeger-query:16685" \ # update this value to point to your jaeger install
  --set telemetry.exporters.collector.exporter.collector.endpoint="otel-collector:4317" \ # update this value to point to your collector install
  --set server.telemetry.dataStore="jaeger"
```

You can now expose the Tracetest service using a `LoadBalancer`, `NodePort` or even a simple `port-forward`:

```
kubectl port-forward svc/tracetest 11633:11633
```

Now Tracetest is available at [http://localhost:11633]

## **Run a Development Build**

You can build and run Tracetest locally using Docker compose. The project provides a Docker compose file and a Makefile with targets to generate the required image.

You will need [Docker Compose](https://docs.docker.com/compose/install/) installed, as well as [GoReleaser-Pro](https://goreleaser.com/install/). 
Note that while the `pro` version is required, no licencing is needed for local builds.

To build the image:

```bash
make build-docker
```

This will build a new image, tagged `kubeshop/tracetest:latest`. That is the default image used in the `docker-compose` file, so now we can start the service:

```bash
docker compose up
```

Once the app is finished loading all the services, you can access your local tracetest at [http://localhost:11633].
