import {SupportedDataStores, TDataStoreService} from 'types/Config.types';

const OpenSearchService = (): TDataStoreService => ({
  getRequest({dataStore = {}}) {
    // todo add datastore mapping
    return Promise.resolve({
      type: SupportedDataStores.OpenSearch,
      ...dataStore,
    });
  },
  validateDraft() {
    return Promise.resolve(false);
  },
  getInitialValues() {
    return {};
  },
});

export default OpenSearchService();
