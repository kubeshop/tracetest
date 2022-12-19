import {SupportedDataStores, TDataStore, TRawDataStore} from 'types/Config.types';

const DataStore = ({
  id = '',
  name = '',
  type = SupportedDataStores.JAEGER,
  isDefault = false,
  openSearch = {},
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
  jaeger,
  tempo,
  createdAt,
});

export default DataStore;
