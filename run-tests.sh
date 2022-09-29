#!/bin/sh

set -ex

STOP=yes
RESTART=yes

TESTS=()

while [[ $# -gt 0 ]]; do
  case $1 in
    --no-stop)
      STOP=no
      shift
      ;;
    --no-restart)
      RESTART=no
      shift
      ;;
    cypress)
      TESTS+=("cypress")
      shift
      ;;
    dogfood)
      TESTS+=("dogfood")
      shift
      ;;
    *)
      echo "Unknown option $1"
      help_message
      exit 1
      ;;
  esac
done

if [ ${#TESTS[@]} -eq 0 ]; then
  echo "missing test type. usage: ./run-tests.sh [--dogfood] [--cypress]"
  exit 1
fi

echo "Running tests: "${TESTS[@]}

opts="-f docker-compose.yaml -f examples/docker-compose.demo.yaml"

if [ "$RESTART" == "yes" ]; then
  docker compose $opts down
fi
docker compose $opts up -d --build --remove-orphans
docker compose $opts -f local-config/docker-compose.testrunner.yaml build ${TESTS[@]}
docker compose $opts -f local-config/docker-compose.testrunner.yaml run ${TESTS[@]}

if [ "$STOP" == "yes" ]; then
  docker compose $opts stop
fi
