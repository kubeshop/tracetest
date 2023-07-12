import {IDataStore, SupportedDataStores, TDataStoreService} from 'types/DataStore.types';
import {TRawOtlpDataStore} from 'models/DataStore.model';

const OtelCollectorService = (): TDataStoreService => ({
  getRequest(draft, dataStoreType = SupportedDataStores.OtelCollector) {
    return Promise.resolve({
      type: dataStoreType,
      name: dataStoreType,
    });
  },
  validateDraft({dataStore = {} as IDataStore, dataStoreType = SupportedDataStores.OtelCollector}) {
    const {isIngestorEnabled = false} =
      (dataStore[dataStoreType || SupportedDataStores.OtelCollector] as TRawOtlpDataStore) ?? {};

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
        [dataStoreType]: {
          isIngestorEnabled: configuredDataStore === dataStoreType,
        },
      },
      dataStoreType,
    };
  },
  shouldTestConnection() {
    return false;
  },
});

export default OtelCollectorService();
