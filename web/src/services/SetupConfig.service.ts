import {SupportedDataStores, TConfig, TDraftConfig, TRawConfig} from 'types/Config.types';
import GrpcClientService from './DataStores/GrpcClient.service';
import OpenSearchService from './DataStores/OpenSearch.service';
import SignalFxService from './DataStores/SignalFx.service';

interface ISetupConfigService {
  getRequest(draft: TDraftConfig): Promise<TRawConfig>;
  getInitialValues(config: TConfig): TDraftConfig;
}

const dataStoreServiceMap = {
  [SupportedDataStores.JAEGER]: GrpcClientService,
  [SupportedDataStores.TEMPO]: GrpcClientService,
  [SupportedDataStores.OpenSearch]: OpenSearchService,
  [SupportedDataStores.SignalFX]: SignalFxService,
} as const;

const SetupConfigService = (): ISetupConfigService => ({
  async getRequest(draft) {
    const dataStoreType = draft.dataStoreType || SupportedDataStores.JAEGER;
    const dataStore = await dataStoreServiceMap[dataStoreType].getRequest(draft);

    const config: TRawConfig = {
      telemetry: {
        dataStores: [dataStore],
      },
      server: {
        telemetry: {
          dataStore: dataStoreType,
        },
      },
    };

    return config;
  },

  getInitialValues({telemetry: {dataStores = []}}) {
    const [dataStore] = dataStores;

    return {
      dataStore,
    };
  },
});

export default SetupConfigService();
