import {IDataStorePluginMap, SupportedDataStores} from 'types/Config.types';
import GrpcClient from './forms/GrpcClient';
import OpenSearch from './forms/OpenSearch';
import SignalFx from './forms/SignalFx/SignalFx';

export const DataStoreComponentMap: IDataStorePluginMap = {
  [SupportedDataStores.JAEGER]: GrpcClient,
  [SupportedDataStores.TEMPO]: GrpcClient,
  [SupportedDataStores.SignalFX]: SignalFx,
  [SupportedDataStores.OpenSearch]: OpenSearch,
};

export default DataStoreComponentMap;
