#!/bin/sh

set -e

TOKEN=$TRACETEST_TOKEN
ENVIRONMENT_ID=$TRACETEST_ENVIRONMENT_ID
ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID
SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY
SESSION_TOKEN=$AWS_SESSION_TOKEN
REGION=$AWS_REGION

apply() {
  echo "Configuring TraceTest"
  tracetest configure --token $TOKEN --environment $ENVIRONMENT_ID

  echo "Applying Resources"
  tracetest apply datastore -f /resources/datastore.yaml
  tracetest apply test -f /resources/test.yaml
}

apply
