#/bin/bash

tracetest configure -t tttoken_<token> # Add your token
tracetest run test -f ./test-api.docker.yaml --required-gates test-specs --output pretty
