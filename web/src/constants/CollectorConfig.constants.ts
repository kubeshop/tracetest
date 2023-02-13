import {SupportedDataStores} from '../types/Config.types';

export const Lightstep = `receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:
    timeout: 100ms

exporters:
  logging:
    logLevel: debug
  # OTLP for Tracetest
  otlp/tt:
    endpoint: tracetest:21321 # Send traces to Tracetest. Read more in docs here:  https://docs.tracetest.io/configuration/connecting-to-data-stores/opentelemetry-collector
    tls:
      insecure: true
  # OTLP for Lightstep
  otlp/ls:
    endpoint: ingest.lightstep.com:443
    headers:
      "lightstep-access-token": "<lightstep_access_token>" # Send traces to Lightstep. Read more in docs here: https://docs.lightstep.com/otel/otel-quick-start

service:
  pipelines:
    traces/tt:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/tt]
    traces/ls:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, otlp/ls]
`;

export const OtelCollector = `receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:
    timeout: 100ms

exporters:
  otlp/1:
    endpoint: tracetest:21321
    tls:
      insecure: true

service:
  pipelines:
    traces/1:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/1]
`;

export const NewRelic = `receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:
    timeout: 100ms

exporters:
  logging:
    logLevel: debug
  # OTLP for Tracetest
  otlp/tt:
    endpoint: tracetest:21321 # Send traces to Tracetest. Read more in docs here:  https://docs.tracetest.io/configuration/connecting-to-data-stores/opentelemetry-collector
    tls:
      insecure: true
  # OTLP for New Relic
  otlp/nr:
    endpoint: otlp.nr-data.net:443
    headers:
      api-key: <new_relic_ingest_licence_key> # Send traces to New Relic.
      # Read more in docs here: https://docs.newrelic.com/docs/more-integrations/open-source-telemetry-integrations/opentelemetry/opentelemetry-setup/#collector-export
      # And here: https://docs.newrelic.com/docs/more-integrations/open-source-telemetry-integrations/opentelemetry/collector/opentelemetry-collector-basic/

service:
  pipelines:
    traces/tt:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/tt]
    traces/nr:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, otlp/nr]
`;

export const Datadog = `receivers:
  otlp:
    protocols:
      http:
      grpc:

processors:
  batch:
    send_batch_max_size: 100
    send_batch_size: 10
    timeout: 10s

exporters:
  # OTLP for Tracetest
  otlp/tt:
    endpoint: tracetest:21321 # Send traces to Tracetest. Read more in docs here:  https://docs.tracetest.io/configuration/connecting-to-data-stores/opentelemetry-collector
    tls:
      insecure: true
  # Datadog exporter
  datadog:
    api:
      site: datadoghq.com
      key: <datadog_API_key> # Add here you API key for Datadog
      # Read more in docs here: https://docs.datadoghq.com/opentelemetry/otel_collector_datadog_exporter
service:
  pipelines:
    traces/tt:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/tt]
    traces/dd:
      receivers: [otlp]
      processors: [batch]
      exporters: [datadog]
`;

export const AWSXRay = `receivers:
  awsxray:
    transport: udp

  processors:
  batch:

exporters:
  logging:
    loglevel: debug
  awsxray:
    region: <aws-region>
  otlp/tt:
    endpoint: tracetest:21321
    tls:
      insecure: true

service:
  pipelines:
    traces/tt:
      receivers: [awsxray]
      processors: [batch]
      exporters: [otlp/tt]
    traces/xr:
      receivers: [awsxray]
      exporters: [awsxray]
`;

export const CollectorConfigMap = {
  [SupportedDataStores.Datadog]: Datadog,
  [SupportedDataStores.Lightstep]: Lightstep,
  [SupportedDataStores.NewRelic]: NewRelic,
  [SupportedDataStores.OtelCollector]: OtelCollector,
  [SupportedDataStores.AWSXRay]: AWSXRay,
} as const;
