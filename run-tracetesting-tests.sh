#!/bin/sh

docker compose -f docker-compose.yaml -f docker-compose.testrunner.yaml build
docker compose -f docker-compose.yaml -f docker-compose.testrunner.yaml run testrunner
docker compose -f docker-compose.yaml stop
