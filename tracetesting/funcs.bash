#!/bin/bash

export TRACETEST_SCHEME="http"
tracetest_main() {
  TRACETEST_ENDPOINT=$TRACETEST_MAIN_ENDPOINT $TRACETEST_CMD $@
}

tracetest_target() {
  >&2 echo TRACETEST_ENDPOINT=$TRACETEST_TARGET_ENDPOINT $TRACETEST_CMD $@
  TRACETEST_ENDPOINT=$TRACETEST_TARGET_ENDPOINT $TRACETEST_CMD $@
}

tracetest_target_curl() {
  reqPath=$1
  shift

  curl "$TRACETEST_SCHEME://$TRACETEST_TARGET_ENDPOINT$reqPath" $@
}

run_test() {
  name=$1
  definition=$2
  tracetest_main test run --definition $definition --wait-for-result --junit results/$name.xml > results/responses/$name.json

  allPassed=$(cat results/responses/$name.json | jq -rc '.testRun.result.allPassed')
  if [ ! "$allPassed" = "true" ]; then
    echo "-> $name FAIL"
    return 1
  else
    echo "-> $name OK"
    return 0
  fi
}
