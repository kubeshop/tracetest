import {FormInstance} from 'antd';
import {Model, TDataStoreSchemas, TConfigSchemas} from 'types/Common.types';
import {THeader} from './Test.types';

export enum ConfigMode {
  NO_TRACING_MODE = 'NO_TRACING',
  READY = 'READY',
}

export enum SupportedDataStores {
  JAEGER = 'jaeger',
  TEMPO = 'tempo',
  OpenSearch = 'openSearch',
  SignalFX = 'signalFx',
  OtelCollector = 'otlp',
}

export type TRawDataStore = TDataStoreSchemas['DataStore'];
export type TDataStore = Model<
  TRawDataStore,
  {
    otlp?: {};
  }
>;

export type TDataStoreConfig = {
  defaultDataStore: TDataStore;
  mode: ConfigMode;
};

export type TRawGRPCClientSettings = TDataStoreSchemas['GRPCClientSettings'];

export type TTestConnectionRequest = TRawDataStore;
export type TTestConnectionResponse = TConfigSchemas['ConnectionResult'];

export interface IGRPCClientSettings extends TRawGRPCClientSettings {
  fileCA: File;
  fileCert: File;
  fileKey: File;
  rawHeaders?: THeader[];
}

interface IDataStore extends TRawDataStore {
  jaeger?: IGRPCClientSettings;
  tempo?: IGRPCClientSettings;
  otlp?: {};
}

export type TDraftDataStore = {
  dataStore?: IDataStore;
  dataStoreType?: SupportedDataStores;
};

export type TDataStoreForm = FormInstance<TDraftDataStore>;

export type TDataStoreService = {
  getRequest(values: TDraftDataStore, dataStoreType?: SupportedDataStores): Promise<TRawDataStore>;
  validateDraft(draft: TDraftDataStore): Promise<boolean>;
  getInitialValues(draft: TDataStoreConfig, dataStoreType?: SupportedDataStores): TDraftDataStore;
};

export interface IDataStorePluginProps {}
export interface IDataStorePluginMap
  extends Record<SupportedDataStores, (props: IDataStorePluginProps) => React.ReactElement> {}
