import {SupportedDataStores, TConfig, TDraftDataStore, TUpdateDataStoreConfigRequest} from 'types/Config.types';
import GrpcClientService from './DataStores/GrpcClient.service';
import OpenSearchService from './DataStores/OpenSearch.service';
import SignalFxService from './DataStores/SignalFx.service';

interface IDataStoreService {
  getRequest(draft: TDraftDataStore): Promise<TUpdateDataStoreConfigRequest>;
  getInitialValues(config: TConfig): TDraftDataStore;
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
    const dataStore = await dataStoreServiceMap[dataStoreType].getRequest(draft);

    const config: TUpdateDataStoreConfigRequest = {
      dataStores: [dataStore],
      defaultDataStore: dataStoreType,
    };

    return config;
  },

  getInitialValues(config) {
    const {
      server: {telemetry: {dataStore: dataStoreType} = {}},
    } = config;
    const type = (dataStoreType || SupportedDataStores.JAEGER) as SupportedDataStores;

    return dataStoreServiceMap[type].getInitialValues(config);
  },

  validateDraft(draft) {
    const dataStoreType = draft.dataStoreType || SupportedDataStores.JAEGER;
    const dataStore = dataStoreServiceMap[dataStoreType];

    return dataStore.validateDraft(draft);
  },
});

export default DataStoreService();
