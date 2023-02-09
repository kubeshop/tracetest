import {
  IBaseClientSettings,
  IGRPCClientSettings,
  SupportedClientTypes,
  SupportedDataStores,
  TDataStoreService,
  TRawGRPCClientSettings,
} from 'types/Config.types';
import DataStore from 'models/DataStore.model';
import GrpcClientService from './GrpcClient.service';

const JaegerService = (): TDataStoreService => ({
  async getRequest({dataStore = {}}) {
    const {grpc = {}} = dataStore[SupportedDataStores.JAEGER] as IBaseClientSettings;
    const grpcRequest = await GrpcClientService.getRequest(grpc as IGRPCClientSettings);

    return Promise.resolve({
      [SupportedDataStores.JAEGER]: grpcRequest,
      name: SupportedDataStores.JAEGER,
      type: SupportedDataStores.JAEGER,
    });
  },
  validateDraft({dataStore = {}}) {
    const {grpc = {}} = dataStore[SupportedDataStores.JAEGER] as IBaseClientSettings;
    return GrpcClientService.validateDraft(grpc as IGRPCClientSettings);
  },
  getInitialValues({defaultDataStore = {name: '', type: SupportedDataStores.JAEGER} as DataStore}) {
    const values = defaultDataStore[SupportedDataStores.JAEGER] as TRawGRPCClientSettings;

    return {
      dataStore: {
        [SupportedDataStores.JAEGER]: {
          grpc: GrpcClientService.getInitialValues(values),
          type: SupportedClientTypes.GRPC,
        },
        name: SupportedDataStores.JAEGER,
        type: SupportedDataStores.JAEGER,
      },
      dataStoreType: SupportedDataStores.JAEGER,
    };
  },
});

export default JaegerService();
