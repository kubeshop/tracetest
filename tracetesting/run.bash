#/bin/bash

export TRACETEST_CLI=${TRACETEST_CLI:-"../cli/tracetest"}
if ! command -v "$TRACETEST_CLI" &> /dev/null; then
  echo "\$TRACETEST_CLI not set to executable. set to $TRACETEST_CLI";
  exit 2
fi

export TARGET_URL=${TARGET_URL:-"http://localhost:8081"}
if [  "$TARGET_URL" = "" ]; then
  echo "\$TARGET_URL not set";
  exit 2
fi


export TRACETEST_MAIN_ENDPOINT=${TRACETEST_MAIN_ENDPOINT:-"localhost:8080"}
export TRACETEST_TARGET_ENDPOINT=${TRACETEST_TARGET_ENDPOINT:-"localhost:8081"}
export DEMO_APP_URL=${DEMO_APP_URL-"http://demo-pokemon-api.demo.svc.cluster.local"}

echo "TRACETEST_CLI: $TRACETEST_CLI"
echo "TARGET_URL: $TARGET_URL"
echo "TRACETEST_MAIN_ENDPOINT: $TRACETEST_MAIN_ENDPOINT"
echo "TRACETEST_TARGET_ENDPOINT: $TRACETEST_TARGET_ENDPOINT"
echo "DEMO_APP_URL: $DEMO_APP_URL"

tee config.main.yml << END
scheme: http
endpoint: $TRACETEST_MAIN_ENDPOINT
END

tee config.target.yml << END
scheme: http
endpoint: $TRACETEST_TARGET_ENDPOINT
END



mkdir -p results/responses

EXIT_STATUS=0
bash ./tests.bash || EXIT_STATUS=$?

exit $EXIT_STATUS
