#!/bin/bash

set -e

TRACETEST_CLI=${TRACETEST_CLI:-"tracetest"}
cmdExitCode=$("$TRACETEST_CLI" &> /dev/null; echo $?)
if [ $cmdExitCode -ne 0 ]; then
  echo "\$TRACETEST_CLI not set to executable. set to $TRACETEST_CLI";
  exit 2
fi

TARGET_URL=${TARGET_URL:-"http://localhost:11633"}
if [  "$TARGET_URL" = "" ]; then
  echo "\$TARGET_URL not set";
  exit 2
fi

DOCKER_COMPOSE="docker compose -f infra/docker-compose.yaml -f ../../examples/docker-compose.demo.yaml"
TRACETEST="tracetest -s $TARGET_URL"

DOCKER_LOG=/tmp/docker-log
# printDockerLog() {
#   $DOCKER_COMPOSE kill
#   echo
#   echo "Error occured. Printing docker log ($DOCKER_LOG):"
#   cat $DOCKER_LOG
# }

# print docker log on any error
# trap 'printDockerLog' ERR

$DOCKER_COMPOSE up > $DOCKER_LOG 2>&1 &
../../scripts/wait-for-port.sh 11633
../../scripts/wait-for-port.sh 8081
sleep 5
$TRACETEST apply test -f tracetest-test.yaml

rm -f ./k6
currentDir=$(pwd)
dir=$(mktemp -d)
cd $dir
go install go.k6.io/xk6/cmd/xk6@latest
xk6 build v0.42.0 --with github.com/kubeshop/xk6-tracetest \
  --replace go.buf.build/grpc/go/prometheus/prometheus=buf.build/gen/go/prometheus/prometheus/protocolbuffers/go@latest \
  --replace go.buf.build/grpc/go/gogo/protobuf=buf.build/gen/go/gogo/protobuf/protocolbuffers/go@latest
mv ./k6 $currentDir
cd $currentDir

./k6 run load-test.js -o xk6-tracetest

$DOCKER_COMPOSE down
