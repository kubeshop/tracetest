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


export TRACETEST_MAIN_ENDPOINT=${TRACETEST_MAIN_ENDPOINT:-"localhost:11633"}
export TRACETEST_TARGET_ENDPOINT=${TRACETEST_TARGET_ENDPOINT:-"localhost:11634"}
export DEMO_APP_URL=${DEMO_APP_URL-"http://demo-pokemon-api.demo"}
export DEMO_APP_GRPC_URL=${DEMO_APP_GRPC_URL-"demo-pokemon-api.demo:8082"}


echo "TRACETEST_CLI: $TRACETEST_CLI"
echo "TARGET_URL: $TARGET_URL"
echo "TRACETEST_MAIN_ENDPOINT: $TRACETEST_MAIN_ENDPOINT"
echo "TRACETEST_TARGET_ENDPOINT: $TRACETEST_TARGET_ENDPOINT"
echo "DEMO_APP_URL: $DEMO_APP_URL"
echo "DEMO_APP_GRPC_URL: $DEMO_APP_GRPC_URL"


cat << EOF > config.main.yml
scheme: http
endpoint: $TRACETEST_MAIN_ENDPOINT
analyticsEnabled: false
EOF

cat << EOF > config.target.yml
scheme: http
endpoint: $TRACETEST_TARGET_ENDPOINT
analyticsEnabled: false
EOF

mkdir -p results/responses

EXIT_STATUS=0
bash ./tests.bash || EXIT_STATUS=$?
bash ./grpc.bash || EXIT_STATUS=$?

exit $EXIT_STATUS
