#!/bin/sh

set -e

TOKEN=$TRACETEST_TOKEN
ENVIRONMENT_ID=$TRACETEST_ENVIRONMENT_ID
SERVER_URL=$SERVER_URL

apply() {
  echo "Configuring TraceTest"
  tracetest configure --token $TOKEN --environment $ENVIRONMENT_ID --server-url $SERVER_URL

  echo "Applying Resources"
  tracetest apply datastore -f /resources/datastore.yaml
  tracetest apply test -f /resources/test.yaml
}

apply
