import {external} from './Generated.types';

export interface IEnv {
  analyticsEnabled: boolean;
  appVersion: string;
  baseApiUrl: string;
  env: string;
  experimentalFeatures: string[];
  measurementId: string;
  serverID: string;
  serverPathPrefix: string;
  segmentLoaded: boolean;
  isTracetestDev: boolean;
  posthogKey: string;
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
export type TKafkaSchemas = external['kafka.yaml']['components']['schemas'];
export type TVariableSetSchemas = external['variableSets.yaml']['components']['schemas'];
export type TExpressionsSchemas = external['expressions.yaml']['components']['schemas'];
export type TTestSuiteSchemas = external['testsuites.yaml']['components']['schemas'];
export type TResourceSchemas = external['resources.yaml']['components']['schemas'];
export type TDataStoreSchemas = external['dataStores.yaml']['components']['schemas'];
export type TConfigSchemas = external['config.yaml']['components']['schemas'];
export type TVariablesSchemas = external['variables.yaml']['components']['schemas'];
export type TTestEventsSchemas = external['testEvents.yaml']['components']['schemas'];
export type TLintersSchemas = external['linters.yaml']['components']['schemas'];
export type TTestRunnerSchemas = external['testRunner.yaml']['components']['schemas'];
export type TWizardSchemas = external['wizards.yaml']['components']['schemas'];

export type TSelector = TTestSchemas['Selector'];

export type Modify<T, R> = Omit<T, keyof R> & R;

export type Model<T, R> = Modify<Required<T>, R>;
