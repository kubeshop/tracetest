import {IMockFactory} from 'types/Common.types';
import {TDataStoreConfig, TRawDataStoreConfig} from 'types/Config.types';
import Config from '../DataStoreConfig.model';

const DataStoreConfigMock: IMockFactory<TDataStoreConfig, TRawDataStoreConfig> = () => ({
  raw(data = {}) {
    return {
      dataStores: [],
      defaultDataStore: '',
      ...data,
    };
  },
  model(data = {}) {
    return Config(this.raw(data));
  },
});

export default DataStoreConfigMock();
