import faker from '@faker-js/faker';
import {TestState} from '../../constants/TestRunResult.constants';
import {IMockFactory} from '../../types/Common.types';
import {IRawTestRunResult, ITestRunResult} from '../../types/TestRunResult.types';
import TestRunResult from '../TestRunResult.model';
import TraceMock from './Trace.mock';

const TestRunResultMock: IMockFactory<ITestRunResult, IRawTestRunResult> = () => ({
  raw(data = {}) {
    return {
      resultId: faker.datatype.uuid(),
      testId: faker.datatype.uuid(),
      traceId: faker.datatype.uuid(),
      spanId: faker.datatype.uuid(),
      createdAt: faker.date.past().toISOString(),
      completedAt: faker.date.past().toISOString(),
      response: {},
      trace: TraceMock.raw(),
      state: TestState.FINISHED,
      assertionResultState: true,
      assertionResult: [],
      ...data,
    };
  },
  model(data = {}) {
    return TestRunResult(this.raw(data));
  },
});

export default TestRunResultMock();
