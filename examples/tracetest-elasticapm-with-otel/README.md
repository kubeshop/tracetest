# Tracetest + Elastic APM + OTel

This repository objective is to show how you can configure your tracetest instance to connect to Elastic stack instance and use it as its tracing backend.

## Steps to start the environment
```bash
docker compose up -d
```

## Open Tracetest UI
Open http://localhost:11633/ and create a new test:
1. Use the "HTTP Request" option. Hit Next.
2. Name your test and add a description. Hit Next.
3. Configure the GET url to be `http://app:8080` since the tests will be running in docker compose network. Hit Create.


## Open Kibana
Open https://localhost:5601 and login using `elastic:changeme` credentials. The credentials can be changed in the `.env` file. Navigate to APM (upper lefthand corner menu) -> Services and you should see the `tracetest` service with the rest of the details.

## Steps to stop the environment
```bash
docker compose down -v
```

## Project structure
* `docker-compose.yml` - docker compose file that starts the whole environment
    * Elastic stack single node cluster with Elasticsearch, Kibana and the APM Server.
    * OTel collector sending the OTel trace data to Elastic stack.
    * Tracetest instance configured to send the trace data to OTel collector.
* `collector-config.yml` - OTel collector configuration file
* `*.js/*.json` - sample NodeJS application listening on port 8080
