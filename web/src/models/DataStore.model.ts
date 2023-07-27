import {SupportedDataStores} from 'types/DataStore.types';
import {Model, TDataStoreSchemas} from 'types/Common.types';

export type TRawDataStore = TDataStoreSchemas['DataStoreResource'];
export type TRawAzureAppInsightsDataStore = TDataStoreSchemas['AzureAppInsights'];

export type TRawOtlpDataStore = {isIngestorEnabled?: boolean};
type DataStore = Model<TRawDataStore, {}>['spec'] & {
  otlp?: TRawOtlpDataStore;
  newrelic?: TRawOtlpDataStore;
  lightstep?: TRawOtlpDataStore;
  datadog?: TRawOtlpDataStore;
  honeycomb?: TRawOtlpDataStore;
  azureappinsights?: TRawAzureAppInsightsDataStore & TRawOtlpDataStore;
  signoz?: TRawOtlpDataStore;
  dynatrace?: TRawOtlpDataStore;
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
  } = {id: '', name: '', type: SupportedDataStores.JAEGER},
}: TRawDataStore): DataStore => ({
  id,
  name,
  type,
  default: isDefault,
  createdAt,
  opensearch,
  elasticapm,
  signalfx,
  jaeger,
  tempo,
  awsxray,
  azureappinsights,
});

export default DataStore;
