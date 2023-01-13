# Tracetest + Elastic native APM agent

This repository objective is to show how you can configure your tracetest instance to connect to Elastic stack instance and use it as its tracing backend.

## Steps to start the environment
```bash
docker compose up -d
```

## Open Tracetest UI
Open http://localhost:11633/ to configure the connection to Elasticsearch:
1. In Settings, configure Elastic APM as the Data Store.
2. Set `traces-apm-default` as the Index name.
3. Add the Address and set it to `https://es01:9200`.
4. Set the Username to `elastic` and password to `changeme`.
5. You will need to download the CA certificate from the docker image and upload it to the config under "Upload CA file".
    * The command to download the `ca.crt` file is:
    `docker cp tracetest-elasticapm-with-elastic-agent-es01-1:/usr/share/elasticsearch/config/certs/ca/ca.crt .`
6. Test the commection and Save it, if all is successful.

Create a new test:
1. Use the "HTTP Request" option. Hit Next.
2. Name your test and add a description. Hit Next.
3. Configure the GET url to be `http://app:8080` since the tests will be running in docker compose network. Hit Create.
4. Running the test should succeed.


## Open Kibana
Open https://localhost:5601 and login using `elastic:changeme` credentials. The credentials can be changed in the `.env` file. Navigate to APM (upper lefthand corner menu) -> Services and you should see the `tracetest` service with the rest of the details.

## Steps to stop the environment
```bash
docker compose down -v

# Remove the built app docker image
docker rmi quick-start-nodejs:latest
```

## Project structure
* `docker-compose.yml` - docker compose file that starts the whole environment
    * Elastic stack single node cluster with Elasticsearch, Kibana and the APM Server.
    * OTel collector to support tracetest (TODO - check if can be removed).
    * Tracetest instance.
* `collector-config.yml` - OTel collector configuration file
* `app.js` - sample NodeJS application listening on port 8080 and instrumented with Elastic Nodejs APM agent.
