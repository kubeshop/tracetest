import {SupportedDataStores, TDataStoreService, TDraftDataStore} from 'types/DataStore.types';
import DataStore, {TRawDataStore} from 'models/DataStore.model';
import ElasticSearchService from './DataStores/ElasticSearch.service';
import OtelCollectorService from './DataStores/OtelCollector.service';
import SignalFxService from './DataStores/SignalFx.service';
import BaseClientService from './DataStores/BaseClient.service';
import JaegerService from './DataStores/Jaeger.service';
import AwsXRayService from './DataStores/AwsXRay.service';
import AzureAppInsightsService from './DataStores/AzureAppInsights.service';
import SumoLogicService from './DataStores/SumoLogic.service';

const dataStoreServiceMap = {
  [SupportedDataStores.AWSXRay]: AwsXRayService,
  [SupportedDataStores.AzureAppInsights]: AzureAppInsightsService,
  [SupportedDataStores.Datadog]: OtelCollectorService,
  [SupportedDataStores.Dynatrace]: OtelCollectorService,
  [SupportedDataStores.ElasticApm]: ElasticSearchService,
  [SupportedDataStores.Honeycomb]: OtelCollectorService,
  [SupportedDataStores.Instana]: OtelCollectorService,
  [SupportedDataStores.JAEGER]: JaegerService,
  [SupportedDataStores.Lightstep]: OtelCollectorService,
  [SupportedDataStores.NewRelic]: OtelCollectorService,
  [SupportedDataStores.OpenSearch]: ElasticSearchService,
  [SupportedDataStores.OtelCollector]: OtelCollectorService,
  [SupportedDataStores.SignalFX]: SignalFxService,
  [SupportedDataStores.Signoz]: OtelCollectorService,
  [SupportedDataStores.SumoLogic]: SumoLogicService,
  [SupportedDataStores.TEMPO]: BaseClientService,
} as const;

interface IDataStoreService {
  getRequest(draft: TDraftDataStore, defaultDataStore: DataStore): Promise<TRawDataStore>;
  getInitialValues(config: DataStore, configuredDataStore?: SupportedDataStores): TDraftDataStore;
  validateDraft(config: TDraftDataStore): Promise<boolean>;
  getIsOtlpBased(draft: TDraftDataStore): boolean;
  _getDataStore(type?: SupportedDataStores): TDataStoreService;
}

const DataStoreService = (): IDataStoreService => ({
  _getDataStore(type = SupportedDataStores.JAEGER) {
    return dataStoreServiceMap[type] || OtelCollectorService;
  },
  async getRequest(draft, defaultDataStore) {
    const dataStoreType = draft.dataStoreType || SupportedDataStores.JAEGER;
    const dataStoreService = this._getDataStore(dataStoreType);
    const dataStoreValues = await dataStoreService.getRequest(draft, dataStoreType);
    const isUpdate = !!defaultDataStore.id;

    const dataStore: DataStore = isUpdate
      ? {id: defaultDataStore.id, ...dataStoreValues, default: true}
      : {
          ...dataStoreValues,
          name: dataStoreType,
          type: dataStoreType as SupportedDataStores,
          default: true,
        };

    return {
      type: 'DataStore',
      spec: dataStore,
    } as TRawDataStore;
  },

  getInitialValues(defaultDataStore, configuredDataStore) {
    const type = (defaultDataStore.type || SupportedDataStores.JAEGER) as SupportedDataStores;
    const dataStore = this._getDataStore(type);

    return {...dataStore.getInitialValues(defaultDataStore, type, configuredDataStore), dataStoreType: type};
  },

  validateDraft(draft) {
    const dataStore = this._getDataStore(draft.dataStoreType);
    return dataStore.validateDraft(draft);
  },

  getIsOtlpBased(draft) {
    const dataStore = this._getDataStore(draft.dataStoreType);

    return dataStore.getIsOtlpBased(draft);
  },
});

export default DataStoreService();
