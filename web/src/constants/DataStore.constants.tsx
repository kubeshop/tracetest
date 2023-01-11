import {SupportedDataStores} from '../types/Config.types';

export const SupportedDataStoresToName = {
  [SupportedDataStores.JAEGER]: 'Jaeger',
  [SupportedDataStores.OpenSearch]: 'OpenSearch',
  [SupportedDataStores.SignalFX]: 'SignalFX',
  [SupportedDataStores.TEMPO]: 'Tempo',
  [SupportedDataStores.OtelCollector]: 'OpenTelemetry',
  [SupportedDataStores.ElasticApm]: 'Elastic APM',
} as const;

export const SupportedDataStoresToDocsLink = {
  [SupportedDataStores.JAEGER]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/jaeger',
  [SupportedDataStores.OpenSearch]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/opensearch',
  [SupportedDataStores.ElasticApm]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/elasticapm',
  [SupportedDataStores.SignalFX]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/signalfx',
  [SupportedDataStores.TEMPO]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/tempo',
  [SupportedDataStores.OtelCollector]:
    'https://docs.tracetest.io/configuration/connecting-to-data-stores/opentelemetry-collector',
} as const;

export const SupportedDataStoresToExplanation: Record<string, React.ReactElement> = {
  [SupportedDataStores.OtelCollector]: (
    <>
      Tracetest can work with any distributed tracing solution that is utilizing the{' '}
      <a href="https://opentelemetry.io/docs/collector/" target="_blank">
        OpenTelemetry Collector
      </a>{' '}
      via a second pipeline. The second pipeline enables your current tracing system to send only Tracetest spans to
      Tracetest, while all other spans continue to go to the backend of your choice.
    </>
  ),
};

export const OtelCollectorConfigSample = `receivers:
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
