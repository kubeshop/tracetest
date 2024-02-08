import {SupportedDataStores} from 'types/DataStore.types';
import {Model, TDataStoreSchemas} from 'types/Common.types';

export type TRawDataStore = TDataStoreSchemas['DataStoreResource'];
export type TRawAzureAppInsightsDataStore = TDataStoreSchemas['AzureAppInsights'];

export type TRawOtlpDataStore = {};
type DataStore = Model<TRawDataStore, {}>['spec'] & {
  agent?: TRawOtlpDataStore;
  azureappinsights?: TRawAzureAppInsightsDataStore & TRawOtlpDataStore;
  datadog?: TRawOtlpDataStore;
  dynatrace?: TRawOtlpDataStore;
  honeycomb?: TRawOtlpDataStore;
  instana?: TRawOtlpDataStore;
  lightstep?: TRawOtlpDataStore;
  newrelic?: TRawOtlpDataStore;
  otlp?: TRawOtlpDataStore;
  signoz?: TRawOtlpDataStore;
};

const DataStore = ({
  spec: {
    id = '',
    name = '',
    type = SupportedDataStores.JAEGER,
    default: isDefault = false,
    createdAt = '',
    opensearch = {},
    elasticapm = {},
    signalfx = {},
    jaeger = {},
    tempo = {},
    awsxray = {},
    azureappinsights = {},
    sumologic = {},
  } = {id: '', name: '', type: SupportedDataStores.JAEGER},
}: TRawDataStore): DataStore => ({
  id,
  name,
  type: (type as string) === 'agent' ? 'otlp' : type,
  default: isDefault,
  createdAt,
  opensearch,
  elasticapm,
  signalfx,
  jaeger,
  tempo,
  awsxray,
  azureappinsights,
  sumologic,
});

export const fromType = (type: SupportedDataStores): DataStore =>
  DataStore({
    spec: {
      type,
      name: 'current',
      id: 'current',
    },
  });

export default DataStore;
