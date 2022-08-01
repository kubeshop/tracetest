#!/bin/bash

source ./funcs.bash

EXIT_STATUS=0


test "grpc_test_create" ./definitions/grpc_test_create.yml || EXIT_STATUS=$?
export TEST_ID=$(tracetest_target test list | jq -rc '.[0].id')
require_not_empty $TEST_ID "requires TEST_ID, got $TEST_ID " || exit $?
test "grpc_test_run" ./definitions/grpc_test_run.yml || EXIT_STATUS=$?


test "grpc_test_create_invalid_metadata" ./definitions/grpc_test_create_invalid_metadata.yml || EXIT_STATUS=$?
export TEST_ID=$(tracetest_target test list | jq -rc '.[0].id')
require_not_empty $TEST_ID "requires TEST_ID, got $TEST_ID " || exit $?
test "grpc_test_run_invalid_metadata" ./definitions/grpc_test_run_invalid_metadata.yml || EXIT_STATUS=$?

exit $EXIT_STATUS
