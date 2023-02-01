import {IMockFactory} from 'types/Common.types';
import {TRawDataStore} from '../DataStore.model';
import DataStoreConfig from '../DataStoreConfig.model';

const DataStoreConfigMock: IMockFactory<
  DataStoreConfig,
  {
    dataStores: TRawDataStore[];
  }
> = () => ({
  raw({dataStores = []} = {}) {
    return {dataStores};
  },
  model({dataStores = []} = {}) {
    return DataStoreConfig(this.raw({dataStores}).dataStores);
  },
});

export default DataStoreConfigMock();
