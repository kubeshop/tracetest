import {SupportedDataStores, TDataStoreService} from 'types/DataStore.types';

const OtelCollectorService = (): TDataStoreService => ({
  getRequest(draft, dataStoreType = SupportedDataStores.OtelCollector) {
    return Promise.resolve({
      type: dataStoreType,
      name: dataStoreType,
    });
  },
  validateDraft() {
    return Promise.resolve(true);
  },
  getInitialValues(draft, dataStoreType = SupportedDataStores.OtelCollector) {
    return {
      dataStore: {
        name: dataStoreType,
        type: dataStoreType,
        [dataStoreType]: {},
      },
      dataStoreType,
    };
  },
  getIsOtlpBased() {
    return true;
  },
  getPublicInfo() {
    return {};
  },
});

export default OtelCollectorService();
