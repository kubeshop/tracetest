import {ConfigMode} from 'types/DataStore.types';
import DataStore, {TRawDataStore} from './DataStore.model';

type DataStoreConfig = {
  defaultDataStore: DataStore;
  mode: ConfigMode;
};

const DataStoreConfig = (rawDataStore: TRawDataStore): DataStoreConfig => {
  const defaultDataStore = DataStore(rawDataStore);
  const isDefaultDataStore = defaultDataStore.default;
  const mode = isDefaultDataStore ? ConfigMode.READY : ConfigMode.NO_TRACING_MODE;

  return {
    defaultDataStore,
    mode,
  };
};

export default DataStoreConfig;
