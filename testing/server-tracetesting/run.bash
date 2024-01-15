#/bin/bash

export TRACETEST_CLI=${TRACETEST_CLI:-"tracetest"}
cmdExitCode=$("$TRACETEST_CLI" &> /dev/null; echo $?)
if [ $cmdExitCode -ne 0 ]; then
  echo "\$TRACETEST_CLI not set to executable. set to $TRACETEST_CLI";
  exit 2
fi

export TARGET_URL=${TARGET_URL:-"http://localhost:11633"}
if [  "$TARGET_URL" = "" ]; then
  echo "\$TARGET_URL not set";
  exit 2
fi

export TRACETEST_ENDPOINT=${TRACETEST_ENDPOINT:-"localhost:11633"}
export DEMO_APP_URL=${DEMO_APP_URL-"http://demo-pokemon-api.demo"}
export DEMO_APP_GRPC_URL=${DEMO_APP_GRPC_URL-"demo-pokemon-api.demo:8082"}
export DEMO_APP_KAFKA_BROKER=${DEMO_APP_KAFKA_BROKER-"stream:9092"}

# TODO: think how to move this id generation to HTTP Test suite
export EXAMPLE_TEST_ID="w2ON-RVVg"

echo "Preparing to run tests on API..."
echo ""

echo "Variable set considered on this run:"
echo "TRACETEST_CLI:          $TRACETEST_CLI"
echo "TARGET_URL:             $TARGET_URL"
echo "TRACETEST_ENDPOINT:     $TRACETEST_ENDPOINT"
echo "DEMO_APP_URL:           $DEMO_APP_URL"
echo "DEMO_APP_GRPC_URL:      $DEMO_APP_GRPC_URL"
echo "DEMO_APP_KAFKA_BROKER:  $DEMO_APP_KAFKA_BROKER"

cat << EOF > tracetesting-vars.yaml
type: VariableSet
spec:
  id: tracetesting-vars
  name: tracetesting-vars
  values:
    - key: TARGET_URL
      value: $TARGET_URL
    - key: DEMO_APP_URL
      value: $DEMO_APP_URL
    - key: DEMO_APP_GRPC_URL
      value: $DEMO_APP_GRPC_URL
    - key: DEMO_APP_KAFKA_BROKER
      value: $DEMO_APP_KAFKA_BROKER
    - key: EXAMPLE_TEST_ID
      value: $EXAMPLE_TEST_ID
EOF

echo "variables set created:"
cat tracetesting-vars.yaml

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

  definition='./features/'$feature'/_test_suite.yml'

  testCMD="$TRACETEST_CLI  --config ./config.yml run testsuite --file $definition --vars ./tracetesting-vars.yaml"
  echo $testCMD
  $testCMD
  return $?
}

echo "Test helpers set."
echo ""

echo "Starting tests..."

EXIT_STATUS=0

# add more test suites here
run_test_suite_for_feature 'http_test' || (EXIT_STATUS=$? && echo "HTTP Test suite failed")
run_test_suite_for_feature 'grpc_test' || (EXIT_STATUS=$? && echo "GRPC Test suite failed")
run_test_suite_for_feature 'kafka_test' || (EXIT_STATUS=$? && echo "Kafka Test suite failed")
run_test_suite_for_feature 'variableset' || (EXIT_STATUS=$? && echo "VariableSet Test suite failed")
run_test_suite_for_feature 'testsuite' || (EXIT_STATUS=$? && echo "TestSuite Test suite failed")

echo ""
echo "Tests done! Exit code: $EXIT_STATUS"

rm tracetesting-vars.yaml

exit $EXIT_STATUS

