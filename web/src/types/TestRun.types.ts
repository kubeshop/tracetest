import {TAssertionResults} from './Assertion.types';
import {Model, THttpSchemas, TTestSchemas} from './Common.types';
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
    request?: THttpSchemas['HTTPRequest'];
    response?: THttpSchemas['HTTPResponse'];
  }
>;

export type TTestRunState = TRawTestRun['state'];
