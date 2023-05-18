#!/bin/bash

##
## EXAMPLE USAGE
## JAEGER_UI_URL="http://10.1.2.11:16686" ./test-kubelet.bash ./test.yaml "kubectl apply -f https://raw.githubusercontent.com/kubernetes/website/main/content/en/examples/controllers/nginx-deployment.yaml" default/nginx
##

JAEGER_UI_URL=${JAEGER_UI_URL:-"http://127.0.0.1:16686"}
DEBUG=${DEBUG}
TEST_FILE=$1
KUBECTL_CMD=$2
POD_PREFIX=$3

echo "-> Parameters:"
echo "      JAEGER_UI_URL: ${JAEGER_UI_URL}"
echo "      TEST_FILE:     ${TEST_FILE}"
echo "      KUBECTL_CMD:   ${KUBECTL_CMD}"
echo "      POD_PREFIX:    ${POD_PREFIX}"

echo

outputFile=$(mktemp)
debugInfo=$(mktemp)
$KUBECTL_CMD > $outputFile 2>  $debugInfo
exitCode=$?

if [ $exitCode -ne 0 ]; then
  echo "!!!!!!!!!!!!!!!!!!!!!!!!!"
  echo "command failed with code $exitCode"

  echo "**** output file: $outputFile"
  echo "**** debug info output file: $debugInfo"
  echo "**** Debug info output"
  cat $debugInfo
  exit $exitCode
fi

echo "-> output"
cat $outputFile
echo

if [ "$DEBUG" == "yes" ]; then
  echo "**** output file: $outputFile"
  echo "**** debug info output file: $debugInfo"
  echo "**** Debug info output"
  cat $debugInfo
fi

echo "-> wait a moment so things propagate correctly"
sleep 5

pods=$(kubectl get pods --all-namespaces -o=jsonpath='{range .items[*]}{.metadata.namespace}/{.metadata.name}{"\n"}{end}' | grep $POD_PREFIX | head -n1)

if [ "$DEBUG" == "yes" ]; then
  echo "**** Pods: $pods"
fi

START=$(date -v-5M +%s)
END=$(date +%s)

curlDebug=$(mktemp)
traces=$(mktemp)
curl -v --get "${JAEGER_UI_URL}/api/traces" \
  --data-urlencode "start=${START}000000" \
  --data-urlencode "end=${END}000000" \
  --data-urlencode "minDuration" \
  --data-urlencode "maxDuration" \
  --data-urlencode "limit=20" \
  --data-urlencode "lookback=5m" \
  --data-urlencode "service=kubelet" \
  --data-urlencode 'tags={"k8s.pod":"'$pods'"}' > $traces 2> $curlDebug

if [ "$DEBUG" == "yes" ]; then
  echo "**** curl debug output file: $curlDebug"
  echo "**** curl response output file: $traces"
  echo "**** cURL request debug output"
  cat $curlDebug

  echo "**** Traces"
  cat $traces
fi

TRACE_ID=$(cat $traces | jq -r '.data | first' | jq -r '.traceID')
echo "-> TraceID": $TRACE_ID

testFile=$(mktemp)
tracetestCommand="tracetest test run --definition $testFile --wait-for-result"
cat $TEST_FILE | sed "8s/.*/      id: $TRACE_ID/" > $testFile

DEBUG=yes
if [ "$DEBUG" == "yes" ]; then
  echo "**** replaced test file: $testFile"
  cat $testFile

  echo "**** tracetest command: $tracetestCommand"
fi

echo "-> Running test"
$tracetestCommand
