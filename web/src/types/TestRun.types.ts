import {TAssertionResults} from './Assertion.types';
import {Model, TTestSchemas, TTriggerSchemas} from './Common.types';
import {TTriggerResult} from './Test.types';
import {TTestRunOutput} from './TestOutput.types';
import {TTrace} from './Trace.types';

export type TRawTestRun = TTestSchemas['TestRun'];

export type TTestRun = Model<
  TRawTestRun,
  {
    result: TAssertionResults;
    trace?: TTrace;
    totalAssertionCount: number;
    failedAssertionCount: number;
    passedAssertionCount: number;
    executionTime: number;
    lastErrorState?: string;
    trigger?: TTriggerSchemas['Trigger'];
    triggerResult?: TTriggerResult;
    outputs?: TTestRunOutput[];
  }
>;

export type TTestRunState = TRawTestRun['state'];
