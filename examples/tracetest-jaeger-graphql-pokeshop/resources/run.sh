#!/bin/sh

set -e

TOKEN=$TRACETEST_TOKEN
ENVIRONMENT_ID=$TRACETEST_ENVIRONMENT_ID
SERVER_URL=$SERVER_URL

run() {
  echo "Configuring Tracetest"
  tracetest configure --token $TOKEN --environment $ENVIRONMENT_ID --server-url $SERVER_URL

  echo "Running Trace-Based Tests..."
  tracetest run test -f /resources/test.yaml
}

run
