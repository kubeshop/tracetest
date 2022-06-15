#!/bin/bash

source ./funcs.bash

EXIT_STATUS=0

# setup
export TEST_ID=$(tracetest_target test run --definition ./definitions/pokemon_list.yml | jq -rc '.test.id')

#run test
run_test "list_tests" "./definitions/tracetest_tests_list.yml"  || EXIT_STATUS=$?

#exit
exit $EXIT_STATUS
