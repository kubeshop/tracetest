import {SupportedDataStores} from '../types/DataStore.types';

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
  [SupportedDataStores.Honeycomb]: 'Honeycomb',
  [SupportedDataStores.AzureAppInsights]: 'Azure Application Insights',
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
  [SupportedDataStores.Honeycomb]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/honeycomb',
  [SupportedDataStores.AzureAppInsights]: 'https://docs.tracetest.io/configuration/connecting-to-data-stores/azure-app-insights',
} as const;

export const SupportedDataStoresToDefaultEndpoint = {
  [SupportedDataStores.JAEGER]: 'jaeger:16685',
  [SupportedDataStores.OpenSearch]: 'http://opensearch:9200',
  [SupportedDataStores.SignalFX]: '',
  [SupportedDataStores.TEMPO]: 'tempo:9095',
  [SupportedDataStores.OtelCollector]: '',
  [SupportedDataStores.ElasticApm]: 'http://elasticsearch:9200',
  [SupportedDataStores.NewRelic]: '',
  [SupportedDataStores.Lightstep]: '',
  [SupportedDataStores.Datadog]: '',
  [SupportedDataStores.AWSXRay]: '',
  [SupportedDataStores.Honeycomb]: '',
  [SupportedDataStores.AzureAppInsights]: '',
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

export const SupportedDataStoresToExplanation: Record<string, React.ReactElement> = {
  [SupportedDataStores.OtelCollector]: collectorExplanation,
  [SupportedDataStores.NewRelic]: collectorExplanation,
  [SupportedDataStores.Lightstep]: collectorExplanation,
  [SupportedDataStores.Datadog]: collectorExplanation,
};

export const NoTestConnectionDataStoreList = [
  SupportedDataStores.OtelCollector,
  SupportedDataStores.Lightstep,
  SupportedDataStores.Datadog,
  SupportedDataStores.NewRelic,
  SupportedDataStores.Honeycomb,
];
