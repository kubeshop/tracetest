import {SupportedDataStores} from 'types/Config.types';
import {Model, TDataStoreSchemas} from 'types/Common.types';

export type RawDataStore = TDataStoreSchemas['DataStore'];
type DataStore = Model<
  RawDataStore,
  {
    otlp?: {};
    newRelic?: {};
    lightstep?: {};
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
  createdAt = '',
}: RawDataStore): DataStore => ({
  id,
  name,
  type,
  isDefault,
  openSearch,
  signalFx,
  elasticApm,
  jaeger,
  tempo,
  createdAt,
});

export default DataStore;
