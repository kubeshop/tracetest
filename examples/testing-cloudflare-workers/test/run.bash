#/bin/bash

tracetest configure -t <YOUR_TRACETEST_API_TOKEN> # add your token here
tracetest run test -f ./test-api.docker.yaml --required-gates test-specs --output pretty
# Add more tests here! :D