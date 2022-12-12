import {SupportedDataStores, TDataStoreConfig, TDraftDataStore, TRawDataStoreConfig} from 'types/Config.types';
import GrpcClientService from './DataStores/GrpcClient.service';
import OpenSearchService from './DataStores/OpenSearch.service';
import SignalFxService from './DataStores/SignalFx.service';

interface IDataStoreService {
  getRequest(draft: TDraftDataStore): Promise<TRawDataStoreConfig>;
  getInitialValues(config: TDataStoreConfig): TDraftDataStore;
  validateDraft(config: TDraftDataStore): Promise<boolean>;
}

const dataStoreServiceMap = {
  [SupportedDataStores.JAEGER]: GrpcClientService,
  [SupportedDataStores.TEMPO]: GrpcClientService,
  [SupportedDataStores.OpenSearch]: OpenSearchService,
  [SupportedDataStores.SignalFX]: SignalFxService,
} as const;

const DataStoreService = (): IDataStoreService => ({
  async getRequest(draft) {
    const dataStoreType = draft.dataStoreType || SupportedDataStores.JAEGER;
    const dataStore = await dataStoreServiceMap[dataStoreType].getRequest(draft, dataStoreType);

    const config: TRawDataStoreConfig = {
      dataStores: [{...dataStore, name: dataStoreType, type: dataStoreType}],
      defaultDataStore: dataStoreType,
    };

    return config;
  },

  getInitialValues(config) {
    const {defaultDataStore = '', dataStores = []} = config;
    const dataStoreType = dataStores.find(({name}) => name === defaultDataStore)?.type;
    const type = (dataStoreType || SupportedDataStores.JAEGER) as SupportedDataStores;

    return {...dataStoreServiceMap[type].getInitialValues(config, type), dataStoreType: type};
  },

  validateDraft(draft) {
    const dataStoreType = draft.dataStoreType || SupportedDataStores.JAEGER;
    const dataStore = dataStoreServiceMap[dataStoreType];

    return dataStore.validateDraft(draft);
  },
});

export default DataStoreService();
