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
  [SupportedDataStores.AzureAppInsights]: 'Azure App Insights',
  [SupportedDataStores.Signoz]: 'Signoz',
  [SupportedDataStores.Dynatrace]: 'Dynatrace',
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
  [SupportedDataStores.AzureAppInsights]:
    'https://docs.tracetest.io/configuration/connecting-to-data-stores/azure-app-insights',
  [SupportedDataStores.Signoz]:
    'https://docs.tracetest.io/configuration/connecting-to-data-stores/signoz',
  [SupportedDataStores.Dynatrace]:
    'https://docs.tracetest.io/configuration/connecting-to-data-stores/dynatrace',
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
  [SupportedDataStores.Signoz]: '',
  [SupportedDataStores.Dynatrace]: 'https://abc12345.live.dynatrace.com/api/v2/otlp',
} as const;
