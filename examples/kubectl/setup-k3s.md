# Test Kubectl commands

This example shows how one can use Tracetest to do trace based testing for kubectl commands. This works by leveraging the `Audit-Id` header sent back from the kubernetes api server for every request.

the `kubectl` command allows to enable a very verbose output, where it includes the response headers. This example script parses that info and maps it to the corresponding Trace ID.

## Requirements

This setup is a bit limited at the moment. It only works with Jaeger backend, becasue it provides a search API that allows us to search spans that include the correct `Audit-Id` tag.

For everything to work, we need the following environment configured:

1. A Kubernetes cluster with `APIServerTracing` enabled and configured
2. A OtelCollector that can be accessed by the k8s api server to send its tracing data
3. A Jager backend to store the traces


## Setup

k3s can be used to easily setup a cluster with the correct configuration. It's easy enought to be configured in any dev environment, even a CICD.

The otelcollector cannot be deployed in the target cluster because it complicates the networking setup needed to have the k8s apiserver communicating with the collector.
An easy solution for this is to have docker compose starting the otelcol/jaeger services, so their ports can be exposed directly to the host network.

### 1. Setup OtelCollector/Jaeger

create a `docker-compose.yaml` file that includes both services:

```yaml
# docker-compose.yaml
services:
  jaeger:
      healthcheck:
          test:
              - CMD
              - wget
              - --spider
              - localhost:16686
          timeout: 3s
          interval: 1s
          retries: 60
      image: jaegertracing/all-in-one:latest
      restart: unless-stopped
      ports:
        - 16686:16686
  otel-collector:
      command:
          - --config
          - /otel-local-config.yaml
      depends_on:
          jaeger:
              condition: service_started
      image: otel/opentelemetry-collector:0.54.0
      ports:
        - 4317:4317
      volumes:
          - ./otel-collector.yaml:/otel-local-config.yaml
```

We need to create a config file in the same directory for the otel-collector so it can communicate with Jaeger:

```yaml
# otel-collector.yaml
receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  # Data sources: traces
  probabilistic_sampler:
    hash_seed: 22
    sampling_percentage: 100

  batch:
    timeout: 100ms

exporters:
  # logging is optional, but useful for making sure the collector is receiving traces
  logging:
    logLevel: debug
  jaeger:
    # this url is valid within the `docker-compose` environmet this collector is running
    endpoint: jaeger:14250 
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [probabilistic_sampler, batch]
      exporters: [jaeger,logging] # logging is optional,
```

With both files created on the same director, we can start them:

```
sudo docker compose up -d
```

### 2. Start a k3s cluster

We need to start a k3s cluster, but we need some special configurations:

1. the k8s version needs to be 1.27+ for traces to work
2. we need to enable the tracing feature flag
3. we need to pass the otel-collector address to the apiserver


First, we'll create the tracing config file, because it's required for the apiserver to start.

> Make sure that the directory exists by running: `sudo mkdir -p /etc/kube-tracing/`

```yaml
# /etc/kube-tracing/apiserver-tracing.yaml
apiVersion: apiserver.config.k8s.io/v1beta1
kind: TracingConfiguration
# default value
endpoint: localhost:4317
samplingRatePerMillion: 1000000 # 100%
EOF
```

This settings are documented [here](https://kubernetes.io/docs/concepts/cluster-administration/system-traces/). The only difference is that we are setting the `samplingRatePerMillion` value to 1.000.000 (meaning, 100%) so that all traces are send to the collector.
This effectively disables sampling. Use this setting with care, and you probablly shouldn't do this in a prod environment.

Now we can install the k3s cluster:

```sh
curl -sfL https://get.k3s.io | \
  INSTALL_K3S_VERSION="v1.27.1+k3s1" \
  INSTALL_K3S_EXEC="--kube-apiserver-arg=feature-gates=APIServerTracing=true --kube-apiserver-arg=tracing-config-file=/etc/kube-tracing/apiserver-tracing.yaml" \
   sh -s - server --cluster-init
```

We use the `INSTALL_K3S_VERSION` env var to set the correct version, and the `INSTALL_K3S_EXEC` to pass flags and settings to the k8s apiserver.

## Running tests

That's it! We have a k8s apiserver ready to do trace based testing. With our [sample script](./test-kubectl.bash) it's very easy.

You only need to set the correct jaeger UI url and the command you want to run. For example:

```sh
JAEGER_UI_URL="http://127.0.0.1:16686" ./test-kubectl.bash "kubectl get namespaces"
```

