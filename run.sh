#!/bin/sh

set -e

opts="-f docker-compose.yaml -f examples/docker-compose.demo.yaml"

restart() {
  make build-docker
  docker compose $opts kill tracetest
  docker compose $opts up -d tracetest
}

logs() {
  docker compose $opts logs -f
}

logstt() {
  docker compose $opts logs -f tracetest
}

ps() {
  docker compose $opts ps
}

down() {
  docker compose $opts down
}

up() {
  make build-docker
  docker compose $opts up -d --remove-orphans
}

test() {

  echo "Running tests: "${TESTS[@]}

  docker compose $opts -f local-config/docker-compose.testrunner.yaml build ${TESTS[@]}
  docker compose $opts -f local-config/docker-compose.testrunner.yaml run ${TESTS[@]}
}

TESTS=()
CMD=()

while [[ $# -gt 0 ]]; do
  case $1 in
    cypress)
      TESTS+=("cypress")
      shift
      ;;
    dogfood)
      TESTS+=("dogfood")
      shift
      ;;
    up)
      CMD+=("up")
      shift
      ;;
    down)
      CMD+=("down")
      shift
      ;;
    test)
      CMD+=("test")
      shift
      ;;
    logstt)
      CMD+=("logstt")
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
  echo "missing command. usage: ./run.sh [up|down|ps|logs|restart|test|dogfood|cypress]"
  exit 1
fi

for cmd in "${CMD[@]}"; do
   $cmd
done
