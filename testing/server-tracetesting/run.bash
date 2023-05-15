#/bin/bash

set -e

export TRACETEST_CLI=${TRACETEST_CLI:-"tracetest"}
if ! command -v "$TRACETEST_CLI" &> /dev/null; then
  echo "\$TRACETEST_CLI not set to executable. set to $TRACETEST_CLI";
  exit 2
fi

export TARGET_URL=${TARGET_URL:-"http://localhost:8081"}
if [  "$TARGET_URL" = "" ]; then
  echo "\$TARGET_URL not set";
  exit 2
fi

export TRACETEST_ENDPOINT=${TRACETEST_ENDPOINT:-"localhost:11633"}
export DEMO_APP_URL=${DEMO_APP_URL-"http://demo-pokemon-api.demo"}
export DEMO_APP_GRPC_URL=${DEMO_APP_GRPC_URL-"demo-pokemon-api.demo:8082"}

# TODO: think how to move this id generation to HTTP Test suite
export EXAMPLE_TEST_ID="w2ON-RVVg"

echo "Preparing to run tests on API..."
echo ""

echo "Environment variables considered on this run:"
echo "TRACETEST_CLI:      $TRACETEST_CLI"
echo "TARGET_URL:         $TARGET_URL"
echo "TRACETEST_ENDPOINT: $TRACETEST_ENDPOINT"
echo "DEMO_APP_URL:       $DEMO_APP_URL"
echo "DEMO_APP_GRPC_URL:  $DEMO_APP_GRPC_URL"

cat << EOF > tracetesting-env.yaml
type: Environment
spec:
  id: tracetesting-env
  name: tracetesting-env
  values:
    - key: TARGET_URL
      value: $TARGET_URL
    - key: DEMO_APP_URL
      value: $DEMO_APP_URL
    - key: DEMO_APP_GRPC_URL
      value: $DEMO_APP_GRPC_URL
    - key: EXAMPLE_TEST_ID
      value: $EXAMPLE_TEST_ID
EOF

echo "Environment variables set:"
cat tracetesting-env.yaml

echo "Setting up tracetest CLI configuration..."
cat << EOF > config.yml
scheme: http
endpoint: $TRACETEST_ENDPOINT
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

  $TRACETEST_CLI --config ./config.yml test run --definition $definition --environment ./tracetesting-env.yaml --wait-for-result --junit $junit_output
  return $?
}

echo "Test helpers set."
echo ""

echo "Starting tests..."

EXIT_STATUS=0

# add more test suites here
run_test_suite_for_feature 'http_test' || EXIT_STATUS=$?
run_test_suite_for_feature 'grpc_test' || EXIT_STATUS=$?
run_test_suite_for_feature 'environment' || EXIT_STATUS=$?
run_test_suite_for_feature 'transaction' || EXIT_STATUS=$?

echo ""
echo "Tests done! Exit code: $EXIT_STATUS"

exit $EXIT_STATUS

