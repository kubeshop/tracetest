import {SupportedDataStores, TDataStore, TRawDataStore} from 'types/Config.types';

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
}: TRawDataStore): TDataStore => ({
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
