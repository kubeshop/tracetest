# Tracetest

This is the Helm chart for [Tracetest](https://github.com/kubeshop/tracetest) installation.

## Installation

### Chart installation

Add repo:

```sh
helm repo add kubeshop https://kubeshop.github.io/helm-charts
helm repo update

```

```sh
helm install tracetest kubeshop/tracetest
```

## Uninstall

```sh
helm delete tracetest
```
