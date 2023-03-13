dataStore:
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
    type: otlp{{ end}}
