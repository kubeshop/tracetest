#/bin/bash

set -e

export TRACETEST_CLI=${TRACETEST_CLI:-"tracetest"}
if ! command -v "$TRACETEST_CLI" &> /dev/null; then
  echo "\$TRACETEST_CLI not set to executable. set to $TRACETEST_CLI";
  exit 2
fi

export TRACETEST_ENDPOINT="localhost:11633"

echo "Preparing to run CLI tests..."
echo ""

echo "Environment variables considered on this run:"
echo "TRACETEST_CLI:      $TRACETEST_CLI"
echo "TRACETEST_ENDPOINT: $TRACETEST_ENDPOINT"

echo "Setting up tracetest CLI configuration..."
cat << EOF > config.yml
scheme: http
endpoint: $TRACETEST_ENDPOINT
analyticsEnabled: false
EOF
echo "tracetest CLI set up."
echo ""

echo "Setting up test helpers..."

run_test() {
  test_definition_file=$1

  $TRACETEST_CLI --config ./config.yml test run --definition $test_definition_file --wait-for-result
  return $?
}

echo "Test helpers set."
echo ""

echo "Starting tests..."

EXIT_STATUS=0

# add more test here
run_test './base/simple-test.yaml' || EXIT_STATUS=$?

echo ""
echo "Tests done! Exit code: $EXIT_STATUS"

exit $EXIT_STATUS
