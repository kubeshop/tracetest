import {ITest} from '../types/Test.types';
import {IAssertionResult} from '../types/Assertion.types';
import {HTTP_METHOD} from '../constants/Common.constants';
import {ITestRunResult} from '../types/TestRunResult.types';
import {TestState} from '../constants/TestRunResult.constants';

interface IProps {
  test: ITest;
  assertionResult: IAssertionResult;
  testRunResult: ITestRunResult;
}

const test: ITest = {
  assertions: [],
  description: '',
  lastTestResult: undefined,
  name: '',
  serviceUnderTest: {
    id: '',
    request: {
      auth: undefined,
      body: '',
      certificate: undefined,
      headers: undefined,
      method: HTTP_METHOD.GET,
      proxy: undefined,
      url: '',
    },
  },
  testId: '',
};

const assertionResult: IAssertionResult = {
  assertion: {assertionId: '', selectors: [], spanAssertions: []},
  spanListAssertionResult: [],
};
const testRunResult: ITestRunResult = {
  assertionResult: [{assertionId: '', spanAssertionResults: []}],
  assertionResultState: false,
  completedAt: '',
  createdAt: new Date().toISOString().toString(),
  response: undefined,
  resultId: '',
  spanId: '',
  state: TestState.CREATED,
  testId: '',
  trace: undefined,
  traceId: '',
};

export const TestingModels: IProps = {
  assertionResult,
  test,
  testRunResult,
};
