import {ConfigMode, SupportedDataStores} from 'types/DataStore.types';
import DataStore, {TRawDataStore} from './DataStore.model';

type DataStoreConfig = {
  defaultDataStore: DataStore;
  mode: ConfigMode;
};

const DataStoreConfig = (dataStores: TRawDataStore[] = []): DataStoreConfig => {
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
