import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';
import {FormInstance} from 'antd';
import {Model, TConfigSchemas} from 'types/Common.types';

export enum ConfigMode {
  NO_TRACING_MODE = 'NO_TRACING',
  READY = 'READY',
}

export enum SupportedDataStores {
  JAEGER = 'jaeger',
  TEMPO = 'tempo',
  OpenSearch = 'opensearch',
  SignalFX = 'signalfx',
}

export type TRawConfig = TConfigSchemas['Config'];
export type TConfig = Model<
  TRawConfig,
  {
    mode: ConfigMode;
  }
>;

export type TRawTelemetryConfig = TConfigSchemas['TelemetryConfig'];
export type TTelemetryConfig = Model<TRawTelemetryConfig, {}>;

export type TRawDataStore = TConfigSchemas['DataStore'];

export type TRawServerConfig = TConfigSchemas['Server'];
export type TServerConfig = Model<TRawServerConfig, {}>;

export type TTestConnectionRequest = TConfigSchemas['TestConnectionRequest'];
export type TTestConnectionResponse = TConfigSchemas['TestConnectionResponse'];

export type TDraftConfig = {
  dataStore?: TRawDataStore;
  dataStoreType?: SupportedDataStores;
};

export interface ISetupConfigState {
  draftConfig: TDraftConfig;
  isFormValid: boolean;
}

export type TDraftConfigForm = FormInstance<TDraftConfig>;

export type TSetupConfigSliceActions = {
  reset: CaseReducer<ISetupConfigState>;
  setDraftConfig: CaseReducer<ISetupConfigState, PayloadAction<{draftConfig: TDraftConfig}>>;
  setIsFormValid: CaseReducer<ISetupConfigState, PayloadAction<{isValid: boolean}>>;
};

export type TDataStoreService = {
  getRequest(values: TDraftConfig, dataStoreType?: SupportedDataStores): Promise<TRawDataStore>;
  validateDraft(draft: TDraftConfig): Promise<boolean>;
  getInitialValues(draft: TConfig): TDraftConfig;
};

export interface IDataStorePluginProps {}
export interface IDataStorePluginMap
  extends Record<SupportedDataStores, (props: IDataStorePluginProps) => React.ReactElement> {}
