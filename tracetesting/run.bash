#/bin/bash

export TRACETEST_CMD=${TRACETEST_CLI:-"../cli/tracetest"}
if ! command -v "$TRACETEST_CLI" &> /dev/null; then
  echo "\$TRACETEST_CLI not set to executable. set to $TRACETEST_CLI";
  exit 2
fi

if [  "$TARGET_URL" = "" ]; then
  echo "\$TARGET_URL not set";
  exit 2
fi

export TRACETEST_MAIN_ENDPOINT="localhost:8080"
export TRACETEST_TARGET_ENDPOINT="localhost:8081"
export TARGET_URL=${TARGET_URL:-"http://localhost:8081"}

mkdir -p results/responses

EXIT_STATUS=0
bash ./tests.bash || EXIT_STATUS=$?

exit $EXIT_STATUS

