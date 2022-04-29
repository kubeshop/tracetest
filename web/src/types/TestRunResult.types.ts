import {TAssertionResultList} from './Assertion.types';
import {ITrace} from './Trace.types';
import {TestState} from '../constants/TestRunResult.constants';

export interface IAttribute {
  id?: string;
  key: string;
  value: string;
  type: 'span' | 'resource';
}

export interface ITestResult {
  id: string;
  successful: {
    id: string;
    operationName: string;
    duration: string;
    numOfSPans: number;
    attributes: Array<IAttribute>;
  };
  failed: {
    id: string;
    operationName: string;
    duration: string;
    numOfSPans: number;
    attributes: Array<IAttribute>;
  };
  timeStamp: Date;
}

export interface ITestRunResult {
  resultId: string;
  testId: string;
  traceId: string;
  spanId: string;
  createdAt: string;
  completedAt: string;
  response: any;
  trace: ITrace;
  state: TestState;
  assertionResultState: boolean;
  assertionResult: TAssertionResultList;
}
