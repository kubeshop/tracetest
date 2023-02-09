postgresConnString: "{{ .psql }}"

poolingConfig:
  maxWaitTimeForTrace: 2m
  retryDelay: 3s

googleAnalytics:
  enabled: {{ .analyticsEnabled }}

demo:
  enabled: [{{ .enabledDemos }}]
  endpoints:
    pokeshopHttp: {{ .pokeshopHttp }}
    pokeshopGrpc: {{ .pokeshopGrpc }}
    otelFrontend: {{ .otelFrontend }}
    otelProductCatalog: {{ .otelProductCatalog }}
    otelCart: {{ .otelCart }}
    otelCheckout: {{ .otelCheckout }}

experimentalFeatures: []
{{ if eq .installBackend "true" }}
telemetry:
  dataStores:{{ if eq .backendType "jaeger" }}
    jaeger:
      type: jaeger
      jaeger:
        endpoint: {{ .backendEndpoint }}
        tls:
          insecure: {{ .backendInsecure }}{{ end}}{{ if eq .backendType "tempo" }}
    tempo:
      type: tempo
      tempo:
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
      type: otlp{{ end}}

server:
  telemetry:
    dataStore: {{ .backendType }}
{{end}}
