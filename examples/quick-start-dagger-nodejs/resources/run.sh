#!/bin/sh

set -e

TOKEN=$TRACETEST_TOKEN
ENVIRONMENT_ID=$TRACETEST_ENVIRONMENT_ID

run() {
  echo "Configuring Tracetest"
  tracetest configure --token $TOKEN --environment $ENVIRONMENT_ID

  echo "Running Trace-Based Tests..."
  tracetest run test -f /resources/test.yaml
}

run
