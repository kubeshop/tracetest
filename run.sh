#!/bin/bash

set -e

export TAG=${TAG:-dev}

opts="-f docker-compose.yaml -f examples/docker-compose.demo.yaml"

help_message() {
  echo "usage: ./run.sh [cypress|tracetests|up|stop|build|down|tracetest-logs|logs|ps|restart]"
}

restart() {
  docker compose $opts kill tracetest
  docker compose $opts up -d tracetest
  docker compose $opts restart otel-collector
}

logs() {
  docker compose $opts logs -f
}

tracetest-logs() {
  docker compose $opts logs -f tracetest
}

ps() {
  docker compose $opts ps
}

down() {
  docker compose $opts down
}

build() {
  make build-docker
  # the previous commands builds the cli binary for linux (because its the os in docker)
  # if the script is run on another os, like macos, we need to rebuild for the binary to match the os
  make dist/tracetest
}

up() {
  docker compose $opts up -d --remove-orphans
}

stop() {
  docker compose $opts stop
}

cypress-ci() {

  echo "Running cypress"


  export CYPRESS_BASE_URL=http://localhost:11633
  export POKEMON_HTTP_ENDPOINT=http://demo-api:8081

  cd web
  npm run cy:ci
}

cypress() {

  echo "Running cypress"


  export CYPRESS_BASE_URL=http://localhost:11633
  export POKEMON_HTTP_ENDPOINT=http://demo-api:8081

  cd web
  npm run cy:run
}

tracetests() {

  echo "Running tracetests"

  SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

  export TRACETEST_CLI=${SCRIPT_DIR}/dist/tracetest
  export TARGET_URL=http://localhost:11633
  export TRACETEST_ENDPOINT=localhost:11633
  export DEMO_APP_URL=http://demo-api:8081
  export DEMO_APP_GRPC_URL=demo-rpc:8082

  cd testing/server-tracetesting
  ./run.bash
}

CMD=()

while [[ $# -gt 0 ]]; do
  case $1 in
    cypress)
      CMD+=("cypress")
      shift
      ;;
    cypress-ci)
      CMD+=("cypress-ci")
      shift
      ;;
    tracetests)
      CMD+=("tracetests")
      shift
      ;;
    up)
      CMD+=("up")
      shift
      ;;
    stop)
      CMD+=("stop")
      shift
      ;;
    build)
      CMD+=("build")
      shift
      ;;
    down)
      CMD+=("down")
      shift
      ;;
    tracetest-logs)
      CMD+=("tracetest-logs")
      shift
      ;;
    logs)
      CMD+=("logs")
      shift
      ;;
    ps)
      CMD+=("ps")
      shift
      ;;
    restart)
      CMD+=("restart")
      shift
      ;;

    *)
      echo "Unknown option $1"
      help_message
      exit 1
      ;;
  esac
done

if [ ${#CMD[@]} -eq 0 ]; then
  echo "Missing command"
  help_message
  exit 1
fi

for cmd in "${CMD[@]}"; do
   $cmd
done
