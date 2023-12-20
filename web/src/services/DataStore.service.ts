import {SupportedDataStores, TDataStoreService, TDraftDataStore} from 'types/DataStore.types';
import DataStore, {TRawDataStore} from 'models/DataStore.model';
import DataStoreConfig from 'models/DataStoreConfig.model';
import ElasticSearchService from './DataStores/ElasticSearch.service';
import OtelCollectorService from './DataStores/OtelCollector.service';
import SignalFxService from './DataStores/SignalFx.service';
import BaseClientService from './DataStores/BaseClient.service';
import JaegerService from './DataStores/Jaeger.service';
import AwsXRayService from './DataStores/AwsXRay.service';
import AzureAppInsightsService from './DataStores/AzureAppInsights.service';
import SumoLogicService from './DataStores/SumoLogic.service';

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
  [SupportedDataStores.Honeycomb]: OtelCollectorService,
  [SupportedDataStores.AWSXRay]: AwsXRayService,
  [SupportedDataStores.AzureAppInsights]: AzureAppInsightsService,
  [SupportedDataStores.Signoz]: OtelCollectorService,
  [SupportedDataStores.Dynatrace]: OtelCollectorService,
  [SupportedDataStores.SumoLogic]: SumoLogicService,
} as const;

interface IDataStoreService {
  getRequest(draft: TDraftDataStore, defaultDataStore: DataStore): Promise<TRawDataStore>;
  getInitialValues(config: DataStoreConfig, configuredDataStore?: SupportedDataStores): TDraftDataStore;
  validateDraft(config: TDraftDataStore): Promise<boolean>;
  shouldTestConnection(draft: TDraftDataStore): boolean;
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

  getInitialValues(dataStoreConfig, configuredDataStore) {
    const {defaultDataStore} = dataStoreConfig;
    const type = (defaultDataStore.type || SupportedDataStores.JAEGER) as SupportedDataStores;
    const dataStore = this._getDataStore(type);

    return {...dataStore.getInitialValues(dataStoreConfig, type, configuredDataStore), dataStoreType: type};
  },

  validateDraft(draft) {
    const dataStore = this._getDataStore(draft.dataStoreType);
    return dataStore.validateDraft(draft);
  },

  shouldTestConnection(draft) {
    const dataStore = this._getDataStore(draft.dataStoreType);

    return dataStore.shouldTestConnection(draft);
  },
});

export default DataStoreService();
