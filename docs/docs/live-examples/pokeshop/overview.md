# Pokeshop API

As a testing ground, the team at Tracetest has implemented a sample instrumented API around the [PokeAPI](https://pokeapi.co/).

The idea is to have a microservice-divided system that could behave like a typical scenario by having async processes ([RabbitMQ](https://www.rabbitmq.com/)), cache layers ([Redis](https://redis.io/)), database storage ([Postgres](https://www.postgresql.org/)), and simple CRUD interfaces for Pokemons.

With this, users can get familiar with the Tracetest tool by focusing on creating assertions, visualizing the trace, and identifying the different data that comes from the Collector ([Jaeger](https://www.jaegertracing.io/)). Users will learn about basic instrumentation practices like what tools to use, what data to send, when, and what suggested standards need to be followed.

- **Source Code**: https://github.com/kubeshop/pokeshop
- **Running it locally**: [instructions](https://github.com/kubeshop/pokeshop/blob/master/docs/installing.md#run-it-locally)
- **Running on Kubernetes**: [instructions](https://github.com/kubeshop/pokeshop/blob/master/docs/installing.md#run-on-a-kubernetes-cluster)

## Use cases

We have three use cases that use each component of this structure and that can be observed via Open Telemetry and tested with Tracetest. Each one is triggered by an API call to their respective endpoint:

- [Add Pokemon](./use-cases/add-pokemon.md): add a new pokemon only relying on user input into the database
- [List Pokemon](./use-cases/list-pokemon.md): lists all Pokemons registered into Pokeshop
- [Import Pokemon](./use-cases/import-pokemon.md): given a Pokemon ID, this endpoint does an async process, going to PokeAPI to get Pokemon data and adding it to the database

## System architecture

The system is divided into two components: 
- an **API** that serves client requests, 
- a **Worker** who deals with background processes.

The communication between the API and Worker is made using a `RabbitMQ` queue, and both services emit telemetry data to Jaeger and communicate with a Postgres database.

A diagram of the system structure can be seen here:

```mermaid
flowchart TD
    A[(Redis)]
    B[(Postgres)]
    C(Node.js API)
    D(RabbitMQ)
    E(Worker)
    F(Jaeger)

    C --> |IORedis| A
    C --> |Sequelize| B
    C --> |ampqlib| D
    D --> |ampqlib| E
    E --> |Sequelize| B
    C --> |OpenTelemetry Node.js SDK| F
    E --> |OpenTelemetry Node.js SDK| F
```

In our live tests, we are deploying it into a single Kubernetes namespace, deployed via a [helm chart](https://github.com/kubeshop/pokeshop/blob/master/docs/installing.md#run-on-a-kubernetes-cluster).

The Pokeshop API is only accessible from within the Kubernetes cluster network as the Tracetest needs to be able to reach it.
