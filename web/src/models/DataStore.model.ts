import {SupportedDataStores} from 'types/DataStore.types';
import {Model, TDataStoreSchemas} from 'types/Common.types';

export type TRawDataStore = TDataStoreSchemas['DataStoreResource'];
export type TRawAzureAppInsightsDataStore = TDataStoreSchemas['AzureAppInsights'];

export type TRawOtlpDataStore = {};
type DataStore = Model<TRawDataStore, {}>['spec'] & {
  otlp?: TRawOtlpDataStore;
  newrelic?: TRawOtlpDataStore;
  lightstep?: TRawOtlpDataStore;
  datadog?: TRawOtlpDataStore;
  honeycomb?: TRawOtlpDataStore;
  azureappinsights?: TRawAzureAppInsightsDataStore & TRawOtlpDataStore;
  signoz?: TRawOtlpDataStore;
  dynatrace?: TRawOtlpDataStore;
  agent?: TRawOtlpDataStore;
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

export default DataStore;
