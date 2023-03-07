#/bin/bash

set -ex

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

extraParams=()

if [ -n "$EXPOSE_HOST" ]; then
  extraParams=("${extraParams[@]}" --set ingress.enabled=true )
  extraParams=("${extraParams[@]}" --set 'ingress.hosts[0].host='$EXPOSE_HOST',ingress.hosts[0].paths[0].path=/,ingress.hosts[0].paths[0].pathType=Prefix' )
fi

if [ -n "$CERT_NAME" ]; then
    extraParams=("${extraParams[@]}" --set ingress.annotations."networking\.gke\.io/managed-certificates"=$CERT_NAME)
    extraParams=("${extraParams[@]}" --set ingress.annotations."networking\.gke\.io/v1beta1\.FrontendConfig"="ssl-redirect")
fi

if [ -n "$BACKEND_CONFIG" ]; then
    extraParams=("${extraParams[@]}" --set service.annotations."cloud\.google\.com/backend-config"='\{\"default\":\"'$BACKEND_CONFIG'\"\}')
fi

helm repo add kubeshop https://kubeshop.github.io/helm-charts
helm repo update
helm upgrade --install $NAME kubeshop/tracetest \
  --namespace $NAME --create-namespace \
  --set image.tag=$TAG \
  --set image.pullPolicy=Always \
  ${extraParams[@]}

PROVISION_FILE=$(cd $(dirname "${BASH_SOURCE:-$0}") && pwd)/provisioning.yaml
kubectl --namespace $NAME create configmap $NAME --from-file=$CONFIG_FILE --from-file=$PROVISION_FILE -o yaml --dry-run=client \
  | envsubst \
  | sed 's#'$(basename $CONFIG_FILE)'#config.yaml#' \
  | kubectl --namespace $NAME replace -f -

kubectl --namespace $NAME delete pods -l app.kubernetes.io/name=tracetest

TIME_OUT=30m
CONDITION='[[ $(kubectl get pods  --namespace '$NAME' -lapp.kubernetes.io/name=tracetest -o jsonpath="{.items[*].status.phase}") = "Running" ]]'
IF_TRUE='echo "pods ready"'
IF_FALSE='echo "pods not ready. retrying"'

ROOT_DIR=$(cd $(dirname "${BASH_SOURCE:-$0}")/.. && pwd)
$ROOT_DIR/scripts/wait.sh "$TIME_OUT" "$CONDITION" "$IF_TRUE" "$IF_FALSE"
