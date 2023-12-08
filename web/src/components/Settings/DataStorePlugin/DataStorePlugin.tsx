import {IDataStorePluginMap, SupportedDataStores} from 'types/DataStore.types';
import Agent from './forms/Agent';
import GrpcClient from './forms/GrpcClient';
import ElasticSearch from './forms/ElasticSearch';
import OpenTelemetryCollector from './forms/OpenTelemetryCollector';
import SignalFx from './forms/SignalFx/SignalFx';
import BaseClient from './forms/BaseClient';
import AwsXRay from './forms/AwsXRay';
import AzureAppInsights from './forms/AzureAppInsights/AzureAppInsights';
import SumoLogic from './forms/SumoLogic';

export const DataStoreComponentMap: IDataStorePluginMap = {
  [SupportedDataStores.Agent]: Agent,
  [SupportedDataStores.JAEGER]: GrpcClient,
  [SupportedDataStores.TEMPO]: BaseClient,
  [SupportedDataStores.SignalFX]: SignalFx,
  [SupportedDataStores.OpenSearch]: ElasticSearch,
  [SupportedDataStores.ElasticApm]: ElasticSearch,
  [SupportedDataStores.OtelCollector]: OpenTelemetryCollector,
  [SupportedDataStores.Lightstep]: OpenTelemetryCollector,
  [SupportedDataStores.Datadog]: OpenTelemetryCollector,
  [SupportedDataStores.NewRelic]: OpenTelemetryCollector,
  [SupportedDataStores.Honeycomb]: OpenTelemetryCollector,
  [SupportedDataStores.AWSXRay]: AwsXRay,
  [SupportedDataStores.AzureAppInsights]: AzureAppInsights,
  [SupportedDataStores.Signoz]: OpenTelemetryCollector,
  [SupportedDataStores.Dynatrace]: OpenTelemetryCollector,
  [SupportedDataStores.SumoLogic]: SumoLogic,
};

export default DataStoreComponentMap;
