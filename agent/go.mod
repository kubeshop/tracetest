module github.com/kubeshop/tracetest/agent

go 1.20

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.7.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.3.0
	github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery v1.1.0
	github.com/aws/aws-sdk-go v1.44.319
	github.com/elastic/go-elasticsearch/v8 v8.9.0
	github.com/fluidtruck/deepcopy v1.0.0
	github.com/fullstorydev/grpcurl v1.8.7
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.3
	github.com/goware/urlx v0.3.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/jhump/protoreflect v1.15.1
	github.com/kubeshop/tracetest/server v0.0.0-20230809150857-6314696222d3
	github.com/opensearch-project/opensearch-go v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/prometheus v1.8.2-0.20211217191541-41f1a8125e66
	github.com/spf13/cobra v1.7.0
	github.com/spf13/viper v1.16.0
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/collector/component v0.82.0
	go.opentelemetry.io/collector/config/configcompression v0.82.0
	go.opentelemetry.io/collector/config/configgrpc v0.82.0
	go.opentelemetry.io/collector/config/configtls v0.82.0
	go.opentelemetry.io/collector/semconv v0.82.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.42.1-0.20230612162650-64be7e574a17
	go.opentelemetry.io/contrib/propagators/aws v1.17.0
	go.opentelemetry.io/contrib/propagators/b3 v1.17.0
	go.opentelemetry.io/contrib/propagators/jaeger v1.17.0
	go.opentelemetry.io/contrib/propagators/ot v1.17.0
	go.opentelemetry.io/otel v1.16.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.16.0
	go.opentelemetry.io/otel/sdk v1.16.0
	go.opentelemetry.io/otel/trace v1.16.0
	go.opentelemetry.io/proto/otlp v1.0.0
	google.golang.org/grpc v1.56.2
	google.golang.org/protobuf v1.31.0
)

require (
	cloud.google.com/go/compute v1.20.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.3.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.0.0 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/alecthomas/participle/v2 v2.0.0-alpha8 // indirect
	github.com/brianvoe/gofakeit/v6 v6.17.0 // indirect
	github.com/bufbuild/protocompile v0.4.0 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/elastic/elastic-transport-go/v8 v8.0.0-20230329154755-1a3c63de0db6 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/goccy/go-yaml v1.11.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.2 // indirect
	github.com/jackc/puddle/v2 v2.2.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/knadh/koanf v1.5.0 // indirect
	github.com/knadh/koanf/v2 v2.0.1 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.1-0.20220423185008-bf980b35cac4 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostynb/go-grpc-compression v1.2.0 // indirect
	github.com/ohler55/ojg v1.14.4 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/segmentio/analytics-go/v3 v3.2.1 // indirect
	github.com/segmentio/backo-go v1.0.0 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/teris-io/shortid v0.0.0-20220617161101-71ec9f2aa569 // indirect
	go.opentelemetry.io/collector v0.82.0 // indirect
	go.opentelemetry.io/collector/config/configauth v0.82.0 // indirect
	go.opentelemetry.io/collector/config/confignet v0.82.0 // indirect
	go.opentelemetry.io/collector/config/configopaque v0.82.0 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.82.0 // indirect
	go.opentelemetry.io/collector/config/internal v0.82.0 // indirect
	go.opentelemetry.io/collector/confmap v0.82.0 // indirect
	go.opentelemetry.io/collector/extension v0.82.0 // indirect
	go.opentelemetry.io/collector/extension/auth v0.82.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.0.0-rcv0014 // indirect
	go.opentelemetry.io/collector/pdata v1.0.0-rcv0014 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.16.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/exp v0.0.0-20221217163422-3c43f8badb15 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/genproto v0.0.0-20230717213848-3f92550aa753 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230706204954-ccb25ca9f130 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230706204954-ccb25ca9f130 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
