import {IMockFactory} from 'types/Common.types';
import {TDataStoreConfig, TRawDataStore} from 'types/Config.types';
import Config from '../DataStoreConfig.model';

const DataStoreConfigMock: IMockFactory<
  TDataStoreConfig,
  {
    dataStores: TRawDataStore[];
  }
> = () => ({
  raw({dataStores = []} = {}) {
    return {dataStores};
  },
  model({dataStores = []} = {}) {
    return Config(this.raw({dataStores}).dataStores);
  },
});

export default DataStoreConfigMock();
