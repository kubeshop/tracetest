#!/bin/bash

# NEEDS TRACETEST_TOKEN to be set in the environment with organization admin access
# https://docs.tracetest.io/concepts/organization-tokens
TRACETEST_TOKEN=$TRACETEST_TOKEN

# configure tracetest
tracetest configure --token $TRACETEST_TOKEN

# create environment
ENVIRONMENT_ID=$(tracetest apply environment -f environment.yaml --output json | jq -r '.spec.id')
echo "Environment ID: $ENVIRONMENT_ID"

# switching to the environment
tracetest configure --environment $ENVIRONMENT_ID

# start agent
tracetest start --api-key $TRACETEST_TOKEN --environment $ENVIRONMENT_ID
