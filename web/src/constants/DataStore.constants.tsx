import {SupportedDataStores} from '../types/Config.types';

export const SupportedDataStoresToName = {
  [SupportedDataStores.JAEGER]: 'Jaeger',
  [SupportedDataStores.OpenSearch]: 'OpenSearch',
  [SupportedDataStores.SignalFX]: 'SignalFX',
  [SupportedDataStores.TEMPO]: 'Tempo',
  [SupportedDataStores.OtelCollector]: 'OpenTelemetry',
  [SupportedDataStores.ElasticApm]: 'Elastic APM',
  [SupportedDataStores.NewRelic]: 'New Relic',
  [SupportedDataStores.Lightstep]: 'Lightstep',
  [SupportedDataStores.Datadog]: 'Datadog',
  [SupportedDataStores.AWSXRay]: 'AWS X-Ray',
} as const;

export const SupportedDataStoresToDocsLink = {
  [SupportedDataStores.JAEGER]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/jaeger',
  [SupportedDataStores.OpenSearch]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/opensearch',
  [SupportedDataStores.ElasticApm]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/elasticapm',
  [SupportedDataStores.NewRelic]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/new-relic',
  [SupportedDataStores.Lightstep]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/lightstep',
  [SupportedDataStores.Datadog]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/datadog',
  [SupportedDataStores.SignalFX]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/signalfx',
  [SupportedDataStores.TEMPO]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/tempo',
  [SupportedDataStores.AWSXRay]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/aws-x-ray',
  [SupportedDataStores.OtelCollector]:
    'https://docs.tracetest.io/configuration/connecting-to-data-stores/opentelemetry-collector',
} as const;

const collectorExplanation = (
  <>
    Tracetest can work with any distributed tracing solution that is utilizing the{' '}
    <a href="https://opentelemetry.io/docs/collector/" target="_blank">
      OpenTelemetry Collector
    </a>{' '}
    via a second pipeline. The second pipeline enables your current tracing system to send only Tracetest spans to
    Tracetest, while all other spans continue to go to the backend of your choice.
  </>
);

const aDotCollectorExplanation = (
  <>
    Tracetest can work with any AWS distributed tracing solution that is utilizing the{' '}
    <a href="https://aws-otel.github.io/docs/getting-started/collector" target="_blank">
      AWS OpenTelemetry Collector
    </a>{' '}
    via a second pipeline. The second pipeline enables your current tracing system to send only Tracetest spans to
    Tracetest, while all other spans continue to go to the backend of your choice.
  </>
);

export const SupportedDataStoresToExplanation: Record<string, React.ReactElement> = {
  [SupportedDataStores.OtelCollector]: collectorExplanation,
  [SupportedDataStores.NewRelic]: collectorExplanation,
  [SupportedDataStores.Lightstep]: collectorExplanation,
  [SupportedDataStores.Datadog]: collectorExplanation,
  [SupportedDataStores.AWSXRay]: aDotCollectorExplanation,
};

export const NoTestConnectionDataStoreList = [
  SupportedDataStores.OtelCollector,
  SupportedDataStores.Lightstep,
  SupportedDataStores.Datadog,
  SupportedDataStores.NewRelic,
  SupportedDataStores.AWSXRay,
];
