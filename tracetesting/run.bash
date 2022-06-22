#/bin/bash

export TRACETEST_CMD=${TRACETEST_CLI:-"../cli/tracetest"}
if ! command -v "$TRACETEST_CLI" &> /dev/null; then
  echo "\$TRACETEST_CLI not set to executable. set to $TRACETEST_CLI";
  exit 2
fi

export TARGET_URL=${TARGET_URL:-"http://localhost:8081"}
if [  "$TARGET_URL" = "" ]; then
  echo "\$TARGET_URL not set";
  exit 2
fi

echo "TRACETEST_CLI: $TRACETEST_CLI"
echo "TARGET_URL: $TARGET_URL"

export TRACETEST_MAIN_ENDPOINT="localhost:8080"
export TRACETEST_TARGET_ENDPOINT="localhost:8081"

export DEMO_APP_URL=${DEMO_APP_URL-"http://demo-pokemon-api.demo.svc.cluster.local"}

mkdir -p results/responses

EXIT_STATUS=0
bash ./tests.bash || EXIT_STATUS=$?

exit $EXIT_STATUS

