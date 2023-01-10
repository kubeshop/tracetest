import {IElasticSearch, SupportedDataStores, TDataStore, TDataStoreService, TRawElasticSearch} from 'types/Config.types';
import Validator from 'utils/Validator';

const ElasticSearchService = (): TDataStoreService => ({
  getRequest(
    {dataStore: {openSearch: {index = '', username = '', password = '', addresses = [], certificate = ''} = {}} = {}},
    dataStoreType = SupportedDataStores.OpenSearch
  ) {
    return Promise.resolve({
      name: dataStoreType,
      type: dataStoreType,
      [dataStoreType]: {
        index,
        username,
        password,
        addresses,
        certificate,
      },
    });
  },
  validateDraft({dataStore: {openSearch: {index = '', addresses = []} = {}} = {}}) {
    const [address] = addresses;
    if (!index || !Validator.url(address)) return Promise.resolve(false);

    return Promise.resolve(true);
  },
  getInitialValues(
    {defaultDataStore = {name: '', type: SupportedDataStores.OpenSearch} as TDataStore},
    dataStoreType = SupportedDataStores.OpenSearch
  ) {
    const {
      index = '',
      username = '',
      password = '',
      addresses = [''],
      certificate = '',
    } = (defaultDataStore[dataStoreType] as TRawElasticSearch) ?? {};

    const draftDataStore: IElasticSearch = {
      index,
      username,
      password,
      addresses,
      certificateFile: certificate ? new File([certificate], 'certificate') : undefined,
    };

    return {
      dataStore: {
        [dataStoreType]: draftDataStore,
        name: dataStoreType,
        type: dataStoreType,
      },
      dataStoreType,
    };
  },
});

export default ElasticSearchService();
