#!/bin/bash

JAEGER_UI_URL=${JAEGER_UI_URL:-"http://127.0.0.1:16686"}
DEBUG=${DEBUG}
KUBECTL_CMD=$1

echo "-> Parameters:"
echo "      JAEGER_UI_URL: ${JAEGER_UI_URL}"
echo "      KUBECTL_CMD: ${KUBECTL_CMD}"
echo


outputFile=$(mktemp)
debugInfo=$(mktemp)
$KUBECTL_CMD -v 9 > $outputFile 2>  $debugInfo

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

AUDIT_ID=$(cat $debugInfo | grep "Audit-Id:" | tr -s ' ' | cut -d' ' -f6)
echo "-> Audit-Id": $AUDIT_ID

START=$(date -v-1H +%s)
END=$(date +%s)

curlDebug=$(mktemp)
traces=$(mktemp)
curl -v --get "${JAEGER_UI_URL}/api/traces" \
  --data-urlencode "start=${START}000000" \
  --data-urlencode "end=${END}000000" \
  --data-urlencode "minDuration" \
  --data-urlencode "maxDuration" \
  --data-urlencode "limit=20" \
  --data-urlencode "lookback=1h" \
  --data-urlencode "service=apiserver" \
  --data-urlencode 'tags={"audit-id":"'$AUDIT_ID'"}' > $traces 2> $curlDebug

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
