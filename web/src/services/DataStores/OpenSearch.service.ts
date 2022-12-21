import {SupportedDataStores, TDataStoreService} from 'types/Config.types';
import Validator from 'utils/Validator';

const OpenSearchService = (): TDataStoreService => ({
  getRequest({dataStore: {openSearch: {index = '', username = '', password = '', addresses = []} = {}} = {}}) {
    return Promise.resolve({
      type: SupportedDataStores.OpenSearch,
      openSearch: {
        index,
        username,
        password,
        addresses,
      },
    });
  },
  validateDraft({dataStore: {openSearch: {index = '', username = '', password = '', addresses = []} = {}} = {}}) {
    const [address] = addresses;
    if (!index || !username || !password || !Validator.url(address)) return Promise.resolve(false);

    return Promise.resolve(true);
  },
  getInitialValues({defaultDataStore: {openSearch = {}} = {}}) {
    const {index = '', username = '', password = '', addresses = ['']} = openSearch;

    return {
      dataStore: {
        openSearch: {
          index,
          username,
          password,
          addresses,
        },
      },
      dataStoreType: SupportedDataStores.OpenSearch,
    };
  },
});

export default OpenSearchService();
