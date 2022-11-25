# Database configuration

Tracetest requires a Postgres instance for it to keep trace of tests, transactions, environments, and so. So, in order to Tracetest even start, it requires you to provide a connection string for it to be able to start using the database.

Check [Postgres official documentation](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING) for more details on how to customize the connection string

#### Example
```yaml
# tracetest.config.yaml

postgresConnString: "host=my-database-url.com user=user-for-tracetest password=password-for-the-user port=5432 sslMode=disabled"
```
