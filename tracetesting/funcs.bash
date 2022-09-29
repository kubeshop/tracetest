#!/bin/bash

tracetest_main() {
   $TRACETEST_CLI_MAIN --config ./config.main.yml $@
}

tracetest_target() {
  $TRACETEST_CLI_TARGET --config ./config.target.yml $@
}

tracetest_target_curl() {
  reqPath=$1
  shift

  curl -sSL "http://$TRACETEST_TARGET_ENDPOINT$reqPath" $@
}

test() {
  name=$1
  definition=$2

  run_test $name $definition
  res=$?
  return $res
}

run_test() {
  name=$1
  definition=$2
  tracetest_main test run --definition $definition --wait-for-result --junit results/$name.xml
  return $?
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
