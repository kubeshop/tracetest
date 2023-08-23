#!/bin/bash
PORT=$1

TIMEOUT=${TIMEOUT:-"30s"}
CONDITION='nc -z -w 1 localhost '$PORT' > /dev/null 2>&1'
IF_TRUE='echo "port '$PORT' ready"'
IF_FALSE='echo "port '$PORT' not available, retry"'

ROOT_DIR=$(cd $(dirname "${BASH_SOURCE:-$0}")/.. && pwd)
$ROOT_DIR/scripts/wait.sh "$TIMEOUT" "$CONDITION" "$IF_TRUE" "$IF_FALSE"
