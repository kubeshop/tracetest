#!/bin/sh

set -e

TOKEN=$TRACETEST_TOKEN
ENVIRONMENT_ID=$TRACETEST_ENVIRONMENT_ID
GRAFANA_HASH=$GRAFANA_AUTH_READ_HASH

apply() {
  echo "Configuring TraceTest"
  tracetest configure --token $TOKEN --environment $ENVIRONMENT_ID

  echo "Applying Resources"
  echo "
type: DataStore
spec:
  id: current
  name: Grafana Tempo Cloud
  type: tempo
  tempo:
    type: http
    http:
      url: https://tempo-us-central1.grafana.net/tempo
      headers:
        Authorization: Basic ${GRAFANA_HASH}
" > /resources/datastore.yaml

  tracetest apply datastore -f /resources/datastore.yaml
  tracetest apply test -f /resources/test.yaml
}

apply
