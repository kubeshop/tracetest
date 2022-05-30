#!/bin/bash

NAMESPACE="tracetest"
TRACE_BACKEND="jaeger"
TRACE_BACKEND_ENDPOINT="jaeger-query:16685"
SKIP_PMA=""

help_message() {
  echo "Tracetest setup script"
  echo
  echo "options:"
  echo "  --help                                         show help message"
  echo "  --namespace [tracetest]                        target installation k8s namespace"
  echo "  --trace-backend [jaeger]                       trace backend (jaeger or tempo)"
  echo "  --trace-backend-endpoint [jaeger-query:16685]  trace backend endpoint"
  echo "  --skip-pma                                     if set, don't install the sample application"
  echo
}

while [[ $# -gt 0 ]]; do
  case $1 in
    --namespace)
      NAMESPACE="$2"
      shift
      shift
      ;;
    --trace-backend)
      TRACE_BACKEND="$2"
      # to lowercase
      TRACE_BACKEND=$(echo "$TRACE_BACKEND" | tr '[:upper:]' '[:lower:]')

      shift
      shift
      ;;
    --trace-backend-endpoint)
      TRACE_BACKEND_ENDPOINT="$2"
      shift
      shift
      ;;
    --skip-pma)
      SKIP_PMA="YES"
      shift # past argument
      ;;
    -h|--help)
      help_message
      exit
      ;;
    -*|--*)
      echo "Unknown option $1"
      help_message
      exit 1
      ;;
  esac
done

echo "----------------------------"
echo "Installing Tracetest"
echo "----------------------------"
echo

helm repo add kubeshop https://kubeshop.github.io/helm-charts
helm repo update

helm upgrade --install tracetest kubeshop/tracetest \
  --namespace $NAMESPACE --create-namespace \
  --set tracingBackend=$TRACE_BACKEND \
  --set ${TRACE_BACKEND}ConnectionConfig.endpoint="$TRACE_BACKEND_ENDPOINT"


if [ "$SKIP_PMA" != "YES" ]; then
    echo
    echo
    echo "----------------------------"
    echo "Installing PokeShop"
    echo "----------------------------"
    echo
    echo
    tmpdir=`mktemp -d`
    curl -L https://github.com/kubeshop/pokeshop/tarball/master | tar -xz --strip-components 1 -C  $tmpdir
    cd $tmpdir/helm-chart
    helm dependency update
    helm upgrade --install demo . \
      --namespace demo --create-namespace \
      -f values.yaml

fi


echo "----------------------------"
echo "Install complete"
echo "----------------------------"
echo
echo "to connect to tracetest, run:"
echo "  kubectl port-forward --namespace ${NAMESPACE} svc/tracetest 8080:8080"
echo "and navigate to http://localhost:8080"
echo
echo "to see tracetest logs: kubectl logs --namespace ${NAMESPACE} -f tracetest svc/tracetest"

