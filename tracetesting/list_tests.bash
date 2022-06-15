#!/bin/bash

NAME="list_tests"
echo "-> $NAME"
source ./funcs.bash

EXIT_STATUS=0

# setup
export TEST_ID=$(tracetest_target test run --definition ./definitions/pokemon_list.yml | jq -rc '.test.id')

#run test
run_test $NAME "./definitions/tracetest_tests_list.yml"  || EXIT_STATUS=$?

# cleanup
tracetest_target_curl "/api/tests/$TEST_ID" -X DELETE

#exit
exit $EXIT_STATUS
