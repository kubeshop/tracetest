#!/bin/bash

set -e

source ./funcs.bash

EXIT_STATUS=0

test "environment_create" ./definitions/environment_create.yml || EXIT_STATUS=$?

export ENV_ID=$(tracetest_target test list -o json | jq -rc '.[0].id')
require_not_empty $ENV_ID "requires ENV_ID, got $ENV_ID " || exit $?

test "environment_list" ./definitions/environment_list.yml || EXIT_STATUS=$?
test "environment_delete" ./definitions/environment_delete.yml || EXIT_STATUS=$?
