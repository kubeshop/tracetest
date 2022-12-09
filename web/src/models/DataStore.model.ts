import {SupportedDataStores, TDataStore, TRawDataStore} from 'types/Config.types';

const DataStore = ({
  name = '',
  type = SupportedDataStores.JAEGER,
  openSearch = {},
  signalFx = {},
  jaeger = {},
  tempo = {},
}: TRawDataStore): TDataStore => ({
  name,
  type,
  openSearch,
  signalFx,
  jaeger,
  tempo,
});

export default DataStore;
