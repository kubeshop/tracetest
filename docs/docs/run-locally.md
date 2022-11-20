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

Now that Tracetest is running, we can expose the dependencies in our cluster to the host machine so they are accessible to the development build.
Tracetests needs postgres to store the tests, results, etc, and access to the trace backend (jaeger, tempo, etc) to fetch traces.
We can use kubectl's port forwarding capabilites for this

```
(trap "kill 0" SIGINT; kubectl port-forward svc/tracetest-postgresql 5432:5432 & kubectl port-forward svc/jaeger-query 16685:16685)
```

### **Start Development Server**

When running the development version, the frontend and backend are built and run separately. You need to have both services running to access the tool.

To start the backend server:

```
make server-run
```

To start the frontend server:
```
cd web
npm install -d
npm start
```

The Tracetest development build is available at [http://localhost:3000]. Note that the port is now `3000` since we are accessing the node development server.
