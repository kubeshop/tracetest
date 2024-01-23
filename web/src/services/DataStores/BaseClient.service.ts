import {
  IBaseClientSettings,
  IGRPCClientSettings,
  IHttpClientSettings,
  SupportedClientTypes,
  SupportedDataStores,
  TDataStoreService,
  TRawBaseClientSettings,
  TRawGRPCClientSettings,
  TRawHttpClientSettings,
} from 'types/DataStore.types';
import HttpClientService from './HttpClient.service';
import GrpcClientService from './GrpcClient.service';

const BaseClientService = (): TDataStoreService => ({
  async getRequest({dataStore = {}}, dataStoreType = SupportedDataStores.TEMPO) {
    const {
      type = SupportedClientTypes.GRPC,
      grpc = {},
      http = {},
    } = dataStore[dataStoreType || SupportedDataStores.TEMPO] as IBaseClientSettings;

    const grpcRequest = await GrpcClientService.getRequest(grpc as IGRPCClientSettings);
    const httpRequest = await HttpClientService.getRequest(http as IHttpClientSettings);

    return Promise.resolve({
      [dataStoreType]: {
        grpc: grpcRequest,
        http: httpRequest,
        type,
      },
      type: dataStoreType,
      name: dataStoreType,
    });
  },
  validateDraft({dataStore = {name: '', type: SupportedDataStores.JAEGER}, dataStoreType}) {
    const {
      type = SupportedClientTypes.GRPC,
      grpc = {},
      http = {},
    } = dataStore[dataStoreType || SupportedDataStores.TEMPO] as IBaseClientSettings;

    switch (type) {
      case SupportedClientTypes.GRPC: {
        return GrpcClientService.validateDraft(grpc as IGRPCClientSettings);
      }
      case SupportedClientTypes.HTTP:
        return HttpClientService.validateDraft(http as IHttpClientSettings);
      default:
        return GrpcClientService.validateDraft(grpc as IGRPCClientSettings);
    }
  },
  getInitialValues(
    defaultDataStore = {name: '', type: SupportedDataStores.JAEGER},
    dataStoreType = SupportedDataStores.TEMPO
  ) {
    const {
      type = 'grpc',
      grpc = {},
      http = {},
    } = defaultDataStore[dataStoreType || SupportedDataStores.TEMPO] as TRawBaseClientSettings;

    return {
      dataStore: {
        [dataStoreType]: {
          grpc: GrpcClientService.getInitialValues(grpc as TRawGRPCClientSettings),
          http: HttpClientService.getInitialValues(http as TRawHttpClientSettings),
          type: type as SupportedClientTypes,
        },
        name: dataStoreType,
        type: dataStoreType,
      },
      dataStoreType,
    };
  },
  getIsOtlpBased() {
    return false;
  },
});

export default BaseClientService();
