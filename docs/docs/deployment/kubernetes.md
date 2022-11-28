# Kubernetes Deployment

You can run Tracetest in a Kubernetes cluster. This setup is ideal for a CI/CD environment, QA teams working on shared environments, etc.
You can use a remote or local (minikube, etc) cluster. We'll even help you setup a local cluster, if you need one.

![Installer using Kubernetes](../img/installer/1_kubernetes_0.7.0.png)

**Tools required (installed if missing)**:
- kubectl
- helm

If you selected to run locally and want the installer to set up [minikube](https://minikube.sigs.k8s.io/docs/) for you:
- Docker

**Requirements**:
- Jaeger or other compatible backend. If missing, the installer will help you configure one.
- OpenTelemetry Collector. If missing, the installer will help you configure one.

**Optionals**:
- [PokeShop demo app](https://github.com/kubeshop/pokeshop/)

**Result**:
- `tracetest` helm chart deployed in the `tracetest` (configurable) namespace.
- [Jaeger](https://www.jaegertracing.io/) instance deployed in the `tracetest` namespace, if selected.
- [Cert Manager](https://cert-manager.io/), if selected.
- [Jaeger Operator](https://www.jaegertracing.io/docs/latest/operator/), if selected.
- [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) deployed in the `tracetest` (configurable) namespace, if selected.
- [PokeShop demo app](https://github.com/kubeshop/pokeshop/) deployed in the `demo` namespace, if selected.

## Access the Tracetest Web UI

Once installed, you can get started by launching Tracetest:

```sh
kubectl port-forward svc/tracetest 11633
```

Then launch a browser to [http://localhost:11633/](http://localhost:11633/).
