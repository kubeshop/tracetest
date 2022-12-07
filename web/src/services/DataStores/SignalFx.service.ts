import {SupportedDataStores, TDataStoreService} from 'types/Config.types';

const SignalFxService = (): TDataStoreService => ({
  getRequest({dataStore = {}}) {
    // todo add datastore mapping
    return Promise.resolve({
      type: SupportedDataStores.SignalFX,
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

export default SignalFxService();
