import {IDataStorePluginMap, SupportedDataStores} from 'types/DataStore.types';
import AwsXRay from './forms/AwsXRay';
import AzureAppInsights from './forms/AzureAppInsights/AzureAppInsights';
import BaseClient from './forms/BaseClient';
import ElasticSearch from './forms/ElasticSearch';
import GrpcClient from './forms/GrpcClient';
import OpenTelemetryCollector from './forms/OpenTelemetryCollector';
import SignalFx from './forms/SignalFx/SignalFx';
import SumoLogic from './forms/SumoLogic';

export const DataStoreComponentMap: IDataStorePluginMap = {
  [SupportedDataStores.AWSXRay]: AwsXRay,
  [SupportedDataStores.AzureAppInsights]: AzureAppInsights,
  [SupportedDataStores.Datadog]: OpenTelemetryCollector,
  [SupportedDataStores.Dynatrace]: OpenTelemetryCollector,
  [SupportedDataStores.ElasticApm]: ElasticSearch,
  [SupportedDataStores.Honeycomb]: OpenTelemetryCollector,
  [SupportedDataStores.Instana]: OpenTelemetryCollector,
  [SupportedDataStores.JAEGER]: GrpcClient,
  [SupportedDataStores.Lightstep]: OpenTelemetryCollector,
  [SupportedDataStores.NewRelic]: OpenTelemetryCollector,
  [SupportedDataStores.OpenSearch]: ElasticSearch,
  [SupportedDataStores.OtelCollector]: OpenTelemetryCollector,
  [SupportedDataStores.SignalFX]: SignalFx,
  [SupportedDataStores.Signoz]: OpenTelemetryCollector,
  [SupportedDataStores.SumoLogic]: SumoLogic,
  [SupportedDataStores.TEMPO]: BaseClient,
};

export default DataStoreComponentMap;
