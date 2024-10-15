export OTEL_SERVICE_NAME=my-service-name
export OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"

opentelemetry-instrument python app.py
