import {FormInstance} from 'antd';
import {Model, TConfigSchemas} from 'types/Common.types';

export enum ConfigMode {
  NO_TRACING_MODE = 'NO_TRACING',
  READY = 'READY',
}

export enum SupportedDataStores {
  JAEGER = 'jaeger',
  TEMPO = 'tempo',
  OpenSearch = 'openSearch',
  SignalFX = 'signalFx',
}

export type TSupportedDataStores = TConfigSchemas['SupportedDataStores'];
export type TRawDataStoreConfig = TConfigSchemas['DataStoreConfig'];
export type TDataStoreConfig = Model<
  TRawDataStoreConfig,
  {
    mode: ConfigMode;
  }
>;

export type TRawDataStore = TConfigSchemas['DataStore'];
export type TDataStore = Model<TRawDataStore, {}>;
export type TRawGRPCClientSettings = TConfigSchemas['GRPCClientSettings'];

export type TTestConnectionRequest = TConfigSchemas['DataStoreTestConnectionRequest'];
export type TTestConnectionResponse = TConfigSchemas['DataStoreTestConnectionResponse'];

export interface IGRPCClientSettings extends TRawGRPCClientSettings {
  fileCA: File;
  fileCert: File;
  fileKey: File;
}

interface IDataStore extends TRawDataStore {
  jaeger?: IGRPCClientSettings;
  tempo?: IGRPCClientSettings;
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
