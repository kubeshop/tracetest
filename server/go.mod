module github.com/kubeshop/tracetest/server

go 1.18

replace k8s.io/client-go => k8s.io/client-go v0.18.0

require (
	github.com/alecthomas/participle/v2 v2.0.0-alpha8
	github.com/denisbrodbeck/machineid v1.0.1
	github.com/fullstorydev/grpcurl v1.8.6
	github.com/gogo/protobuf v1.3.2
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
	github.com/j2gg0s/otsql v0.14.0
	github.com/jhump/protoreflect v1.12.0
	github.com/lib/pq v1.10.5
	github.com/mitchellh/mapstructure v1.4.3
	github.com/orlangure/gnomock v0.20.0
	github.com/prometheus/prometheus v1.8.2-0.20211217191541-41f1a8125e66
	github.com/segmentio/analytics-go/v3 v3.2.1
	github.com/stretchr/testify v1.7.1
	go.opentelemetry.io/collector v0.44.0
	go.opentelemetry.io/contrib/propagators/aws v1.5.0
	go.opentelemetry.io/contrib/propagators/b3 v1.5.0
	go.opentelemetry.io/contrib/propagators/jaeger v1.5.0
	go.opentelemetry.io/contrib/propagators/ot v1.5.0
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.6.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.3.0
	go.opentelemetry.io/otel/sdk v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
	go.opentelemetry.io/proto/otlp v0.12.0
	google.golang.org/grpc v1.47.0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/containerd/containerd v1.6.6 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.14+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-redis/redis/v7 v7.4.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/klauspost/compress v1.15.6 // indirect
	github.com/knadh/koanf v1.4.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/mostynb/go-grpc-compression v1.1.16 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20211202183452-c5a74bcca799 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/segmentio/backo-go v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/streadway/amqp v1.0.0 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	go.opentelemetry.io/collector/model v0.44.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.28.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.6.0 // indirect
	go.opentelemetry.io/otel/metric v0.26.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/net v0.0.0-20220615171555-694bf12d69de // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/sys v0.0.0-20220614162138-6c1b26c55098 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220615141314-f1464d18c36b // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

// Temporary fix until we manage to merge the patch to the gnomock repo (https://github.com/orlangure/gnomock/pull/534)
replace github.com/orlangure/gnomock v0.20.0 => github.com/mathnogueira/gnomock v0.21.0

replace github.com/orlangure/gnomock/preset/postgres v0.20.0 => github.com/mathnogueira/gnomock/preset/postgres v0.21.0

replace github.com/orlangure/gnomock/preset/redis v0.20.0 => github.com/mathnogueira/gnomock/preset/redis v0.21.0

replace github.com/orlangure/gnomock/preset/rabbitmq v0.20.0 => github.com/mathnogueira/gnomock/preset/rabbitmq v0.21.0
