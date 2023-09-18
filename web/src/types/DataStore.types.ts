import {FormInstance} from 'antd';
import {Model, TDataStoreSchemas, TConfigSchemas} from 'types/Common.types';
import ConnectionTestStep from 'models/ConnectionResultStep.model';
import DataStore, {TRawOtlpDataStore} from 'models/DataStore.model';
import DataStoreConfig from 'models/DataStoreConfig.model';
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
  Agent = 'agent',
  JAEGER = 'jaeger',
  TEMPO = 'tempo',
  OtelCollector = 'otlp',
  NewRelic = 'newrelic',
  Lightstep = 'lightstep',
  OpenSearch = 'opensearch',
  ElasticApm = 'elasticapm',
  SignalFX = 'signalfx',
  Datadog = 'datadog',
  AWSXRay = 'awsxray',
  Honeycomb = 'honeycomb',
  AzureAppInsights = 'azureappinsights',
  Signoz = 'signoz',
  Dynatrace = 'dynatrace',
}

export enum SupportedClientTypes {
  GRPC = 'grpc',
  HTTP = 'http',
}

export type TCollectorDataStores =
  | SupportedDataStores.NewRelic
  | SupportedDataStores.OtelCollector
  | SupportedDataStores.Lightstep
  | SupportedDataStores.Datadog
  | SupportedDataStores.Signoz
  | SupportedDataStores.Dynatrace;

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
  jaeger?: IBaseClientSettings;
  tempo?: IBaseClientSettings;
  opensearch?: IElasticSearch;
  elasticapm?: IElasticSearch;
  otlp?: TRawOtlpDataStore;
  lightstep?: TRawOtlpDataStore;
  newrelic?: TRawOtlpDataStore;
  datadog?: TRawOtlpDataStore;
  honeycomb?: TRawOtlpDataStore;
  signoz?: TRawOtlpDataStore;
  dynatrace?: TRawOtlpDataStore;
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
    draft: DataStoreConfig,
    dataStoreType?: SupportedDataStores,
    configuredDataStore?: SupportedDataStores
  ): TDraftDataStore;
  shouldTestConnection(draft: TDraftDataStore): boolean;
};

export interface IDataStorePluginProps {}
export interface IDataStorePluginMap
  extends Record<SupportedDataStores, (props: IDataStorePluginProps) => React.ReactElement> {}
