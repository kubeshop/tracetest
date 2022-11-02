#!/bin/bash

set -e

source ./funcs.bash

EXIT_STATUS=0

# create a test to register it as a step
test "transaction_step_create" ./definitions/test_create.yml || EXIT_STATUS=$?

export TEST_ID=$(tracetest_target test list -o json | jq -rc '.[0].id')
require_not_empty $TEST_ID "requires TEST_ID, got $TEST_ID " || exit $?

test "transaction_create" ./definitions/transaction_create.yml || EXIT_STATUS=$?

export TRANSACTION_ID=$(tracetest_target_curl "/api/transactions" -X GET | jq -rc '.[0].id')
require_not_empty $TRANSACTION_ID "requires TRANSACTION_ID, got $TRANSACTION_ID " || exit $?

test "transaction_list" ./definitions/transaction_list.yml || EXIT_STATUS=$?
test "transaction_delete" ./definitions/transaction_delete.yml || EXIT_STATUS=$?

exit $EXIT_STATUS
