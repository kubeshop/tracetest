import {TAssertionResults} from './Assertion.types';
import {Model, Modify, TTestSchemas, TTriggerSchemas} from './Common.types';
import {TEnvironment} from './Environment.types';
import {TTriggerResult} from './Test.types';
import {TTestRunOutput} from './TestOutput.types';
import {TTrace} from './Trace.types';

export type TTestRunState = NonNullable<TTestSchemas['TestRun']['state'] | 'WAITING' | 'SKIPPED'>;

export type TRawTestRun = Modify<
  TTestSchemas['TestRun'],
  {
    state?: TTestRunState;
  }
>;

export type TTestRun = Model<
  TRawTestRun,
  {
    result: TAssertionResults;
    trace?: TTrace;
    totalAssertionCount: number;
    failedAssertionCount: number;
    passedAssertionCount: number;
    executionTime: number;
    triggerTime: number;
    lastErrorState?: string;
    trigger?: TTriggerSchemas['Trigger'];
    triggerResult?: TTriggerResult;
    outputs?: TTestRunOutput[];
    environment?: TEnvironment;
    state: TTestRunState;
  }
>;
