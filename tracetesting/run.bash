#/bin/bash

set -e

export TRACETEST_CLI_MAIN=${TRACETEST_CLI_MAIN:-"tracetest"}
if ! command -v "$TRACETEST_CLI_MAIN" &> /dev/null; then
  echo "\$TRACETEST_CLI_MAIN not set to executable. set to $TRACETEST_CLI_MAIN";
  exit 2
fi

export TARGET_URL=${TARGET_URL:-"http://localhost:8081"}
if [  "$TARGET_URL" = "" ]; then
  echo "\$TARGET_URL not set";
  exit 2
fi

export TRACETEST_MAIN_ENDPOINT=${TRACETEST_MAIN_ENDPOINT:-"localhost:11633"}
export DEMO_APP_URL=${DEMO_APP_URL-"http://demo-pokemon-api.demo"}
export DEMO_APP_GRPC_URL=${DEMO_APP_GRPC_URL-"demo-pokemon-api.demo:8082"}

echo "Preparing to run tests on API..."
echo ""

echo "Environment variables considered on this run:"
echo "TRACETEST_CLI_MAIN:      $TRACETEST_CLI_MAIN"
echo "TARGET_URL:              $TARGET_URL"
echo "TRACETEST_MAIN_ENDPOINT: $TRACETEST_MAIN_ENDPOINT"
echo "DEMO_APP_URL:            $DEMO_APP_URL"
echo "DEMO_APP_GRPC_URL:       $DEMO_APP_GRPC_URL"

cat << EOF > .main.env
TRACETEST_CLI_MAIN=$TRACETEST_CLI_MAIN
TARGET_URL=$TARGET_URL
TRACETEST_MAIN_ENDPOINT=$TRACETEST_MAIN_ENDPOINT
DEMO_APP_URL=$DEMO_APP_URL
DEMO_APP_GRPC_URL=$DEMO_APP_GRPC_URL
EOF

echo ""

echo "Setting up tracetest CLI configuration..."
cat << EOF > config.main.yml
scheme: http
endpoint: $TRACETEST_MAIN_ENDPOINT
analyticsEnabled: false
EOF
echo "tracetest CLI set up."
echo ""

echo "Setting up test helpers..."

mkdir -p results/responses

run_test_suite_for_feature() {
  feature=$1

  junit_output='results/'$feature'_test_suite.xml'
  definition='./features/'$feature'/_test_suite.yml'

  $TRACETEST_CLI_MAIN --config ./config.main.yml test run --definition $definition --environment ./.main.env --wait-for-result --junit $junit_output
  return $?
}

echo "Test helpers set."
echo ""

echo "Starting tests..."

EXIT_STATUS=0

# add more test suites here
run_test_suite_for_feature 'grpc_test' || EXIT_STATUS=$?
run_test_suite_for_feature 'environment' || EXIT_STATUS=$?
run_test_suite_for_feature 'transaction' || EXIT_STATUS=$?

echo ""
echo "Tests done! Exit code: $EXIT_STATUS"

exit $EXIT_STATUS

