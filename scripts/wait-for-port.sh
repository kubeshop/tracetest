#!/bin/bash
PORT=$1
timeout 30s bash -c 'until nc -z -w 1 localhost '$PORT' > /dev/null 2>&1; do echo "port '$PORT' not available, retry"; sleep 1; done; echo "port '$PORT' ready"'
