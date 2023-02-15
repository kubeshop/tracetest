import {SupportedDataStores, TDraftDataStore} from 'types/Config.types';
import DataStore, {TRawDataStore} from 'models/DataStore.model';
import DataStoreConfig from 'models/DataStoreConfig.model';
import ElasticSearchService from './DataStores/ElasticSearch.service';
import OtelCollectorService from './DataStores/OtelCollector.service';
import SignalFxService from './DataStores/SignalFx.service';
import BaseClientService from './DataStores/BaseClient.service';
import JaegerService from './DataStores/Jaeger.service';
import AwsXRayService from './DataStores/AwsXRay.service';

interface IDataStoreService {
  getRequest(draft: TDraftDataStore, defaultDataStore: DataStore): Promise<TRawDataStore>;
  getInitialValues(config: DataStoreConfig): TDraftDataStore;
  validateDraft(config: TDraftDataStore): Promise<boolean>;
}

const dataStoreServiceMap = {
  [SupportedDataStores.JAEGER]: JaegerService,
  [SupportedDataStores.TEMPO]: BaseClientService,
  [SupportedDataStores.OpenSearch]: ElasticSearchService,
  [SupportedDataStores.ElasticApm]: ElasticSearchService,
  [SupportedDataStores.SignalFX]: SignalFxService,
  [SupportedDataStores.OtelCollector]: OtelCollectorService,
  [SupportedDataStores.NewRelic]: OtelCollectorService,
  [SupportedDataStores.Lightstep]: OtelCollectorService,
  [SupportedDataStores.Datadog]: OtelCollectorService,
  [SupportedDataStores.AWSXRay]: AwsXRayService,
} as const;

const DataStoreService = (): IDataStoreService => ({
  async getRequest(draft, defaultDataStore) {
    const dataStoreType = draft.dataStoreType || SupportedDataStores.JAEGER;
    const dataStoreValues = await dataStoreServiceMap[dataStoreType].getRequest(draft, dataStoreType);
    const isUpdate = !!defaultDataStore.id;

    const dataStore: TRawDataStore = isUpdate
      ? {id: defaultDataStore.id, ...dataStoreValues, isDefault: true}
      : {
          ...dataStoreValues,
          name: dataStoreType,
          type: dataStoreType as SupportedDataStores,
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
