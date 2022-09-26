#!/bin/sh

set -e

STOP=yes
RESTART=yes

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
    -*|--*)
      echo "Unknown option $1"
      help_message
      exit 1
      ;;
  esac
done


opts="-f docker-compose.yaml -f examples/docker-compose.demo.yaml"

if [ "$RESTART" == "yes" ]; then
  docker compose $opts down
fi
docker compose $opts up -d --build --remove-orphans
docker compose $opts -f local-config/docker-compose.testrunner.yaml build
docker compose $opts -f local-config/docker-compose.testrunner.yaml run testrunner

if [ "$STOP" == "yes" ]; then
  docker compose $opts stop
fi
