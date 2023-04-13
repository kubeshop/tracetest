#!/bin/bash
PORT=$1

CONDITION='nc -z -w 1 localhost '$PORT' > /dev/null 2>&1'
IF_TRUE='echo "port '$PORT' ready"'
IF_FALSE='echo "port '$PORT' not available, retry"'

set -ex
bash -c "until ${CONDITION}; do ${IF_FALSE}; sleep 1; done; ${IF_TRUE}"
