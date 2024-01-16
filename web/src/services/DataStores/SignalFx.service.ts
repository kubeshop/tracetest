import {SupportedDataStores, TDataStoreService} from 'types/DataStore.types';

const SignalFxService = (): TDataStoreService => ({
  getRequest({dataStore: {signalfx: {realm = '', token = ''} = {}} = {}}) {
    return Promise.resolve({
      type: SupportedDataStores.SignalFX,
      name: SupportedDataStores.SignalFX,
      signalfx: {
        realm,
        token,
      },
    });
  },
  validateDraft({dataStore: {signalfx: {realm = '', token = ''} = {}} = {}}) {
    if (!realm || !token) return Promise.resolve(false);

    return Promise.resolve(true);
  },
  getInitialValues({signalfx = {}}) {
    const {realm = '', token = ''} = signalfx;

    return {
      dataStore: {name: SupportedDataStores.SignalFX, type: SupportedDataStores.SignalFX, signalfx: {realm, token}},
      dataStoreType: SupportedDataStores.SignalFX,
    };
  },
  shouldTestConnection() {
    return true;
  },
});

export default SignalFxService();
