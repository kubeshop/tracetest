import {HTTP_METHOD} from '../constants/Common.constants';
import {IAssertion} from './Assertion.types';
import {ITestRunResult} from './TestRunResult.types';

export interface IHTTPRequest {
  url: string;
  method: HTTP_METHOD;
  headers?: Array<{[key: string]: string}>;
  body?: string;
  auth?: any;
  proxy?: any;
  certificate?: any;
}

export interface ITest {
  testId: string;
  name: string;
  description: string;
  serviceUnderTest: {
    id: string;
    request: IHTTPRequest;
  };
  assertions: Array<IAssertion>;
  lastTestResult: ITestRunResult;
}
