#/bin/bash

set -e

export TAG=${TAG:-"latest"}
export TEST_ENV=${TEST_ENV:-"local"}

if ! command -v "$TEST_ENV" &> /dev/null; then
  echo "\$TEST_ENV not set. stopping script";
  exit 2
fi

if [ $TEST_ENV = "local" ]; then
  export TRACETEST_ENDPOINT="localhost:11633"
  export TRACETEST_CLI_COMMAND="tracetest"
else
  export TRACETEST_ENDPOINT="host.docker.internal:11633"
  export TRACETEST_CLI_COMMAND="docker run --volume $PWD/tests:/app/tests --entrypoint tracetest kubeshop/tracetest:$TAG"
fi

echo "Preparing to run CLI tests..."
echo ""

echo "Environment variables considered on this run:"
echo "TAG:                   $TAG"
echo "TEST_ENV:              $TEST_ENV"
echo "TRACETEST_ENDPOINT:    $TRACETEST_ENDPOINT"
echo "TRACETEST_CLI_COMMAND: $TRACETEST_CLI_COMMAND"

echo "Setting up tracetest CLI configuration..."
cat << EOF > config.yml
scheme: http
endpoint: $TRACETEST_ENDPOINT
analyticsEnabled: false
EOF
echo "tracetest CLI set up."
echo ""

echo "Setting up test helpers..."

run_cli_command() {
  test_definition_file=$args

  $TRACETEST_CLI_COMMAND --config ./config.yml $args
  return $?
}

echo "Test helpers set."
echo ""

echo "Starting tests..."

EXIT_STATUS=0

run_cli_command '--help' || EXIT_STATUS=$?
run_cli_command 'version' || EXIT_STATUS=$?
run_cli_command 'test run --definition ./tests/simple-test.yaml --wait-for-result' || EXIT_STATUS=$?

echo ""
echo "Tests done! Exit code: $EXIT_STATUS"

exit $EXIT_STATUS
