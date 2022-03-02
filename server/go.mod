module github.com/kubeshop/tracetest/server

go 1.13

require (
	github.com/gogo/protobuf v1.3.2
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.4
	github.com/mitchellh/mapstructure v1.4.3
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/collector v0.44.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.28.0
	go.opentelemetry.io/otel v1.4.1
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.3.0
	go.opentelemetry.io/otel/sdk v1.4.1
	go.opentelemetry.io/otel/trace v1.4.1
	go.opentelemetry.io/proto/otlp v0.12.0
	go.uber.org/goleak v1.1.12 // indirect
	google.golang.org/grpc v1.44.0
	gopkg.in/yaml.v2 v2.4.0
)
