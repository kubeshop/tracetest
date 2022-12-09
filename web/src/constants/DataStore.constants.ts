import {SupportedDataStores} from '../types/Config.types';

export const SupportedDataStoresToName = {
  [SupportedDataStores.JAEGER]: 'Jaeger',
  [SupportedDataStores.OpenSearch]: 'OpenSearch',
  [SupportedDataStores.SignalFX]: 'SignalFX',
  [SupportedDataStores.TEMPO]: 'Tempo',
} as const;

export const SupportedDataStoresToDocsLink = {
  [SupportedDataStores.JAEGER]: 'https://docs.tracetest.io/run-locally/#installing-jaeger',
  [SupportedDataStores.OpenSearch]: 'https://docs.tracetest.io/run-locally/#installing-jaeger',
  [SupportedDataStores.SignalFX]: 'https://docs.tracetest.io/run-locally/#installing-jaeger',
  [SupportedDataStores.TEMPO]: 'https://docs.tracetest.io/run-locally/#installing-jaeger',
} as const;
