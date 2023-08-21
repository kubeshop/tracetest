{{ if eq .installBackend "true" }}---
type: DataStore
spec:
  name: {{ .backendType }}
  type: {{ .backendType }}{{ if eq .backendType "jaeger" }}
  jaeger:
    type: jaeger
    jaeger:
      endpoint: {{ .backendEndpoint }}
      tls:
        insecure: {{ .backendInsecure }}{{ end}}{{ if eq .backendType "tempo" }}
  tempo:
    type: tempo
    tempo:
      type: grpc
      grpc:
        endpoint: {{ .backendEndpoint }}
        tls:
          insecure: {{ .backendInsecure }}{{ end}}{{ if eq .backendType "opensearch" }}
  opensearch:
    type: opensearch
    opensearch:
      addresses: {{ .backendAddresses }}
      index: {{ .backendIndex }}{{ end}}{{ if eq .backendType "signalfx" }}
  signalfx:
    type: signalfx
    signalfx:
      token: {{ .backendToken }}
      realm: {{ .backendRealm }}{{ end}}{{ if eq .backendType "otlp" }}
  otlp:
    type: otlp{{ end}}{{ end}}
---
type: Config
spec:
  analyticsEnabled: {{ .analyticsEnabled }}
---
type: PollingProfile
spec:
  name: Custom Profile
  strategy: periodic
  default: true
  periodic:
    timeout: 2m
    retryDelay: 3s
{{ if eq .enablePokeshopDemo "true" }}---
type: Demo
spec:
  name: pokeshop
  type: pokeshop
  enabled: true
  pokeshop:
    httpEndpoint: {{ .pokeshopHttp }}
    grpcEndpoint: {{ .pokeshopGrpc }}
    kafkaBroker: {{ .pokeshopKafka }}{{end}}{{ if eq .enableOtelDemo "true" }}
---
type: Demo
spec:
  name: otel
  type: otelstore
  enabled: true
  opentelemetryStore:
    frontendEndpoint: {{ .otelFrontend }}
    productCatalogEndpoint: {{ .otelProductCatalog }}
    cartEndpoint: {{ .otelCart }}
    checkoutEndpoint: {{ .otelCheckout }}{{end}}
