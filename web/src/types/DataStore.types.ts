import {FormInstance} from 'antd';
import {Model, TDataStoreSchemas, TConfigSchemas} from 'types/Common.types';
import ConnectionTestStep from 'models/ConnectionResultStep.model';
import DataStore, {TRawOtlpDataStore} from 'models/DataStore.model';
import {THeader} from './Test.types';

export enum ConfigMode {
  NO_TRACING_MODE = 'NO_TRACING',
  READY = 'READY',
}

export enum ConnectionTypes {
  Collector = 'collector',
  Direct = 'direct',
}

export enum SupportedDataStores {
  AWSXRay = 'awsxray',
  AzureAppInsights = 'azureappinsights',
  Datadog = 'datadog',
  Dynatrace = 'dynatrace',
  ElasticApm = 'elasticapm',
  Honeycomb = 'honeycomb',
  Instana = 'instana',
  JAEGER = 'jaeger',
  Lightstep = 'lightstep',
  NewRelic = 'newrelic',
  OpenSearch = 'opensearch',
  OtelCollector = 'otlp',
  SignalFX = 'signalfx',
  Signoz = 'signoz',
  SumoLogic = 'sumologic',
  TEMPO = 'tempo',
}

export enum SupportedClientTypes {
  GRPC = 'grpc',
  HTTP = 'http',
}

export type TCollectorDataStores =
  | SupportedDataStores.Datadog
  | SupportedDataStores.Dynatrace
  | SupportedDataStores.Instana
  | SupportedDataStores.Lightstep
  | SupportedDataStores.NewRelic
  | SupportedDataStores.OtelCollector
  | SupportedDataStores.Signoz;

export type TRawGRPCClientSettings = TDataStoreSchemas['GRPCClientSettings'];
export type TRawElasticSearch = TDataStoreSchemas['ElasticSearch'];
export type TRawBaseClientSettings = TDataStoreSchemas['BaseClient'];
export type TRawHttpClientSettings = TDataStoreSchemas['HTTPClientSettings'];

export type TTestConnectionRequest = DataStore;
export type TRawConnectionResult = TConfigSchemas['ConnectionResult'];
export type TConnectionResult = Model<
  TRawConnectionResult,
  {
    allPassed: boolean;
    portCheck: ConnectionTestStep;
    authentication: ConnectionTestStep;
    connectivity: ConnectionTestStep;
    fetchTraces: ConnectionTestStep;
  }
>;

export type TTestConnectionResponse = TConfigSchemas['ConnectionResult'];

export interface IGRPCClientSettings extends TRawGRPCClientSettings {
  fileCA?: File;
  fileCert?: File;
  fileKey?: File;
  rawHeaders?: THeader[];
}

export interface IBaseClientSettings extends TRawBaseClientSettings {
  type: SupportedClientTypes;
}

export interface IHttpClientSettings extends TRawHttpClientSettings {
  fileCA?: File;
  fileCert?: File;
  fileKey?: File;
  rawHeaders?: THeader[];
}

export interface IElasticSearch extends TRawElasticSearch {
  certificateFile?: File;
}

export type IDataStore = DataStore & {
  agent?: TRawOtlpDataStore;
  datadog?: TRawOtlpDataStore;
  dynatrace?: TRawOtlpDataStore;
  elasticapm?: IElasticSearch;
  honeycomb?: TRawOtlpDataStore;
  instana?: TRawOtlpDataStore;
  jaeger?: IBaseClientSettings;
  lightstep?: TRawOtlpDataStore;
  newrelic?: TRawOtlpDataStore;
  opensearch?: IElasticSearch;
  otlp?: TRawOtlpDataStore;
  signoz?: TRawOtlpDataStore;
  tempo?: IBaseClientSettings;
};

export type TDraftDataStore = {
  dataStore?: IDataStore;
  dataStoreType?: SupportedDataStores;
};

export type TDataStoreForm = FormInstance<TDraftDataStore>;

export type TDataStoreService = {
  getRequest(values: TDraftDataStore, dataStoreType?: SupportedDataStores): Promise<DataStore>;
  validateDraft(draft: TDraftDataStore): Promise<boolean>;
  getInitialValues(
    draft: DataStore,
    dataStoreType?: SupportedDataStores,
    configuredDataStore?: SupportedDataStores
  ): TDraftDataStore;
  getIsOtlpBased(draft: TDraftDataStore): boolean;
};

export interface IDataStorePluginProps {}
export interface IDataStorePluginMap
  extends Record<SupportedDataStores, (props: IDataStorePluginProps) => React.ReactElement> {}
