#!/bin/bash

NAME="list_tests"
echo "-> $NAME"
source ./funcs.bash

EXIT_STATUS=0

# setup
export OUT=$(tracetest_target test run --definition ./definitions/pokemon_list.yml)
echo $OUT
TEST_ID=$(echo $OUT | jq -rc '.test.id')
echo $TEST_ID

#run test
run_test $NAME "./definitions/tracetest_tests_list.yml"  || EXIT_STATUS=$?

# cleanup
tracetest_target_curl "/api/tests/$TEST_ID" -X DELETE

#exit
exit $EXIT_STATUS
