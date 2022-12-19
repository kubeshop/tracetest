import {SupportedDataStores, TDataStore, TDataStoreConfig, TDraftDataStore, TRawDataStore} from 'types/Config.types';
import GrpcClientService from './DataStores/GrpcClient.service';
import OpenSearchService from './DataStores/OpenSearch.service';
import OtelCollectorService from './DataStores/OtelCollector.service';
import SignalFxService from './DataStores/SignalFx.service';

interface IDataStoreService {
  getRequest(draft: TDraftDataStore, defaultDataStore: TDataStore): Promise<TRawDataStore>;
  getInitialValues(config: TDataStoreConfig): TDraftDataStore;
  validateDraft(config: TDraftDataStore): Promise<boolean>;
}

const dataStoreServiceMap = {
  [SupportedDataStores.JAEGER]: GrpcClientService,
  [SupportedDataStores.TEMPO]: GrpcClientService,
  [SupportedDataStores.OpenSearch]: OpenSearchService,
  [SupportedDataStores.SignalFX]: SignalFxService,
  [SupportedDataStores.OtelCollector]: OtelCollectorService,
} as const;

const DataStoreService = (): IDataStoreService => ({
  async getRequest(draft, defaultDataStore) {
    const dataStoreType = draft.dataStoreType || SupportedDataStores.JAEGER;
    const dataStoreValues = await dataStoreServiceMap[dataStoreType].getRequest(draft, dataStoreType);
    const isUpdate = !!defaultDataStore.id && defaultDataStore.type === dataStoreType;

    const dataStore: TRawDataStore = isUpdate
      ? {...defaultDataStore, ...dataStoreValues, isDefault: true}
      : {
          ...dataStoreValues,
          name: dataStoreType,
          type: dataStoreType,
          isDefault: true,
        };

    return dataStore;
  },

  getInitialValues(dataStoreConfig) {
    const {defaultDataStore} = dataStoreConfig;
    const type = (defaultDataStore.type || SupportedDataStores.JAEGER) as SupportedDataStores;

    return {...dataStoreServiceMap[type].getInitialValues(dataStoreConfig, type), dataStoreType: type};
  },

  validateDraft(draft) {
    const dataStoreType = draft.dataStoreType || SupportedDataStores.JAEGER;
    const dataStore = dataStoreServiceMap[dataStoreType];

    return dataStore.validateDraft(draft);
  },
});

export default DataStoreService();
