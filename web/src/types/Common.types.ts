import {external} from './Generated.types';

export interface IEnv {
  analyticsEnabled: boolean;
  appVersion: string;
  env: string;
  experimentalFeatures: string[];
  measurementId: string;
  serverID: string;
  serverPathPrefix: string;
  segmentLoaded: boolean;
  isTracetestDev: boolean;
}

export interface IMockFactory<T, R> {
  (): {
    raw(data?: Partial<R>): R;
    model(data?: Partial<R>): T;
  };
}

export type THttpSchemas = external['http.yaml']['components']['schemas'];
export type TTraceSchemas = external['trace.yaml']['components']['schemas'];
export type TTestSchemas = external['tests.yaml']['components']['schemas'];
export type TTriggerSchemas = external['triggers.yaml']['components']['schemas'];
export type TGrpcSchemas = external['grpc.yaml']['components']['schemas'];
export type TTraceIDSchemas = external['traceid.yaml']['components']['schemas'];
export type TEnvironmentSchemas = external['environments.yaml']['components']['schemas'];
export type TExpressionsSchemas = external['expressions.yaml']['components']['schemas'];
export type TTransactionsSchemas = external['transactions.yaml']['components']['schemas'];
export type TResourceSchemas = external['resources.yaml']['components']['schemas'];
export type TDataStoreSchemas = external['dataStores.yaml']['components']['schemas'];
export type TConfigSchemas = external['config.yaml']['components']['schemas'];
export type TVariablesSchemas = external['variables.yaml']['components']['schemas'];
export type TTestEventsSchemas = external['testEvents.yaml']['components']['schemas'];
export type TLintersSchemas = external['linterns.yaml']['components']['schemas'];

export type TSelector = TTestSchemas['Selector'];

export type Modify<T, R> = Omit<T, keyof R> & R;

export type Model<T, R> = Modify<Required<T>, R>;
