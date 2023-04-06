# Configuring the Tracetest Server

Tracetest requires a very minimal configuration to be launched, needing just the connection information to connect with the PostgreSQL database which is installed as part of the server install. There are a couple ways to provide this database connection information.

For Docker-based installs, the server configuration file is placed in the ./tracetest/tracetest.yaml file by default when you run the 'tracetest server install' command and select the 'Using Docker Compose' option. The configuration file is mounted to `/app/config.yaml` within the Tracetest Docker container. When Tracetest is run with a 'docker compose -f tracetest/docker-compose.yaml  up -d' command, the server will use the contents of this file to connect to the Postgres database. All other configuration data is stored in the Postgres instance.

This is an example of a tracetest.yaml file:

```yaml
postgres:
  host: postgres
  user: postgres
  password: postgres
  port: 5432
  dbname: postgres
  params: sslmode=disable
```

Alternatively, we support setting a series of environment variables which can contain the connection information for the Postgres instance. If these environment variables are set, they will be used by the Tracetest server to connect to the database.

The list of environment variables and example values is:
- TRACETEST_POSTGRES_HOST - example: postgres
- TRACETEST_POSTGRES_PORT - example: 5432
- TRACETEST_POSTGRES_DBNAME - example: postgres
- TRACETEST_POSTGRES_USER - example: postgres
- TRACETEST_POSTGRES_PASSWORD - example: postgres

You can also 'hydrate' the server with a number of resources the first time it is launched by using [provisioning](./provisioning).

