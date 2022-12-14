import {SupportedDataStores, TDataStoreService} from 'types/Config.types';

const OtelCollectorService = (): TDataStoreService => ({
  getRequest() {
    return Promise.resolve({
      type: SupportedDataStores.OtelCollector,
    });
  },
  validateDraft() {
    return Promise.resolve(true);
  },
  getInitialValues() {
    return {
      dataStore: {},
      dataStoreType: SupportedDataStores.OtelCollector,
    };
  },
});

export default OtelCollectorService();
