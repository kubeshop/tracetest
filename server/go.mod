module github.com/kubeshop/tracetest

go 1.13

replace k8s.io/client-go => k8s.io/client-go v0.18.0

require (
	github.com/alecthomas/participle/v2 v2.0.0-alpha8
	github.com/denisbrodbeck/machineid v1.0.1
	github.com/gogo/protobuf v1.3.2
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
	github.com/j2gg0s/otsql v0.14.0
	github.com/lib/pq v1.10.5
	github.com/mitchellh/mapstructure v1.4.3
	github.com/orlangure/gnomock v0.20.0
	github.com/prometheus/prometheus v1.8.2-0.20211217191541-41f1a8125e66
	github.com/stretchr/testify v1.7.1
	go.opentelemetry.io/collector v0.44.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.29.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.28.0
	go.opentelemetry.io/contrib/propagators/aws v1.5.0
	go.opentelemetry.io/contrib/propagators/b3 v1.5.0
	go.opentelemetry.io/contrib/propagators/jaeger v1.5.0
	go.opentelemetry.io/contrib/propagators/ot v1.5.0
	go.opentelemetry.io/otel v1.5.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.5.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.5.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.3.0
	go.opentelemetry.io/otel/sdk v1.5.0
	go.opentelemetry.io/otel/trace v1.5.0
	go.opentelemetry.io/proto/otlp v0.12.0
	google.golang.org/grpc v1.45.0
	gopkg.in/yaml.v2 v2.4.0
)

// Temporary fix until we manage to merge the patch to the gnomock repo (https://github.com/orlangure/gnomock/pull/534)
replace github.com/orlangure/gnomock v0.20.0 => github.com/mathnogueira/gnomock v0.21.0

replace github.com/orlangure/gnomock/preset/postgres v0.20.0 => github.com/mathnogueira/gnomock/preset/postgres v0.21.0

replace github.com/orlangure/gnomock/preset/redis v0.20.0 => github.com/mathnogueira/gnomock/preset/redis v0.21.0

replace github.com/orlangure/gnomock/preset/rabbitmq v0.20.0 => github.com/mathnogueira/gnomock/preset/rabbitmq v0.21.0
