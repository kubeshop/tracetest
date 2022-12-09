import {SupportedDataStores, TDataStoreService} from 'types/Config.types';

const SignalFxService = (): TDataStoreService => ({
  getRequest({dataStore: {signalFx: {realm = '', token = ''} = {}} = {}}) {
    return Promise.resolve({
      type: SupportedDataStores.SignalFX,
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
  getInitialValues({dataStores: [{signalFx = {}}] = []}) {
    const {realm = '', token = ''} = signalFx;

    return {
      dataStore: {signalFx: {realm, token}},
      dataStoreType: SupportedDataStores.SignalFX,
    };
  },
});

export default SignalFxService();
