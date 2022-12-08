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

export type TRawConfig = TConfigSchemas['Config'];
export type TConfig = Model<
  TRawConfig,
  {
    mode: ConfigMode;
  }
>;

export type TSupportedDataStores = TConfigSchemas['SupportedDataStores'];
export type TRawTelemetryConfig = TConfigSchemas['TelemetryConfig'];
export type TTelemetryConfig = Model<TRawTelemetryConfig, {}>;
export type TUpdateDataStoreConfigRequest = TConfigSchemas['UpdateDataStoreConfigRequest'];

export type TRawDataStore = TConfigSchemas['DataStore'];

export type TRawServerConfig = TConfigSchemas['Server'];
export type TServerConfig = Model<TRawServerConfig, {}>;

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
  getInitialValues(draft: TConfig): TDraftDataStore;
};

export interface IDataStorePluginProps {}
export interface IDataStorePluginMap
  extends Record<SupportedDataStores, (props: IDataStorePluginProps) => React.ReactElement> {}
