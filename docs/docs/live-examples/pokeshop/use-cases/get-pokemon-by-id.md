# Pokeshop API - Get Pokemon by ID

This use case retrieves a specific Pokemon from the cache if it is cached, or from the database (Postgres) also populating the cache. The idea of this query is to showcase a straightforward scenario, where the API layer receives a request from the outside and needs to evaluate a different behavior depending of his dependencies.

```mermaid
sequenceDiagram
    participant Endpoint as GET /pokemon/:id
    participant API as API
    participant Database as Postgres
    participant Cache as Redis
    
    Endpoint->>API: request

    API->>Cache: query cache
    Cache-->>API: cache response

    alt cache hit
      API-->>Endpoint: 200 OK <br> <Pokemon object>
    else cache miss
      API->>Database: get specific of pokemon
      Database-->>API: pokemon

      API->>Cache: populate cache
      Cache-->>API: ok

      API-->>Endpoint: 200 OK <br> <Pokemon object>
    end
```

You can trigger this use case by calling the endpoint `GET /pokemon/25` without payload, and should receive a payload similar to this: 
```json
{
  "id":  25,
  "name":  "pikachu",
  "type":  "electric",
  "imageUrl":  "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png",
  "isFeatured":  true
}
```

## Building a test for those scenarios

Using Tracetest, we can [create a test](../../../web-ui/creating-tests.md) that will execute an API call on `GET /pokemon/25` and validate two scenarios:
1. **an API call with a cache hit**
   - The API should return a valid result with HTTP 200 OK
   - The cache should be queried
   - The database should not be queried
2. **an API call with a cache miss**
   - The API should return a valid result with HTTP 200 OK
   - The cache should be queried 
   - The cache should be populated
   - The database should be queried

### Traces

Running these tests for the first time will create an Observability trace with two different shapes, depending on the cache situation.

1. **cache miss**, where we can see spans from the API, database, and cache:
![](../images/get-pokemon-by-id-trace-cachemiss.png)

2. **cache hit**, where we can see spans from the API and cache:
![](../images/get-pokemon-by-id-trace-cachehit.png)


### Assertions

With this trace, now we can build [assertions](../../../concepts/assertions.md) on Tracetest and validate API, cache, and database responses:

- [both cases] The API should return a valid result with HTTP 200 OK
![](../images/get-pokemon-by-id-api-test-spec.png)

- [both cases] The cache should be queried
![](../images/get-pokemon-by-id-redis-query-test-spec.png)

- [cache hit] The database should not be queried
TODO

- [cache miss] The cache should be populated
![](../images/get-pokemon-by-id-redis-set-test-spec.png)

- [cache miss] The database should be queried
![](../images/get-pokemon-by-id-db-query-test-spec.png)

And that's all, now you can validate this entire use case.

### Test Definition

If you want to replicate this entire test on Tracetest see by yourself, you can replicate these steps on our Web UI or using our CLI, saving the following test definition as the file `test-definition.yml` and later running:

```sh
tracetest test -d test-definition.yml --wait-for-results
```

```yaml
type: Test
spec:
  name: Pokeshop - Get Pokemon by ID [cache miss]
  description: Get a Pokemon by ID
  trigger:
    type: http
    httpRequest:
      url: http://demo-pokemon-api.demo/pokemon/${env:POKEMON_ID}
      method: GET
      headers:
      - key: Content-Type
        value: application/json
  specs:
  - selector: span[tracetest.span.type="http" http.method="GET"]
    assertions:
    - attr:http.status_code         =         200
    - attr:http.response.body   contains     '"id":${env:POKEMON_ID}'
  - selector: span[tracetest.span.type="database" db.system="redis" db.operation="get"]
    assertions:
    - attr:name   =   "get pokemon-${env:POKEMON_ID}"
  - selector: span[tracetest.span.type="database" db.system="redis" db.operation="set"]
    assertions:
    - attr:name  =  "set pokemon-${env:POKEMON_ID}"
  - selector: span[tracetest.span.type="database" name="findOne pokeshop.pokemon"
      db.system="postgres" db.name="pokeshop" db.operation="findOne" db.sql.table="pokemon"]
    assertions:
    - attr:name  =  "findOne pokeshop.pokemon"
```
