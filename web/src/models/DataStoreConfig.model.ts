import {ConfigMode, TDataStoreConfig, TRawDataStoreConfig} from 'types/Config.types';
import DataStore from './DataStore.model';

const DataStoreConfig = ({dataStores = [], defaultDataStore = ''}: TRawDataStoreConfig): TDataStoreConfig => {
  const dataStoreList = dataStores.map(rawDataStore => DataStore(rawDataStore));

  const mode =
    (Boolean(defaultDataStore && dataStoreList.find(({name}) => name === defaultDataStore)) && ConfigMode.READY) ||
    ConfigMode.NO_TRACING_MODE;

  return {
    dataStores: dataStoreList,
    defaultDataStore,
    mode,
  };
};

export default DataStoreConfig;
