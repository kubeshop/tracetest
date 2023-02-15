import {FormInstance} from 'antd';
import {Model, TDataStoreSchemas, TConfigSchemas} from 'types/Common.types';
import ConnectionTestStep from 'models/ConnectionResultStep.model';
import {TRawDataStore} from 'models/DataStore.model';
import DataStoreConfig from 'models/DataStoreConfig.model';
import {THeader} from './Test.types';

export enum ConfigMode {
  NO_TRACING_MODE = 'NO_TRACING',
  READY = 'READY',
}

export enum SupportedDataStores {
  JAEGER = 'jaeger',
  TEMPO = 'tempo',
  OtelCollector = 'otlp',
  NewRelic = 'newRelic',
  Lightstep = 'lightstep',
  OpenSearch = 'openSearch',
  ElasticApm = 'elasticApm',
  SignalFX = 'signalFx',
  Datadog = 'datadog',
  AWSXRay = 'awsxray',
}

export enum SupportedClientTypes {
  GRPC = 'grpc',
  HTTP = 'http',
}

export type TCollectorDataStores =
  | SupportedDataStores.NewRelic
  | SupportedDataStores.OtelCollector
  | SupportedDataStores.Lightstep
  | SupportedDataStores.Datadog;

export type TRawGRPCClientSettings = TDataStoreSchemas['GRPCClientSettings'];
export type TRawElasticSearch = TDataStoreSchemas['ElasticSearch'];
export type TRawBaseClientSettings = TDataStoreSchemas['BaseClient'];
export type TRawHttpClientSettings = TDataStoreSchemas['HTTPClientSettings'];

export type TTestConnectionRequest = TRawDataStore;
export type TRawConnectionResult = TConfigSchemas['ConnectionResult'];
export type TConnectionResult = Model<
  TRawConnectionResult,
  {
    allPassed: boolean;
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

type IDataStore = TRawDataStore & {
  jaeger?: IBaseClientSettings;
  tempo?: IBaseClientSettings;
  openSearch?: IElasticSearch;
  elasticApm?: IElasticSearch;
  otlp?: {};
  lightstep?: {};
  newRelic?: {};
  datadog?: {};
}

export type TDraftDataStore = {
  dataStore?: IDataStore;
  dataStoreType?: SupportedDataStores;
};

export type TDataStoreForm = FormInstance<TDraftDataStore>;

export type TDataStoreService = {
  getRequest(values: TDraftDataStore, dataStoreType?: SupportedDataStores): Promise<TRawDataStore>;
  validateDraft(draft: TDraftDataStore): Promise<boolean>;
  getInitialValues(draft: DataStoreConfig, dataStoreType?: SupportedDataStores): TDraftDataStore;
};

export interface IDataStorePluginProps {}
export interface IDataStorePluginMap
  extends Record<SupportedDataStores, (props: IDataStorePluginProps) => React.ReactElement> {}
