#!/bin/sh

set -e

TOKEN=$TRACETEST_API_KEY
ENVIRONMENT_ID=$TRACETEST_ENVIRONMENT_ID

apply() {
  echo "Configuring TraceTest"
  tracetest configure --token $TOKEN --environment $ENVIRONMENT_ID

  echo "Applying Resources"
  tracetest apply datastore -f /resources/datastore.yaml
  tracetest apply test -f /resources/test.yaml
}

apply
