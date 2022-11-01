#!/bin/bash

set -e

source ./funcs.bash

EXIT_STATUS=0

export TEST_ID=$(tracetest_target test list -o json | jq -rc '.[0].id')
require_not_empty $TEST_ID "requires TEST_ID, got $TEST_ID " || exit $?

test "transaction_create" ./definitions/transaction_create.yml || EXIT_STATUS=$?

export TRANSACTION_ID=$(tracetest_target test list -o json | jq -rc '.[0].id')
require_not_empty $TRANSACTION_ID "requires TRANSACTION_ID, got $TRANSACTION_ID " || exit $?

test "transaction_list" ./definitions/transaction_list.yml || EXIT_STATUS=$?
test "transaction_delete" ./definitions/environment_delete.yml || EXIT_STATUS=$?

exit $EXIT_STATUS
