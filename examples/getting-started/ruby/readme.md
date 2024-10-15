# Step-by-step

1. Install Dependencies

    ```bash
    gem install opentelemetry-sdk
                opentelemetry-exporter-otlp
                opentelemetry-instrumentation-all
    ```

2. Configure OpenTelemetry

      ```ruby
      require 'opentelemetry/sdk'
      require 'opentelemetry/exporter/otlp'
      require 'opentelemetry/instrumentation/all'

      OpenTelemetry::SDK.configure do |c|
          c.use_all() # enables all instrumentation!
      end
      ```

3. Start the App

      ```bash
      export OTEL_SERVICE_NAME=my-service-name
      export OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
      export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"
      export OTEL_EXPORTER_OTLP_HEADERS="x-tracetest-token=<token>"

      rails server -p 8080
      ```
