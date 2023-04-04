import {TTestSchemas} from './Common.types';

export enum RunErrorTypes {
  MissingVariables = 'missingVariables',
  Unknown = 'unknown',
}

export type TTestRunState = NonNullable<TTestSchemas['TestRun']['state'] | 'WAITING' | 'SKIPPED' | 'FAILED'>;
