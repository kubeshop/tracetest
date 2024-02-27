import {SupportedDataStores, TDataStoreService} from 'types/DataStore.types';

const SumoLogicService = (): TDataStoreService => ({
  getRequest({dataStore: {sumologic: {url = '', accessID = '', accessKey = ''} = {}} = {}}) {
    return Promise.resolve({
      type: SupportedDataStores.SumoLogic,
      name: SupportedDataStores.SumoLogic,
      sumologic: {
        url,
        accessID,
        accessKey,
      },
    });
  },
  validateDraft({dataStore: {sumologic: {url = '', accessID = '', accessKey = ''} = {}} = {}}) {
    if (!url || !accessID || !accessKey) return Promise.resolve(false);

    return Promise.resolve(true);
  },
  getInitialValues({sumologic = {}}) {
    const {url = '', accessID = '', accessKey = ''} = sumologic;

    return {
      dataStore: {
        name: SupportedDataStores.SumoLogic,
        type: SupportedDataStores.SumoLogic,
        sumologic: {url, accessID, accessKey},
      },
      dataStoreType: SupportedDataStores.SumoLogic,
    };
  },
  getIsOtlpBased() {
    return false;
  },
  getPublicInfo({sumologic = {}}) {
    const {url = ''} = sumologic;

    return {
      URL: url,
    };
  },
});

export default SumoLogicService();
