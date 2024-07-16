#!/bin/sh

set -e

TOKEN=$TRACETEST_API_KEY
ENVIRONMENT_ID=$TRACETEST_ENVIRONMENT_ID

run() {
  echo "Configuring Tracetest"
  tracetest configure --token $TOKEN --environment $ENVIRONMENT_ID

  echo "Running Trace-Based Tests..."
  tracetest run test -f /resources/import-pokemon.yaml
}

run
