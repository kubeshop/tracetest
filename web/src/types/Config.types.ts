import {Model, TConfigSchemas} from 'types/Common.types';

export enum ConfigMode {
  NO_TRACING_MODE = 'NO_TRACING',
  READY = 'READY',
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

export type TRawServerConfig = TConfigSchemas['Server'];
export type TServerConfig = Model<TRawServerConfig, {}>;

export type TTestConnectionRequest = TConfigSchemas['TestConnectionRequest'];
export type TTestConnectionResponse = TConfigSchemas['TestConnectionResponse'];
