#/bin/bash

tracetest configure -t tttoken_698efda22b86a3be
tracetest run test -f ./test-api.docker.yaml --required-gates test-specs --output pretty
