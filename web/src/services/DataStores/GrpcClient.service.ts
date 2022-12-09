import {SupportedDataStores, TDataStoreService} from 'types/Config.types';

const GrpcClientService = (): TDataStoreService => ({
  getRequest({dataStore = {}}, dataStoreType = SupportedDataStores.JAEGER) {
    return Promise.resolve({
      type: dataStoreType,
      ...dataStore,
    });
  },
  validateDraft({dataStore = {}}) {
    return Promise.resolve(false);
  },
  getInitialValues({telemetry: {dataStores = []} = {}}) {
    return {
      dataStore: {
        jaeger: {},
      },
      dataStoreType: SupportedDataStores.JAEGER,
    };
  },
});

export default GrpcClientService();
