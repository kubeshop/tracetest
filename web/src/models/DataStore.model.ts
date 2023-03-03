import {SupportedDataStores} from 'types/DataStore.types';
import {Model, TDataStoreSchemas} from 'types/Common.types';

export type TRawDataStore = TDataStoreSchemas['DataStore'];
type DataStore = Model<
  TRawDataStore,
  {
    otlp?: {};
    newRelic?: {};
    lightstep?: {};
    datadog?: {};
  }
>;

const DataStore = ({
  id = '',
  name = '',
  type = SupportedDataStores.JAEGER,
  isDefault = false,
  openSearch = {},
  elasticApm = {},
  signalFx = {},
  jaeger = {},
  tempo = {},
  awsxray = {},
  createdAt = '',
}: TRawDataStore): DataStore => ({
  id,
  name,
  type,
  isDefault,
  openSearch,
  signalFx,
  elasticApm,
  jaeger,
  tempo,
  awsxray,
  createdAt,
});

export default DataStore;
