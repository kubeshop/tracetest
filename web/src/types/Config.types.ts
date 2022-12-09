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
export type TDataStoreConfig = Model<TRawDataStoreConfig, {
  mode: ConfigMode;
}>;

export type TRawDataStore = TConfigSchemas['DataStore'];
export type TDataStore = Model<TRawDataStore, {}>;

export type TTestConnectionRequest = TConfigSchemas['TestConnectionRequest'];
export type TTestConnectionResponse = TConfigSchemas['TestConnectionResponse'];

export type TDraftDataStore = {
  dataStore?: TRawDataStore;
  dataStoreType?: SupportedDataStores;
};

export type TDataStoreForm = FormInstance<TDraftDataStore>;

export type TDataStoreService = {
  getRequest(values: TDraftDataStore, dataStoreType?: SupportedDataStores): Promise<TRawDataStore>;
  validateDraft(draft: TDraftDataStore): Promise<boolean>;
  getInitialValues(draft: TDataStoreConfig): TDraftDataStore;
};

export interface IDataStorePluginProps {}
export interface IDataStorePluginMap
  extends Record<SupportedDataStores, (props: IDataStorePluginProps) => React.ReactElement> {}
