#/bin/bash

set -e

export TAG=${TAG:-"latest"}
export TRACETEST_ENDPOINT="host.docker.internal:11633"
export TRACETEST_CLI="tracetest"

echo "Preparing to run CLI tests..."
echo ""

echo "Environment variables considered on this run:"
echo "TAG:                $TAG"
echo "TRACETEST_ENDPOINT: $TRACETEST_ENDPOINT"
echo "TRACETEST_CLI:      $TRACETEST_CLI"

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

run_test() {
  test_definition_file=./files/$1

  docker run --volume $(PWD):/app/files --entrypoint $TRACETEST_CLI kubeshop/tracetest:$TAG --config ./files/config.yml test run --definition $test_definition_file --wait-for-result
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
