#/bin/bash

if [ -z "$NAME" ];then
  echo '$NAME is required'
  exit 1
fi

if [ -z "$TAG" ];then
  echo '$TAG is required'
  exit 1
fi

if [ -z "$CONFIG_FILE" ];then
  echo '$CONFIG_FILE is required'
  exit 1
fi

echo $NAME
echo $TAG
echo $CONFIG_FILE
basename $CONFIG_FILE

helm repo add kubeshop https://kubeshop.github.io/helm-charts
helm repo update
helm upgrade --install $NAME kubeshop/tracetest \
  --namespace $NAME --create-namespace \
  --set image.tag=$TAG \
  --set image.pullPolicy=Always \
  --set ingress.enabled=false

kubectl --namespace $NAME create configmap $NAME --from-file $CONFIG_FILE --dry-run=client \
  | envsubst \
  | sed 's#'$(basename $CONFIG_FILE)'#config.yaml#' \
  | kubectl --namespace $NAME replace -f -

kubectl --namespace $NAME delete pods -l app.kubernetes.io/name=tracetest

kubectl --namespace $NAME wait --for=condition=ready pod -l app.kubernetes.io/name=tracetest --timeout 10m
