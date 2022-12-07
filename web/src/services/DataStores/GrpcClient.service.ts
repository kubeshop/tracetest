import {SupportedDataStores, TDataStoreService} from 'types/Config.types';

const GrpcClientService = (): TDataStoreService => ({
  getRequest({dataStore = {}}, dataStoreType = SupportedDataStores.JAEGER) {
    // todo add datastore mapping
    return Promise.resolve({
      type: dataStoreType,
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

export default GrpcClientService();
