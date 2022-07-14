import {external} from './Generated.types';

export type TRecursivePartial<T> = {
  [P in keyof T]?: TRecursivePartial<T[P]>;
};

export interface IEnv {
  measurementId: string;
  analyticsEnabled: string;
  serverPathPrefix: string;
}

export type Modify<T, R> = Omit<T, keyof R> & R;

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

export type TStructure = TTestSchemas['SpanSelector'];
export type TFilter = TTestSchemas['SelectorFilter'];

export type Model<T, R> = Modify<Required<T>, R>;

export interface IAnalytics {
  identify(userId: string, traits: {name: string; email: string}): void;
  track(event: string, traits: {[key: string]: any}): void;
  page(pageName: string): void;
}
