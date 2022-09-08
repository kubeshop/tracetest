#!/bin/sh

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


if [ "$RESTART" == "yes" ]; then
  docker compose -f docker-compose.yaml -f local-config/docker-compose.testrunner.yaml down
fi
docker compose -f docker-compose.yaml up -d --build --remove-orphans
docker compose -f docker-compose.yaml -f local-config/docker-compose.testrunner.yaml build
docker compose -f docker-compose.yaml -f local-config/docker-compose.testrunner.yaml run testrunner

if [ "$STOP" == "yes" ]; then
  docker compose -f docker-compose.yaml stop
fi
