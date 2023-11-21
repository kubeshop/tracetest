import {SupportedDataStores} from 'types/DataStore.types';

export const tracetest = `# OTLP for Tracetest
  otlp/tracetest:
    endpoint: tracetest:4317 # Send traces to Tracetest. Read more in docs here:  https://docs.tracetest.io/configuration/connecting-to-data-stores/opentelemetry-collector
    tls:
      insecure: true`;

export const Lightstep = (traceTestBlock: string) => `receivers:
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

  ${traceTestBlock}

  # OTLP for Lightstep
  otlp/lightstep:
    endpoint: ingest.lightstep.com:443
    headers:
      "lightstep-access-token": "<lightstep_access_token>" # Send traces to Lightstep. Read more in docs here: https://docs.lightstep.com/otel/otel-quick-start

service:
  pipelines:
    traces/tracetest:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/tracetest]
    traces/lightstep:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, otlp/lightstep]
`;

export const OtelCollector = (traceTestBlock: string) => `receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:
    timeout: 100ms

exporters:
  ${traceTestBlock}

service:
  pipelines:
    traces/1:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/tracetest]
`;

export const NewRelic = (traceTestBlock: string) => `receivers:
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

  ${traceTestBlock}

  # OTLP for New Relic
  otlp/newrelic:
    endpoint: otlp.nr-data.net:443
    headers:
      api-key: <new_relic_ingest_licence_key> # Send traces to New Relic.
      # Read more in docs here: https://docs.newrelic.com/docs/more-integrations/open-source-telemetry-integrations/opentelemetry/opentelemetry-setup/#collector-export
      # And here: https://docs.newrelic.com/docs/more-integrations/open-source-telemetry-integrations/opentelemetry/collector/opentelemetry-collector-basic/

service:
  pipelines:
    traces/tracetest:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/tracetest]
    traces/newrelic:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, otlp/newrelic]
`;

export const Datadog = (traceTestBlock: string) => `receivers:
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
  ${traceTestBlock}

  # Datadog exporter
  datadog:
    api:
      site: datadoghq.com
      key: <datadog_API_key> # Add here you API key for Datadog
      # Read more in docs here: https://docs.datadoghq.com/opentelemetry/otel_collector_datadog_exporter
service:
  pipelines:
    traces/tracetest:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/tracetest]
    traces/datadog:
      receivers: [otlp]
      processors: [batch]
      exporters: [datadog]
`;

export const Honeycomb = (traceTestBlock: string) => `receivers:
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

  ${traceTestBlock}

  # OTLP for Honeycomb
  otlp/honeycomb:
    endpoint: "api.honeycomb.io:443"
    headers:
      "x-honeycomb-team": "YOUR_API_KEY"
      # Read more in docs here: https://docs.honeycomb.io/getting-data-in/otel-collector/

service:
  pipelines:
    traces/tracetest:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/tracetest]
    traces/honeycomb:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, otlp/honeycomb]
`;

export const AzureAppInsights = (traceTestBlock: string) => `receivers:
otlp:
  protocols:
    grpc:
    http:

processors:
  batch:

exporters:
  azuremonitor:
    instrumentation_key: <your-instrumentation-key>

  ${traceTestBlock}

service:
  pipelines:
    traces/tracetest:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/tracetest]
    traces/appinsights:
      receivers: [otlp]
      exporters: [azuremonitor]
`;

export const Signoz = (traceTestBlock: string) => `receivers:
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

  ${traceTestBlock}

  # OTLP for Signoz
  otlp/signoz:
    endpoint: address-to-your-signoz-server:4317 # Send traces to Signoz. Read more in docs here: https://signoz.io/docs/tutorial/opentelemetry-binary-usage-in-virtual-machine/#opentelemetry-collector-configuration
    tls:
      insecure: true

service:
  pipelines:
    traces/tracetest:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, otlp/tracetest]
    traces/signoz:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, otlp/signoz]
`;

export const Dynatrace = (traceTestBlock: string) => `receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:
    timeout: 100ms

exporters:
  logging:
    verbosity: detailed

  ${traceTestBlock}

  # OTLP for Dynatrace
  otlphttp/dynatrace:
    endpoint: https://abc12345.live.dynatrace.com/api/v2/otlp # Send traces to Dynatrace. Read more in docs here: https://www.dynatrace.com/support/help/extend-dynatrace/opentelemetry/collector#configuration
    headers:
      Authorization: "Api-Token dt0c01.sample.secret"  # Requires "openTelemetryTrace.ingest" permission

service:
  pipelines:
    traces/tracetest:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, otlp/tracetest]
    traces/dynatrace:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, otlphttp/dynatrace]
`;

export const CollectorConfigMap = {
  [SupportedDataStores.Datadog]: Datadog(tracetest),
  [SupportedDataStores.Lightstep]: Lightstep(tracetest),
  [SupportedDataStores.NewRelic]: NewRelic(tracetest),
  [SupportedDataStores.OtelCollector]: OtelCollector(tracetest),
  [SupportedDataStores.Honeycomb]: Honeycomb(tracetest),
  [SupportedDataStores.AzureAppInsights]: AzureAppInsights(tracetest),
  [SupportedDataStores.Signoz]: Signoz(tracetest),
  [SupportedDataStores.Dynatrace]: Dynatrace(tracetest),
  [SupportedDataStores.Agent]: OtelCollector(tracetest),
} as const;

export const CollectorConfigFunctionMap = {
  [SupportedDataStores.Datadog]: Datadog,
  [SupportedDataStores.Lightstep]: Lightstep,
  [SupportedDataStores.NewRelic]: NewRelic,
  [SupportedDataStores.OtelCollector]: OtelCollector,
  [SupportedDataStores.Honeycomb]: Honeycomb,
  [SupportedDataStores.AzureAppInsights]: AzureAppInsights,
  [SupportedDataStores.Signoz]: Signoz,
  [SupportedDataStores.Dynatrace]: Dynatrace,
  [SupportedDataStores.Agent]: OtelCollector,
} as const;
