# OpenTelemetry Astronomy Shop Demo

This system implements an Astronomy shop in a set of microservices in different languages with OpenTelemetry enabled, intended to be used as an example of OpenTelemetry instrumentation and observability.

- **Source Code**: https://github.com/open-telemetry/opentelemetry-demo
- **Running it Locally**: [Instructions](https://github.com/open-telemetry/opentelemetry-demo/blob/main/docs/docker_deployment.md#run-docker-compose)
- **Running on Kubernetes**: [Instructions](https://github.com/open-telemetry/opentelemetry-demo/blob/main/docs/kubernetes_deployment.md)

## Running with Tracetest

To run the this demo locally with Tracetest, first clone OpenTelemetry demo repo in your machine in any folder:
```sh
git clone https://github.com/open-telemetry/opentelemetry-demo.git
```

And then, run in that folder:
```sh
docker compose up --no-build
```

After a few minutes, the store should be running normally in your machine, to test it go to a browser and access: [http://localhost:8080](http://localhost:8080)

Now, to start Tracetest connected with this demo, download the contents of the [Running Tracetest with OpenTelemetry store demo](https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-open-telemetry-store-demo) in any folder, and then run `docker compose up`.

After that, Tracetest will start on [http://localhost:11633](http://localhost:11633) and you can start creating tests.

## Use Cases

- [Add item into shopping cart](./use-cases/add-item-into-shopping-cart.md): Simulate a user choosing an item and adding it to the shopping cart.
- [Check shopping cart content](./use-cases/check-shopping-cart-contents.md): Simulate a user choosing different products and checking the shopping cart later. 
- [Checkout](./use-cases/checkout.md): Simulates a user choosing a product and later doing a checkout of that product, with billing and shipping info.
- [Get recommended products](./use-cases/get-recommended-products.md): Simulates a user querying for recommended products.

## System Architecture

This demonstration environment consists of a series of microservices, handling each aspect of the store, such as Product Catalog, Payment, Currency, etc.

A detailed description of these services can be seen [here](https://opentelemetry.io/docs/demo/services/)
and the architecture diagrams can be seen [here](https://opentelemetry.io/docs/demo/architecture/).
