#/bin/bash

# Add a Tracetest token here
# https://docs.tracetest.io/concepts/environment-tokens
tracetest configure -t tttoken_<token>
tracetest run test -f ./api.pokemon.spec.docker.yaml --required-gates test-specs --output pretty
