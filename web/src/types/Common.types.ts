import {external} from './Generated.types';

export type TRecursivePartial<T> = {
  [P in keyof T]?: TRecursivePartial<T[P]>;
};

export interface IEnv {
  analyticsEnabled: boolean;
  appVersion: string;
  demoEnabled: string[];
  demoEndpoints: {[key: string]: string};
  env: string;
  experimentalFeatures: string[];
  measurementId: string;
  serverID: string;
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
  identify(traits: Record<string, any>): void;
  track(event: string, traits: Record<string, any>): void;
  page(pageName: string, traits: Record<string, any>): void;
}

export declare type RecursivePartial<T> = T extends object
  ? {
      [P in keyof T]?: T[P] extends (infer U)[]
        ? RecursivePartial<U>[]
        : T[P] extends object
        ? RecursivePartial<T[P]>
        : T[P];
    }
  : any;
