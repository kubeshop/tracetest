# OpenTelemetry Astronomy Shop Demo

The OpenTelemetry Demo is an example application published by the OpenTelemtry CNCF project. It implements an Astronomy shop in a set of microservices in different languages with OpenTelemetry enabled, intended to be used as an example of OpenTelemetry instrumentation and observability. The Tracetest team has made several key contributions to this project, including providing a full suite of end to end tests.

We will provide a full recipe below for running the full demo as well as running the associated Tracetests via Docker. Here are other references you may find useful:

- **Source Code**: https://github.com/open-telemetry/opentelemetry-demo
- **Running it locally in Docker**: [Instructions](https://opentelemetry.io/docs/demo/docker-deployment/)
- **Running on Kubernetes**: [Instructions](https://opentelemetry.io/docs/demo/kubernetes-deployment/)

:::info
Tracetest is part of the official testing harness in the latest version of the OpenTelemetry Demo. Read more in the OpenTelemetry docs, [here](https://opentelemetry.io/docs/demo/tests/).

Or, check out the hands-on workshop on YouTube!

<iframe width="100%" height="250" src="https://www.youtube.com/embed/2MSDy3XHjtE?si=T0ItpwRyE7HbJu5V" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>

:::

## Running the OpenTelemetry Astronomy Shop Demo in Docker

### Prerequisites

- Docker
- Docker Compose v2.0.0+
- 4 GB of RAM for the application

### Get and run the demo

1. Clone the Demo repository:

    ```shell
    git clone https://github.com/open-telemetry/opentelemetry-demo.git
    ```

2. Change to the demo folder:

    ```shell
    cd opentelemetry-demo/
    ```

3. Run docker compose[^1] to start the demo:

    ```shell
    docker compose up --no-build
    ```

    > **Notes:**
    >
    > - The `--no-build` flag is used to fetch released docker images from
    >   [ghcr](https://ghcr.io/open-telemetry/demo) instead of building from
    >   source. Removing the `--no-build` command line option will rebuild all
    >   images from source. It may take more than 20 minutes to build if the
    >   flag is omitted.
    > - If you're running on Apple Silicon, run `docker compose build`[^1] in
    >   order to create local images vs. pulling them from the repository.

## Verify the web store and Telemetry

Once the images are built and containers are started you can access:

- Web store: <http://localhost:8080/>
- Grafana: <http://localhost:8080/grafana/>
- Feature Flags UI: <http://localhost:8080/feature/>
- Load Generator UI: <http://localhost:8080/loadgen/>
- Jaeger UI: <http://localhost:8080/jaeger/ui/>

## Running Tracetests

The Tracetest tests for the OpenTelemetry Demo can be found in the official repo here:

- **Instructions to run (also shown below in this recipe)**: [Running Tracetest Tests](https://github.com/open-telemetry/opentelemetry-demo/tree/main/test#testing-services-with-trace-based-tests)
- **Full source of all tests**: [Source](https://github.com/open-telemetry/opentelemetry-demo/tree/main/test/tracetesting)

To run the entire Test Suite of trace-based tests, run the command:

```sh
make run-tracetesting
#or
docker compose run traceBasedTests
```

To run tests for specific services, pass the name of the service as a
parameter (using the folder names located [here](https://github.com/open-telemetry/opentelemetry-demo/tree/main/test/tracetesting):

```sh
make run-tracetesting SERVICES_TO_TEST="service-1 service-2 ..."
#or
docker compose run traceBasedTests "service-1 service-2 ..."
```

For instance, if you need to run the tests for `ad-service` and
`payment-service`, you can run them with:

```sh
make run-tracetesting SERVICES_TO_TEST="ad-service payment-service"
```

Tracetest will be started on [http://localhost:11633](http://localhost:11633) as part of running these tests and you can view any of the tests, Test Suites, prior runs, or create and run your own tests. It is a great testbed to explore Tracetest!

## Use Cases

- [Add item into shopping cart](./use-cases/add-item-into-shopping-cart.md): Simulate a user choosing an item and adding it to the shopping cart.
- [Check shopping cart content](./use-cases/check-shopping-cart-contents.md): Simulate a user choosing different products and checking the shopping cart later. 
- [Checkout](./use-cases/checkout.md): Simulates a user choosing a product and later doing a checkout of that product, with billing and shipping info.
- [Get recommended products](./use-cases/get-recommended-products.md): Simulates a user querying for recommended products.

## System Architecture

This demonstration environment consists of a series of microservices, handling each aspect of the store, such as Product Catalog, Payment, Currency, etc.

A detailed description of these services can be seen [here](https://opentelemetry.io/docs/demo/services/)
and the architecture diagrams can be seen [here](https://opentelemetry.io/docs/demo/architecture/).
