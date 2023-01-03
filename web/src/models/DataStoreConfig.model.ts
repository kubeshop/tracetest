import {ConfigMode, SupportedDataStores, TDataStoreConfig, TRawDataStore} from 'types/Config.types';
import DataStore from './DataStore.model';

const DataStoreConfig = (dataStores: TRawDataStore[] = []): TDataStoreConfig => {
  const dataStoreList = dataStores.map(rawDataStore => DataStore(rawDataStore));
  const defaultDataStore = dataStoreList.find(({isDefault}) => isDefault);
  const mode = (!!defaultDataStore && ConfigMode.READY) || ConfigMode.NO_TRACING_MODE;

  return {
    defaultDataStore:
      defaultDataStore ??
      DataStore({
        name: 'default',
        type: SupportedDataStores.JAEGER,
      }),
    mode,
  };
};

export default DataStoreConfig;
