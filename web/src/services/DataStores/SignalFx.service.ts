import {SupportedDataStores, TDataStoreService} from 'types/DataStore.types';

const SignalFxService = (): TDataStoreService => ({
  getRequest({dataStore: {signalFx: {realm = '', token = ''} = {}} = {}}) {
    return Promise.resolve({
      type: SupportedDataStores.SignalFX,
      name: SupportedDataStores.SignalFX,
      signalFx: {
        realm,
        token,
      },
    });
  },
  validateDraft({dataStore: {signalFx: {realm = '', token = ''} = {}} = {}}) {
    if (!realm || !token) return Promise.resolve(false);

    return Promise.resolve(true);
  },
  getInitialValues({defaultDataStore: {signalFx = {}} = {}}) {
    const {realm = '', token = ''} = signalFx;

    return {
      dataStore: {name: SupportedDataStores.SignalFX, type: SupportedDataStores.SignalFX, signalFx: {realm, token}},
      dataStoreType: SupportedDataStores.SignalFX,
    };
  },
});

export default SignalFxService();
