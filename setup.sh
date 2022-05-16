#!/bin/bash

NAMESPACE="tracetest"
SKIP_PMA="no"

while [[ $# -gt 0 ]]; do
  case $1 in
    -n|--namespace)
      EXTENSION="$2"
      shift
      shift
      ;;
    --skip-pma)
      SKIP_PMA="YES"
      shift # past argument
      shift # past value
      ;;
    -*|--*)
      echo "Unknown option $1"
      exit 1
      ;;
  esac
done

echo "----------------------------"
echo "Installing Tracetest"
echo "----------------------------"
echo 
echo

helm repo add kubeshop https://kubeshop.github.io/helm-charts
helm repo update

helm install tracetest kubeshop/tracetest \
  --namespace $NAMESPACE \
  --set tracingBackend=jaeger \
  --set jaegerConnectionConfig.endpoint="jaeger-query:16685"


if [ "$SKIP_PMA" = "no" ]; then
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
    helm install -n $NAMESPACE -f values.yaml .

fi
