import {IDataStorePluginMap, SupportedDataStores} from 'types/DataStore.types';
import GrpcClient from './forms/GrpcClient';
import ElasticSearch from './forms/ElasticSearch';
import OpenTelemetryCollector from './forms/OpenTelemetryCollector';
import SignalFx from './forms/SignalFx/SignalFx';
import BaseClient from './forms/BaseClient';
import AwsXRay from './forms/AwsXRay';

export const DataStoreComponentMap: IDataStorePluginMap = {
  [SupportedDataStores.JAEGER]: GrpcClient,
  [SupportedDataStores.TEMPO]: BaseClient,
  [SupportedDataStores.SignalFX]: SignalFx,
  [SupportedDataStores.OpenSearch]: ElasticSearch,
  [SupportedDataStores.ElasticApm]: ElasticSearch,
  [SupportedDataStores.OtelCollector]: OpenTelemetryCollector,
  [SupportedDataStores.Lightstep]: OpenTelemetryCollector,
  [SupportedDataStores.Datadog]: OpenTelemetryCollector,
  [SupportedDataStores.NewRelic]: OpenTelemetryCollector,
  [SupportedDataStores.AWSXRay]: AwsXRay,
};

export default DataStoreComponentMap;
