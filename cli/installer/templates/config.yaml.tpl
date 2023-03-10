postgres:
  host: {{ .pHost }}
  user: {{ .pUser }}
  password: {{ .pPasswd }}
  port: 5432
  dbname: postgres
  params: sslmode=disable

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
