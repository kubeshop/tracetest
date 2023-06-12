import {SupportedDataStores, TDataStoreService} from 'types/DataStore.types';

const OtelCollectorService = (): TDataStoreService => ({
  getRequest(draft, dataStoreType = SupportedDataStores.OtelCollector) {
    return Promise.resolve({
      type: dataStoreType,
      name: dataStoreType,
    });
  },
  validateDraft({dataStore: {isIngestorEnabled = false} = {}}) {
    return Promise.resolve(isIngestorEnabled);
  },
  getInitialValues(
    draft,
    dataStoreType = SupportedDataStores.OtelCollector,
    configuredDataStore = SupportedDataStores.OtelCollector
  ) {
    return {
      dataStore: {
        name: dataStoreType,
        type: dataStoreType,
        isIngestorEnabled: configuredDataStore === dataStoreType,
      },
      dataStoreType,
    };
  },
  shouldTestConnection() {
    return false;
  },
});

export default OtelCollectorService();
