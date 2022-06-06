#!/bin/bash

set -e

NAMESPACE="tracetest"
TRACE_BACKEND="jaeger"
TRACE_BACKEND_ENDPOINT="jaeger-query:16685"
SKIP_PMA=""
SKIP_JAEGER=""

help_message() {
  echo "Tracetest setup script"
  echo
  echo "options:"
  echo "  --help                                         show help message"
  echo "  --namespace [tracetest]                        target installation k8s namespace"
  echo "  --trace-backend [jaeger]                       trace backend (jaeger or tempo)"
  echo "  --trace-backend-endpoint [jaeger-query:16685]  trace backend endpoint"
  echo "  --skip-pma                                     if set, don't install the sample application"
  echo "  --skip-jaeger                                  if set, don't install jaeger"
  echo
}

crd_exists() {
  kubectl get crd $1 > /dev/null 2>&1
  return
}

install_cert_manager(){
  kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.8.0/cert-manager.yaml
  echo
  echo "--> waiting for cert-manager"
  kubectl wait --for=condition=ready pod -l app=webhook --namespace cert-manager
  echo "--> cert-manager ready. Create self signed ClusterIssuer"
  cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned
spec:
  selfSigned: {}
EOF
  echo "--> ClusterIssuer created"
}

install_jaeger() {
  if crd_exists "jaegers.jaegertracing.io"; then
    echo "--> CRD 'jaegers.jaegertracing.io' already exists. Assume jaeger-operator is installed"
  else
    set +e #ignore errors here
    kubectl create namespace observability
    set -e

    kubectl create -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.32.0/jaeger-operator.yaml -n observability
    echo
    echo "--> waiting for jaeger-operator"
    kubectl wait --for=condition=ready pod -l name=jaeger-operator --namespace observability --timeout 5m
    sleep 5
    echo "--> jaeger-operator ready"
  fi

  echo "--> creating jaeger instance"
  cat <<EOF | kubectl apply --namespace $NAMESPACE -f -
apiVersion: jaegertracing.io/v1
kind: Jaeger
metadata:
  name: jaeger
EOF
   echo "--> jaeger instance created"
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
    --skip-jaeger)
      SKIP_JAEGER="YES"
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


echo "--> install tracetest to namespace $NAMESPACE"
helm upgrade --install tracetest kubeshop/tracetest \
  --namespace $NAMESPACE --create-namespace \
  --set tracingBackend=$TRACE_BACKEND \
  --set ${TRACE_BACKEND}ConnectionConfig.endpoint="$TRACE_BACKEND_ENDPOINT"

if [ "$SKIP_JAEGER" != "YES" ]; then
    echo
    echo
    echo "----------------------------"
    echo "Installing Jaeger"
    echo "----------------------------"
    echo
    echo

    echo "--> install cert-manager"
    if crd_exists "certificates.cert-manager.io"; then
      echo "--> CRD 'certificates.cert-manager.io' already exists. Assume cert-manager is installed"
    else
      install_cert_manager
    fi

    echo
    echo "--> install jaeger"
    install_jaeger
fi


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
      -f values.yaml \
      --set postgres.auth.username=ashketchum,postgres.auth.password=squirtle123,postgres.auth.database=pokeshop \
      --set rabbitmq.auth.username=guest,rabbitmq.auth.password=guest,rabbitmq.auth.erlangCookie=secretcookie \
      --set 'env[5].value=jaeger-agent.'$NAMESPACE'.svc.cluster.local'

    echo "--> waiting for demo app to be ready"
    kubectl wait --for=condition=ready pod --namespace demo demo-rabbitmq-0
    kubectl wait --for=condition=ready pod --namespace demo -l app.kubernetes.io/name=pokemon-api
    echo "--> demo app ready"
    rm -rf $tmpdir
fi


echo "----------------------------"
echo "Installing Tracetest"
echo "----------------------------"
echo

helm repo add kubeshop https://kubeshop.github.io/helm-charts
helm repo update

echo
echo
echo "----------------------------"
echo "Install complete"
echo "----------------------------"
echo
echo "to connect to tracetest, run:"
echo "  kubectl port-forward --namespace $NAMESPACE svc/tracetest 8080:8080"
echo "and navigate to http://localhost:8080"
echo
echo "to see tracetest logs:"
echo "  kubectl logs --namespace $NAMESPACE -f svc/tracetest"

