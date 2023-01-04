import {SupportedDataStores, TDataStoreService} from 'types/Config.types';

const OtelCollectorService = (): TDataStoreService => ({
  getRequest() {
    return Promise.resolve({
      type: SupportedDataStores.OtelCollector,
      name: SupportedDataStores.OtelCollector,
    });
  },
  validateDraft() {
    return Promise.resolve(true);
  },
  getInitialValues() {
    return {
      dataStore: {
        name: SupportedDataStores.OtelCollector,
        type: SupportedDataStores.OtelCollector,
      },
      dataStoreType: SupportedDataStores.OtelCollector,
    };
  },
});

export default OtelCollectorService();
