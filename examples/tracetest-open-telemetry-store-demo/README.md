# Running Tracetest with Open Telemetry store demo

This folder contains the assets to run Tracetest with Open Telemetry Community store demo locally.
To see the complete documentation go to: [https://docs.tracetest.io/live-examples/opentelemetry-store/overview](https://docs.tracetest.io/live-examples/opentelemetry-store/overview)

## Running locally

To run the OpenTelemetry demo locally just clone their repo in your machine:
```sh
git clone https://github.com/open-telemetry/opentelemetry-demo.git
```

And then, run:
```sh
docker compose up --no-build
```

After a few minutes, the store should be running normally on your machine, to test it go to a browser and access: [http://localhost:8080](http://localhost:8080)

To start Tracetest with this demo, just go to this folder, and run `docker compose up`.

After that, Tracetest will start on [http://localhost:11633](http://localhost:11633) and you can start creating tests.
