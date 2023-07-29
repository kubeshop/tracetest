# Installing Tracetest and Signoz on a Kubernetes cluster

```sh
# create test cluster
k3d cluster create tracetest-signoz
```

```sh
helm repo add signoz https://charts.signoz.io
helm repo add kubeshop https://kubeshop.github.io/helm-charts
helm repo update
``````

```sh
# install signoz
helm install signoz signoz/signoz --namespace observability --create-namespace
```

```sh
# create local configuration for collector
cat << EOF > opentelemetry-collector-resources.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: collector-config
data:
  collector.yaml: |
    receivers:
      otlp:
        protocols:
          grpc:
          http:

    processors:
      batch:
        timeout: 100ms

      # Data sources: traces
      probabilistic_sampler:
        hash_seed: 22
        sampling_percentage: 100

    exporters:
      # OTLP for Tracetest
      otlp/tracetest:
        endpoint: tracetest.tracetest.svc.cluster.local:4317
        tls:
          insecure: true
      # OTLP for Signoz
      otlp/signoz:
        endpoint: signoz-otel-collector.observability.svc.cluster.local:4317
        tls:
          insecure: true

    service:
      pipelines:
        traces:
          receivers: [otlp]
          processors: [probabilistic_sampler, batch]
          exporters: [otlp/signoz,otlp/tracetest]
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: otel-collector
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: otel-collector
  template:
    metadata:
      labels:
        app.kubernetes.io/name: otel-collector
    spec:
      containers:
        - name: otelcol
          args:
            - --config=/conf/collector.yaml
          image: otel/opentelemetry-collector-contrib:0.67.0
          volumeMounts:
            - mountPath: /conf
              name: collector-config
          resources:
            requests:
              cpu: 250m
              memory: 512Mi
      volumes:
        - configMap:
            items:
              - key: collector.yaml
                path: collector.yaml
            name: collector-config
          name: collector-config
---
apiVersion: v1
kind: Service
metadata:
  name: otel-collector
spec:
  ports:
    - name: grpc-otlp
      port: 4317
      protocol: TCP
      targetPort: 4317
  selector:
    app.kubernetes.io/name: otel-collector
  type: ClusterIP
EOF
```

```sh
kubectl create namespace tracetest
# installing opentelemetry-collector
kubectl -n tracetest apply -f opentelemetry-collector-resources.yaml
```

```sh
cat << EOF > tracetest-values.yaml
provisioning: |
  type: DataStore
  spec:
    id: current
    name: Signoz
    type: signoz

env:
  tracetestDev: "true"

telemetry:
  exporters:
    collector:
      exporter:
        type: collector
        collector:
          endpoint: otel-collector.tracetest.svc.cluster.local:4317
EOF
```

```sh
# install tracetest
helm install tracetest kubeshop/tracetest --namespace tracetest --create-namespace --values tracetest-values.yaml
```

# open port-forward in another console session
kubectl port-forward --namespace tracetest svc/tracetest 11633 & \
kubectl port-forward --namespace observability svc/signoz-frontend 3301:3301 & \
echo "Press CTRL-C to stop port forwarding and exit the script"
wait


# run tracetest version to check if everything is ok
tracetest version
# it should return an output like:
# CLI: v0.13.0
# Server: v0.13.0
# ✔️ Version match

# update tracetest data store to use signoz
tracetest apply datastore -f ./tracetest-datastore.yaml



# montar comando de teste para garantir que está tudo ok
