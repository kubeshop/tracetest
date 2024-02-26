import {IElasticSearch, SupportedDataStores, TDataStoreService, TRawElasticSearch} from 'types/DataStore.types';
import Validator from 'utils/Validator';

const ElasticSearchService = (): TDataStoreService => ({
  async getRequest({dataStore = {}}, dataStoreType = SupportedDataStores.OpenSearch) {
    const values = dataStore[dataStoreType || SupportedDataStores.OpenSearch] as IElasticSearch;
    const {
      certificateFile = new File([''], 'certificate'),
      addresses = [],
      index = '',
      username = '',
      password = '',
      insecureSkipVerify = false,
    } = values;
    const certificate = await certificateFile.text();
    return Promise.resolve({
      name: dataStoreType,
      type: dataStoreType,
      [dataStoreType]: {
        index,
        username,
        password,
        addresses,
        certificate,
        insecureSkipVerify,
      },
    });
  },
  validateDraft({dataStore = {}, dataStoreType}) {
    const values = dataStore[dataStoreType || SupportedDataStores.OpenSearch] as IElasticSearch;
    const {addresses = [], index = ''} = values;
    const [address] = addresses;
    if (!index || !Validator.url(address)) return Promise.resolve(false);

    return Promise.resolve(true);
  },
  getInitialValues(
    defaultDataStore = {name: '', type: SupportedDataStores.OpenSearch},
    dataStoreType = SupportedDataStores.OpenSearch
  ) {
    const {
      index = '',
      username = '',
      password = '',
      addresses = [''],
      certificate = '',
      insecureSkipVerify = false,
    } = (defaultDataStore[dataStoreType] as TRawElasticSearch) ?? {};

    const draftDataStore: IElasticSearch = {
      index,
      username,
      password,
      addresses,
      certificateFile: certificate ? new File([certificate], 'certificate') : undefined,
      insecureSkipVerify,
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
  getIsOtlpBased() {
    return false;
  },
  getPublicInfo({opensearch = {}}) {
    const {addresses = [], index = ''} = opensearch;

    return {
      Address: addresses.join(', '),
      Index: index,
    };
  },
});

export default ElasticSearchService();
