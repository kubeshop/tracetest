import {TestState} from '../constants/TestRunResult.constants';
import {TAssertionResultList} from './Assertion.types';
import {Modify} from './Common.types';
import {IRawTrace, ITrace} from './Trace.types';

export interface IAttribute {
  id?: string;
  key: string;
  value: string;
  type: 'span' | 'resource';
}

export interface IRawTestRunResult {
  resultId: string;
  testId: string;
  traceId: string;
  spanId: string;
  createdAt: string;
  completedAt: string;
  response: any;
  trace?: IRawTrace;
  state: TestState;
  assertionResultState: boolean;
  assertionResult: TAssertionResultList;
}

export type ITestRunResult = Modify<IRawTestRunResult, {trace?: ITrace}>;
