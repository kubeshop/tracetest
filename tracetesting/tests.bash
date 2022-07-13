#!/bin/bash

source ./funcs.bash

EXIT_STATUS=0

# ensure test not exists
tracetest_target_curl "/api/tests/383d3dce-7b60-4a61-bdea-87f47263af5d" -X DELETE
test "test_create_with_id_notexists" ./definitions/test_create_with_id_notexists.yml || EXIT_STATUS=$?
test "test_create_with_id_exists" ./definitions/test_create_with_id_exists.yml || EXIT_STATUS=$?
tracetest_target_curl "/api/tests/383d3dce-7b60-4a61-bdea-87f47263af5d" -X DELETE

test "test_create" ./definitions/test_create.yml || EXIT_STATUS=$?

export TEST_ID=$(tracetest_target test list | jq -rc '.[0].id')
require_not_empty $TEST_ID "requires TEST_ID, got $TEST_ID " || exit $?

test "test_list" ./definitions/test_list.yml || EXIT_STATUS=$?

test "test_run" ./definitions/test_run.yml || EXIT_STATUS=$?

export RUN_ID=$(tracetest_target_curl "/api/tests/$TEST_ID/run?take=1&skip=0" | jq -rc '.[0].id')
require_not_empty $RUN_ID "requires RUN_ID, got $RUN_ID " || exit $?

test "test_rerun" ./definitions/test_rerun.yml || EXIT_STATUS=$?
test "test_run_delete" ./definitions/test_run_delete.yml || EXIT_STATUS=$?
test "test_delete" ./definitions/test_delete.yml || EXIT_STATUS=$?

exit $EXIT_STATUS
