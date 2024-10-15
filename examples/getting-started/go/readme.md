# Step-by-step

1. Install Dependencies

    ```bash
    go get \
      github.com/gorilla/mux v1.8.1 \
      go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.56.0 \
      go.opentelemetry.io/otel v1.31.0 \
      go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.31.0 \
      go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.31.0 \
      go.opentelemetry.io/otel/sdk v1.31.0 \
      go.opentelemetry.io/otel/trace v1.31.0
    ```

2. Start the App

    ```bash
    export OTEL_SERVICE_NAME=my-service-name
    export OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
    export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"
    export OTEL_EXPORTER_OTLP_HEADERS="x-tracetest-token=bla"

    go run .
    ```
