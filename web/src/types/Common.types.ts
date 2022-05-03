import {external, components} from './Generated.types';

export type TRecursivePartial<T> = {
  [P in keyof T]?: TRecursivePartial<T[P]>;
};

export interface IEnv {
  measurementId: string;
  analyticsEnabled: string;
}

export type Modify<T, R> = Omit<T, keyof R> & R;

export type TraceSchemas = external['trace.yaml']['components']['schemas'];
export type Schemas = components['schemas'];
