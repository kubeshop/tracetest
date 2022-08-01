#!/bin/bash

tracetest_main() {
   $TRACETEST_CLI --config ./config.main.yml $@
}

tracetest_target() {
  $TRACETEST_CLI --config ./config.target.yml $@
}

tracetest_target_curl() {
  reqPath=$1
  shift

  curl -sSL "http://$TRACETEST_TARGET_ENDPOINT$reqPath" $@
}

test() {
  name=$1
  definition=$2

  echo -n "-> $name "
  run_test $name $definition
  res=$?
  if [ "$res" = 0 ]; then
    echo -n "OK"
  else
    echo "FAIL"
    echo "$name.json:"
    cat results/responses/$name.json
    echo
    echo "$name.xml:"
    cat results/$name.xml

  fi
  echo
  return $res
}

run_test() {
  name=$1
  definition=$2
  tracetest_main test run --definition $definition --wait-for-result --junit results/$name.xml > results/responses/$name.json

  allPassed=$(cat results/responses/$name.json | jq -rc '.testRun.result.allPassed')
  if [ ! "$allPassed" = "true" ]; then
    return 1
  else
    return 0
  fi
}

require_not_empty() {
  case $1 in
    "" | "null")
      echo $2
      return 1
      ;;

    *)
      return 0
      ;;
  esac
}
